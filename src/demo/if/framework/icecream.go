package framework

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type Icecream struct {
	DB *gorm.DB
}

var icecream *Icecream
var once sync.Once

func GetIcecream() *Icecream {
	once.Do(func() {
		icecream = new(Icecream)
	})
	return icecream
}
