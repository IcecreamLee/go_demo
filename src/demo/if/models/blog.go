package models

import (
	"Demo/if/framework"
	"time"
)

type Blog struct {
	framework.Model
	Title    string
	Date     time.Time `gorm:"default:'0000-00-00'"`
	Content  string
	Type     int
	Views    int
	IsDelete int
	AddAt    int
	UpdateAt int
	Test     float32
}

func NewBlog() *Blog {
	return &Blog{}
}

func (Blog) TableName() string {
	return "blog"
}

func (Blog) GetTable() string {
	return "blog"
}

func (blog *Blog) GetBlogs() []Blog {
	var blogs []Blog
	db, _ := framework.GetDB()
	_ = db.Find(&blogs)

	for i, blog := range blogs {
		blogs[i].Content = ""
		blog.Content = ""
	}
	return blogs
}
