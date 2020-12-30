# TopEngine
   高性能一体化设计的web服务框架，使用TopEngine无需Nginx等其他web服务的介入。
 
 ### 如何使用?
   在project目录下编写你的项目, 编译后直接运行即可。
   
   <strong>框架封装的方法如下:</strong><br>
   1.HandleRoute(dr *common.DynamicRoute) {}<br>
   //该函数用于在project包中添加用户的自定义路由<br>
   //在project/te_example.go中具有使用示例<br><br>
   2.dr.AddRoute(url string, methods []string, func() {})<br>
   //该方法用于注册路由已经路由的处理方法<br>
   //在project/te_example.go中具有使用示例<br><br>
   
   <strong>路由传入请求体(req common.Request)内置方法:</strong><br>
   1.req.Header<br>
   //map数据, 存放用户请求头部信息 => 后续会更新封装方法<br>
   
   2.req.Body<br>
   //[]byte类型数据, 存放用户请求消息中的消息体。<br>
   //例如用户提交的json or file<br>
   
   3.req.Fetch(key string)<br>
   //用于从用户请求消息中抓取用户提交的表单数据<br>
   //抓取失败则返回空字符串<br><br>
   
   <strong>路由传出响应体(res common.Response)内置方法:</strong><br>
   1.res.StatusCode<br>
   //string类型数据, 用于设定服务器的返回状态码<br>
   //例如 res.StatusCode = "200 OK"<br>
   //小提示:我们已经封装好了完整的状态码仓库，你可以直接使用 common.StatusCode[200]<br>
   //例如: res.StatusCode = common.StatusCode[403]<br>
   //或者: res.setCode(403)<br>
   
   2.res.Body<br>
   //[]byte类型数据 用于存放响应结构体<br>
   //换句话说, 你需要把想让用户看到的数据扔到这里<br>
   
   3.res.SetReturnJson()<br>
   //设定服务器将content-type变为json<br>
   
   4.res.Header<br>
   //map类型数据 服务器额外响应头部信息<br>
   //请务必注意 使用前需要make初始化一下<br>
   //make(map[string]string)<br>

### 特性
   1.每个请求由TopEngine自身处理消耗的总时长平均不超过1ms.
   
   2.可扩展性 即插即用特性.

 ### 注意事项
   目前TopEngine Willow为Early Access版本, 常用HTTP协议大部分功能已经实现, 以下是尚未完成的功能列表:
   
   1.WebSocket协议升级
   
   2.HTTP 1.0版本支持
   
   3.DELETE等不常见HTTP协议方法
