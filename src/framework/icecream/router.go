package icecream

var Routers = make(map[string]ControllerInterface)

func AddRouter(uri string, c ControllerInterface) {
	Routers[uri] = c
}
