package handler

import (
	"strings"
)

func HeaderToMap(header string) (headMap map[string]string, body []byte, ok bool) {
	headMap = make(map[string]string)
	ok = true

	//首先大体上解析header和body
	cutHeader := strings.Split(header, "\r\n\r\n")
	if len(cutHeader) >= 2 {
		body = []byte(cutHeader[1])
	}

	headers := strings.Split(cutHeader[0], "\r\n")
	//获得请求头部KV信息
	for _, headStr := range headers[1:] {
		cutHeader = strings.Split(headStr, ": ")
		headMap[cutHeader[0]] = cutHeader[1]
	}

	//获得请求头首行的请求方式和资源及http版本
	cutHeader = strings.Split(headers[0], " ")
	if len(cutHeader) != 3 {
		ok = false
		return
	}
	headMap["Method"] = cutHeader[0]
	headMap["Resource"] = cutHeader[1]
	headMap["Http-Version"] = cutHeader[2]
	return
}
