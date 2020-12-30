package handler

import (
	"TopEngine/common"
	"fmt"
	"strconv"
	"strings"
)

func (hm *HandleMethods) Post(resourceName string, headMap map[string]string, reqBody []byte) (statusCode string, header map[string]string, result []byte) {
	header = make(map[string]string)
	var (
		arguments = make(map[string]string)
	)

	//POST 解析
	if headMap["Content-Type"] != "application/json" {
		for _, segment := range strings.Split(string(reqBody), "&") {
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
			if !valInArray("POST", dr.Routes[url].Method) {
				statusCode = common.StatusCode[405]
				result = []byte(fmt.Sprintf(config.ErrorPage, common.StatusCode[405], common.StatusCode[405], config.ServerName))
				header["Content-Length"] = strconv.Itoa(len(result))
				header["Content-Type"] = "text/html"
				return
			}

			//构建请求体 执行动态路由
			res := dr.Routes[url].Call(common.Request{
				Header:   headMap,
				Argument: arguments,
				Method:   "POST",
				Body:     reqBody,
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
				header["Content-Type"] = "text/html; charset=utf-8"
			}
			if statusCode == "" {
				statusCode = common.StatusCode[200]
			}
			return
		}
	}

	return
}
