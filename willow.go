package main

import (
	"TopEngine/common"
	"TopEngine/handler"
	"TopEngine/socket"
	"fmt"
	"sync"
)

//接下来的设计
//LRU Local Memory Cache (LLMC)
//广播更新RPC
//设计一个RPC controller
//完善分段传输
//WebSocket协议升级
//443 https ##DONE##
//ETAG  ##DONE##

// 性能阈值 { 在不记录非TopEngine占用处理时长的前提下 每请求不超过1ms的平均处理时长 }

func main() {
	//http.ListenAndServeTLS()
	config := common.ReadConf()
	wg := sync.WaitGroup{}
	for _, server := range config.ServerList {
		wg.Add(1)
		go func(params handler.Config) {
			sock := socket.Socket{}
			sock.InitSocket(fmt.Sprintf("%s:%d", params.Addr, params.Port), params)
			sock.InComingConnection()
		}(server)
	}
	handler.InitRoute()
	wg.Wait()
}
