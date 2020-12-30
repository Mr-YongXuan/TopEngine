package handler

import (
	"TopEngine/common"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

func (hm *HandleMethods) Options(resourceName string, headMap map[string]string) (statusCode string, header map[string]string, body []byte) {
	header = make(map[string]string)
	var (
		result  []byte
		readBuf []byte
		err     error
	)

	//动态路由处理
	for url := range dr.Routes {
		if url == resourceName {
			header["Content-Length"] = strconv.Itoa(len(result))
			header["Content-Type"] = "text/html; charset=utf-8"
			header["Allow"] = ""
			for _, method := range dr.Routes[url].Method {
				header["Allow"] += method + ", "
			}
			header["Allow"] += "OPTIONS"

			//构建动态路由头部
			statusCode = common.StatusCode[200]
			res := dr.Routes[url].Call(common.Request{
				Header: headMap,
				Method: "OPTIONS",
			})
			//合并响应消息
			if res.Header != nil {
				for k, v := range res.Header {
					header[k] = v
				}
			}
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
		header["Content-Type"] = "text/html"
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

	if !fileExt {
		statusCode = common.StatusCode[404]
		result = []byte(fmt.Sprintf(config.ErrorPage, common.StatusCode[404], common.StatusCode[404], config.ServerName))
		header["Content-Length"] = strconv.Itoa(len(result))
		header["Content-Type"] = "text/html"
		return
	}

	if readBuf, err = ioutil.ReadFile(absPath); err != nil {
		statusCode = common.StatusCode[403]
		result = []byte(fmt.Sprintf(config.ErrorPage, common.StatusCode[403], common.StatusCode[403], config.ServerName))
		header["Content-Length"] = strconv.Itoa(len(result))
		header["Content-Type"] = "text/html"
		return
	}
	statusCode = common.StatusCode[200]
	result = readBuf
	header["Content-Length"] = strconv.Itoa(len(result))
	header["Allow"] = "GET, POST, HEAD, OPTIONS"
	//创建该资源的mime
	header["Content-Type"] = mime.GetMime(absPath)
	if strings.Split(header["Content-Type"], "/")[0] != "text" {
		header["Accept-Ranges"] = "bytes"
	}

	return
}
