package core

import "framework/icecream"

type Controller struct {
	Request *icecream.IRequest
}

func (c *Controller) Get() {

}

func (c *Controller) Post() {

}

type ControllerInterface interface {
	Get()
	Post()
}
