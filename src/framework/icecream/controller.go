package icecream

type Controller struct {
	Request *IRequest
}

func (c *Controller) Get() {

}

func (c *Controller) Post() {

}

type ControllerInterface interface {
	Get()
	Post()
}
