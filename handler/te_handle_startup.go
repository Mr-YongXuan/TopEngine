package handler

import (
	"TopEngine/common"
	"TopEngine/project"
)

var dr = &common.DynamicRoute{}
var mime = &common.MimeTable{}

//在此处调用你编写好的路由
func InitRoute() {
	//初始化MIME以及路由表
	mime.InitMime()
	dr.Routes = make(map[string]common.StoreRoute)
	//我在willow/project下创建了job1.go文件
	//于是我通过以下语法将编写好的路由进行导入到TopEngine中
	project.HandleRoute(dr)
}
