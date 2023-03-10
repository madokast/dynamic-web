# dynamic-web

类似spring的web框架

根据函数类型自动注入对象

## 使用方法

1. 定义请求体和返回体格式，要求是可以 JSON 化的 struct 类型
2. 定义处理逻辑，是一个函数，入参为（请求体指针），返回值为（返回体指针 + error）
3. 使用 `web.Do(处理逻辑)` 注册路由

```go
// 请求体
type HelloRequest struct {
	Name string `json:"name"`
}

// 返回体
type HelloResponse struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

// 处理逻辑
func Hello(request *HelloRequest) (*HelloResponse, error) {
	if len(request.Name) == 0 {
		return nil, errors.New("name is empty")
	} else {
		return &HelloResponse{
			Name: request.Name,
			Msg:  "hello",
		}, nil
	}
}

func main() {
	http.HandleFunc("/hello", web.Do(Hello))
	http.ListenAndServe(":8080", nil)
}
```
