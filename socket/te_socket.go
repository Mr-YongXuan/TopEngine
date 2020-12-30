package socket

import (
	"TopEngine/common"
	"TopEngine/handler"
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"sync"
)

var config = common.ReadConf()
var locked = sync.Mutex{}
var logbox = common.CreateLog()

type Socket struct {
	listener    net.Listener
	connections map[net.Conn]int
	running     bool
	timeout     int
	serverName  string
	version     string
	handle      handler.HandleMethods
}

func (sock *Socket) InitSocket(address string, params handler.Config) {
	var err error
	sock.timeout = 30
	sock.handle.InitMethods(params)
	//决断 加密或不加密
	if params.Cert && params.CertCRT != "" && params.CertKEY != "" {
		//启用https
		var crt tls.Certificate
		if crt, err = tls.LoadX509KeyPair(params.CertCRT, params.CertKEY); err != nil {
			fmt.Printf("create listener failed, reason:%s\n", err)
			return
		}
		tlsConfig := &tls.Config{}
		tlsConfig.NextProtos = append(tlsConfig.NextProtos, "http/1.1")
		tlsConfig.Certificates = []tls.Certificate{crt}
		tlsConfig.Time = time.Now
		tlsConfig.Rand = rand.Reader
		if sock.listener, err = net.Listen("tcp", address); err != nil {
			fmt.Printf("create listener failed, reason:%s\n", err)
			return
		}
		sock.listener = tls.NewListener(sock.listener, tlsConfig)
		fmt.Println("TopEngine: server running at https://" + address)

	} else {
		//不启用https
		if sock.listener, err = net.Listen("tcp", address); err != nil {
			fmt.Printf("create listener failed, reason:%s\n", err)
			return
		}
		fmt.Println("TopEngine: server running at http://" + address)
	}

	sock.running = true
	go sock.keepAlive()
}

func (sock *Socket) InComingConnection() {
	var (
		err error
		cli net.Conn
	)

	for sock.running {
		if cli, err = sock.listener.Accept(); err != nil {
			fmt.Printf("can not accept client connection! reason:%s\n", err)
			continue
		}
		//拉起用户协程
		go sock.HandleProcesses(cli)
	}
}

//TODO 优化点
func (sock *Socket) keepAlive() {
	sock.connections = make(map[net.Conn]int)

	for sock.running {
		if len(sock.connections) != 0 {
			//每秒执行一次 超时的连接执行断开操作
			for cli, sec := range sock.connections {
				if sec <= 0 {
					cli.Close()
					common.Protect(func() {
						delete(sock.connections, cli)
					})
				} else {
					sock.connections[cli]--
				}

				time.Sleep(1 * time.Second)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (sock *Socket) HandleProcesses(cli net.Conn) {
	var (
		err error
		//处理数据存储用 - s
		handleHeader map[string]string
		code         string
		body         []byte
		//处理数据存储用 - e
		cliRecv  = make([]byte, 1024)
		response = make(map[string]string)
		line     = map[string]string{
			"Http-Version": "HTTP/1.1", //TODO 稍后修正
		}
	)

	for sock.running {
		if _, err = cli.Read(cliRecv); err != nil {
			cli.Close()
			common.Protect(func() {
				delete(sock.connections, cli)
			})
			//用户断开连接
			return
		}
		startTime := time.Now()
		//step.1 解析头部信息
		headMap, reqBody, ok := handler.HeaderToMap(string(cliRecv))
		if !ok {
			//405
			code = common.StatusCode[405]
			body = []byte(fmt.Sprintf(config.ErrorPage, common.StatusCode[405], common.StatusCode[405], config.ServerName))
			response["Content-Length"] = strconv.Itoa(len(body))
			response["Content-Type"] = "text/html"
		}

		//step.2 构建响应消息
		response["Server"] = config.ServerName
		response["Date"] = time.Now().UTC().Format(http.TimeFormat)
		response["Connection"] = "keep-alive"

		//step.3 资源请求处理及提交处理过程
		switch headMap["Method"] {
		case "GET":
			//执行资源调用过程
			code, handleHeader, body = sock.handle.Get(headMap["Resource"], headMap)

		case "HEAD":
			code, handleHeader, body = sock.handle.Head(headMap["Resource"], headMap)

		case "POST":
			code, handleHeader, body = sock.handle.Post(headMap["Resource"], headMap, reqBody)

		case "OPTIONS":
			code, handleHeader, body = sock.handle.Options(headMap["Resource"], headMap)
		}

		//合并数据
		line["StatusCode"] = code
		if handleHeader != nil {
			for k, v := range handleHeader {
				if k == "TE_SPLIT" {continue}
				response[k] = v
			}
		}

		//step.4 处理结束后合并结果将消息返回给客户端
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("%s %s\r\n", line["Http-Version"], line["StatusCode"]))
		for key, value := range response {
			buffer.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}

		if handleHeader["TE_SPLIT"] == "" {
			buffer.WriteString(fmt.Sprintf("\r\n%s", body))
			if _, err = cli.Write(buffer.Bytes()); err != nil {
				//发送失败则断开连接 并祈祷它的下一次到来
				cli.Close()
			}
		} else {
			//分块传输 先读取1024块
			f, _ := os.Open(handleHeader["TE_SPLIT"])
			buf := make([]byte, 4096)
			bfRd := bufio.NewReader(f)
			_, err = bfRd.Read(buf)
			buffer.WriteString(fmt.Sprintf("\r\n%s", buf))
			if _, err = cli.Write(buffer.Bytes()); err != nil {
				//发送失败则断开连接 并祈祷它的下一次到来
				cli.Close()
			}
			for {
				if _, err = bfRd.Read(buf); err == io.EOF {
					break
				}
				if _, err = cli.Write(buf); err != nil {
					//发送失败则断开连接 并祈祷它的下一次到来
					cli.Close()
				}
			}
		}

		//step.5 检查本次连接是否需要保持
		value, ok := headMap["Connection"]
		if ok {
			if value == "keep-alive" {
				//保持连接
				common.Protect(func() {
					sock.connections[cli] = sock.timeout
				})
			} else {
				//注意 此处仍旧需要判定其他连接方法
				cli.Close()
				break
			}
		}

		//性能调试计时器 REMOVE BEFORE DEPLOYMENT IF DONT NEED
		res := time.Since(startTime).Microseconds()
		if res < 1000 {
			logbox.Add(0, fmt.Sprintf("%s %s %s %s ==> Usage:%dus", cli.RemoteAddr(), line["StatusCode"], headMap["Method"], headMap["Resource"], res))
		} else {
			res /= 1000
			logbox.Add(0, fmt.Sprintf("%s %s %s %s ==> Usage:%dms", cli.RemoteAddr(), line["StatusCode"], headMap["Method"], headMap["Resource"], res))
		}
	}
}
