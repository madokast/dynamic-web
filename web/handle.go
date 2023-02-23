package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// 日志，有需要请替换
func logError(msg string) { _ = fmt.Errorf(msg) }

type httpHandleFunc = func(w http.ResponseWriter, r *http.Request)

type Handle struct{ Method interface{} }

func Do[Request interface{}, Response interface{}](method func(request *Request) (*Response, error)) httpHandleFunc {
	h := &Handle{Method: method}
	return h.HTTPHandle()
}

func (h *Handle) HTTPHandle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 拿到的是 Request 而不是 *Request
		requestType := h.requestType()
		// 这里 new 了后，是 *Request
		request := reflect.New(requestType).Interface()
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil && !strings.Contains(err.Error(), "EOF") {
			handleError(w, err)
			return
		}
		called := reflect.ValueOf(h.Method).Call([]reflect.Value{reflect.ValueOf(request)})
		res := called[0].Interface()
		ierr := called[1].Interface()
		if err, ok := ierr.(error); ok && err != nil {
			handleError(w, err)
		} else {
			handleSuccess(w, res)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(400)
	w.Header().Set("content-type", "application/json")
	_, err = w.Write([]byte(fmt.Sprintf("{\"msg\":\"%s\"}", err.Error())))
	if err != nil {
		logError(err.Error())
	}
}

func handleSuccess(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logError(err.Error())
		_, _ = w.Write([]byte("{\"msg\":\"data marshal error\"}"))
	}
}

func (h *Handle) requestType() reflect.Type {
	// 最后取 Elem 相当于解引用
	// 本来函数 Method 是 func(*T)
	// 这里返回的是 T 类型
	return reflect.TypeOf(h.Method).In(0).Elem()
}
