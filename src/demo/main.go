package main

import "fmt"

type user struct {
	name  string
	email string
}

func (u user) notify() {
	fmt.Printf("sending user email to %s<%s>\n", u.name, u.email)
}

func (u *user) changeEmail(email string) {
	u.email = email
}

func main() {
	bob := user{"bob", "bob@email.com"}
	bob.notify()

	bill := &user{"bill", "bill@email.com"}
	bill.notify()

	bob.changeEmail("bob@new.com")
	bob.notify()
	bill.changeEmail("bill@new.com")
	bill.notify()

	fmt.Printf("#v", bob)
	fmt.Printf("#v", bill)
	fmt.Printf("#v", &bill)
}
