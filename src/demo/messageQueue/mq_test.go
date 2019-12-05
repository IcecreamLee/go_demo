package main

import (
	"math/rand"
	"testing"
)

func TestMessageStack_push(t *testing.T) {
	messageStack := MessageStack{}
	for i = 0; i < 100; i++ {
		message := Message{ID: i, Priority: rand.Intn(10)}
		messageStack.push(message)
	}

	for i, mesage := range messageStack.Messages {
		if i != 0 && mesage.Priority == messageStack.Messages[i-1].Priority && mesage.ID < messageStack.Messages[i-1].ID {
			t.Error("message push position error")
		}
		if i != len(messageStack.Messages)-1 && mesage.Priority < messageStack.Messages[i+1].Priority {
			t.Error("message push position error")
		}
	}
}
