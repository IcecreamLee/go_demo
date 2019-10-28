package routers

import (
	"Demo/if/controllers"
	"net/http"
)

func init() {
	http.HandleFunc("/", controllers.Index)
}
