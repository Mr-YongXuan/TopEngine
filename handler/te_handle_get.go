package handler

import (
	"TopEngine/common"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func (hm *HandleMethods) Get(resourceName string, headMap map[string]string) (statusCode string, header map[string]string, result []byte) {
	header = make(map[string]string)
	var (
		readBuf   []byte
		err       error
		arguments = make(map[string]string)
	)

	//GET 请求参数处理
	pubTmp := strings.Split(resourceName, "?")
	resourceName = pubTmp[0]
	if len(pubTmp) > 1 {
		for _, segment := range strings.Split(pubTmp[1], "&") {
			kv := strings.Split(segment, "=")
			if len(kv) != 2 {
				continue
			}
			arguments[kv[0]] = kv[1]
		}
	}

	//动态路由处理
	for url := range dr.Routes {
		if url == resourceName {
			//检查请求方法是否被允许
			if !valInArray("GET", dr.Routes[url].Method) {
				statusCode = common.StatusCode[405]
				result = []byte(fmt.Sprintf(config.ErrorPage, common.StatusCode[405], common.StatusCode[405], config.ServerName))
				header["Content-Length"] = strconv.Itoa(len(result))
				header["Content-Type"] = "text/html; " + hm.Params.Charset
				return
			}

			//构建请求体 执行动态路由
			res := dr.Routes[url].Call(common.Request{
				Header:   headMap,
				Argument: arguments,
				Method:   "GET",
			})
			//合并响应消息
			if res.Header != nil {
				for k, v := range res.Header {
					header[k] = v
				}
			}
			//重构 & 发送
			result = res.Body
			header["Content-Length"] = strconv.Itoa(len(result))
			if _, ok := header["Content-Type"]; !ok {
				header["Content-Type"] = "text/html; " + hm.Params.Charset
			}
			statusCode = common.StatusCode[200]
			return
		}
	}

	absPath, _ := os.Getwd()
	absPath = path.Join(absPath, hm.Params.Root, resourceName)
	fileExt := true
	fileStat := isFile(absPath)

	//文件不存在 || 目录不存在
	if fileStat == 0 {
		statusCode = common.StatusCode[404]
		result = []byte(fmt.Sprintf(config.ErrorPage, common.StatusCode[404], common.StatusCode[404], config.ServerName))
		header["Content-Length"] = strconv.Itoa(len(result))
		header["Content-Type"] = "text/html; " + hm.Params.Charset
		return
	}

	//目标为目录
	if fileStat == 1 {
		fileExt = false
		for _, index := range hm.Params.IndexPage {
			if isFile(path.Join(absPath, index)) == 2 {
				absPath = path.Join(absPath, index)
				fileExt = true
				break
			}
		}
	}

	//目录中不存在索引文件
	if !fileExt {
		statusCode = common.StatusCode[404]
		result = []byte(fmt.Sprintf(config.ErrorPage, common.StatusCode[404], common.StatusCode[404], config.ServerName))
		header["Content-Length"] = strconv.Itoa(len(result))
		header["Content-Type"] = "text/html; " + hm.Params.Charset
		return
	}

	//缓存检查 304
	if val, ok := headMap["If-None-Match"]; ok {
		if val == fmt.Sprintf("\"%x-%d\"", getLastMod(absPath), getFileSize(absPath)) {
			//return 304
			statusCode = common.StatusCode[304]
			return
		}

	} else if val, ok := headMap["If-Modify-Since"]; ok {
		if val == time.Unix(getLastMod(absPath), 0).Format(http.TimeFormat) {
			//return 304
			statusCode = common.StatusCode[304]
			return
		}
	}

	//缓存未命中
	//决断 分块or整体发送
	if getFileSize(absPath) > 1024 {
		header["TE_SPLIT"] = absPath
	} else {
		header["TE_SPLIT"] = ""
		if readBuf, err = ioutil.ReadFile(absPath); err != nil {
			statusCode = common.StatusCode[403]
			result = []byte(fmt.Sprintf(config.ErrorPage, common.StatusCode[403], common.StatusCode[403], config.ServerName))
			header["Content-Length"] = strconv.Itoa(len(result))
			header["Content-Type"] = "text/html; " + hm.Params.Charset
			return
		}
		result = readBuf
	}
	statusCode = common.StatusCode[200]
	header["Content-Length"] = fmt.Sprintf("%d", getFileSize(absPath))

	//创建该资源的mime
	header["Content-Type"] = mime.GetMime(absPath)
	if strings.Split(header["Content-Type"], "/")[0] != "text" {
		header["Accept-Ranges"] = "bytes"
	} else {
		header["Content-Type"] += "; " + hm.Params.Charset
	}

	//最后修改时间
	header["Last-Modified"] = time.Unix(getLastMod(absPath), 0).Format(http.TimeFormat)

	//构建资源ETAG标签
	header["ETag"] = fmt.Sprintf("\"%x-%s\"", getLastMod(absPath), header["Content-Length"])

	return
}
