package main

import "testing"

func TestMessageStack_insert(t *testing.T) {
	message := Message{Priority: 1}
	messageStack := MessageStack{}
	index := messageStack.push(message)
	if index != 0 {
		t.Error("find error:", index)
	}
	message = Message{Priority: 2}
	index = messageStack.push(message)
	if index != 0 {
		t.Error("find error:", index)
	}
	message = Message{Priority: 3}
	index = messageStack.push(message)
	if index != 0 {
		t.Error("find error:", index)
	}
	message = Message{Priority: 4}
	index = messageStack.push(message)
	if index != 0 {
		t.Error("find error:", index)
	}
	message = Message{Priority: 4}
	index = messageStack.push(message)
	if index != 1 {
		t.Error("find error:", index)
	}
	message = Message{Priority: 3}
	index = messageStack.push(message)
	if index != 3 {
		t.Error("find error:", index)
	}
	message = Message{Priority: 1}
	index = messageStack.push(message)
	if index != 6 {
		t.Error("find error:", index)
	}
}
