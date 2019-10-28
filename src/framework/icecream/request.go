package icecream

import (
	"fmt"
	"net/http"
)

func NewRequest(w http.ResponseWriter, r *http.Request) {
	request := new(IRequest)
	request.Request = r
	request.Response = w
	request.exec()

}

type IRequest struct {
	Request  *http.Request
	Response http.ResponseWriter
}

func (request *IRequest) exec() {
	if _, ok := Routers[request.Request.URL.Path]; ok {
		if request.Request.Method == "GET" {

			fmt.Printf("Controller: %#v\n", Routers[request.Request.URL.Path])
		}
	}
	fmt.Printf("URL: %s", request.Request.URL.Path)
}
