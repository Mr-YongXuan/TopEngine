package common

type DynamicRoute struct {
	Routes map[string]StoreRoute
}

type StoreRoute struct {
	Method []string
	Call   Callback
}

type Request struct {
	Header   map[string]string
	Argument map[string]string
	Method   string
	Body     []byte
}

type Response struct {
	StatusCode string
	Header     map[string]string
	Body       []byte
}

type Callback func(req Request) (res Response)

func (dr *DynamicRoute) AddRoute(url string, method []string, call Callback) {
	dr.Routes[url] = StoreRoute{
		Method: method,
		Call:   call,
	}
}

func (req *Request) Fetch(key string) (res string) {
	res, _ = req.Argument[key]
	return
}

func (res *Response) SetCode(sc int) {
	res.StatusCode = StatusCode[sc]
}

func (res *Response) SetReturnJson() {
	res.Header["Content-Type"] = "application/json"
}
