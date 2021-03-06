// Beego (http://beego.me/)
// @description beego is an open-source, high-performance web framework for the Go programming language.
// @link        http://github.com/zhaocloud/beego for the canonical source repository
// @license     http://github.com/zhaocloud/beego/blob/master/LICENSE
// @authors     zhaocloud

package testing

import (
	"github.com/zhaocloud/beego/config"
	"github.com/zhaocloud/beego/httplib"
)

var port = ""
var baseUrl = "http://localhost:"

// beego test request client
type TestHttpRequest struct {
	httplib.BeegoHttpRequest
}

func getPort() string {
	if port == "" {
		config, err := config.NewConfig("ini", "../conf/app.conf")
		if err != nil {
			return "8080"
		}
		port = config.String("httpport")
		return port
	}
	return port
}

// returns test client in GET method
func Get(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Get(baseUrl + getPort() + path)}
}

// returns test client in POST method
func Post(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Post(baseUrl + getPort() + path)}
}

// returns test client in PUT method
func Put(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Put(baseUrl + getPort() + path)}
}

// returns test client in DELETE method
func Delete(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Delete(baseUrl + getPort() + path)}
}

// returns test client in HEAD method
func Head(path string) *TestHttpRequest {
	return &TestHttpRequest{*httplib.Head(baseUrl + getPort() + path)}
}
