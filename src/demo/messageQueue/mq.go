package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// 消息队列单个消息结构体
type Message struct {
	ID           int
	Topic        string
	Data         string
	Priority     int
	Status       int
	Result       int
	ExecuteTime  int64
	ExecutedTime int64
}

type MessageStack struct {
	Messages []Message
	Lock     sync.RWMutex
}

func (m *MessageStack) insert(message Message) int {
	i := 0
	left := 0
	right := len(m.Messages) - 1
	rightIndex := right

	for {
		if right < 0 {
			i = 0
			break
		} else if left >= right {
			if message.Priority > m.Messages[right].Priority {
				i = right
				break
			}
			i = right + 1
			break
		}

		middle := int(math.Floor(float64((left + right) / 2)))

		if message.Priority == m.Messages[middle].Priority && message.Priority > m.Messages[middle+1].Priority {
			i = middle + 1
			break
		} else if message.Priority > m.Messages[middle].Priority {
			right = middle - 1
		} else {
			left = middle + 1
		}
	}

	if rightIndex < i {
		m.Messages = append(m.Messages, message)
		return i
	}

	messagesTmp := append([]Message{message}, m.Messages[i:]...)
	messagesTmp2 := m.Messages[:i]
	m.Messages = append(messagesTmp2, messagesTmp...)
	return i
}

// 延迟队列消息
type DelayJobStack struct {
	Messages []Message
	Lock     sync.RWMutex
}

func (m *DelayJobStack) insert(message Message) int {
	i := 0
	left := 0
	right := len(m.Messages) - 1
	rightIndex := right

	for {
		if right < 0 {
			i = 0
			break
		} else if left >= right {
			if message.ExecuteTime < m.Messages[right].ExecuteTime {
				i = right
				break
			}
			i = right + 1
			break
		}

		middle := int(math.Floor(float64((left + right) / 2)))

		if message.ExecuteTime == m.Messages[middle].ExecuteTime && message.ExecuteTime < m.Messages[middle+1].ExecuteTime {
			i = middle + 1
			break
		} else if message.ExecuteTime < m.Messages[middle].ExecuteTime {
			right = middle - 1
		} else {
			left = middle + 1
		}
	}

	if rightIndex < i {
		m.Messages = append(m.Messages, message)
		return i
	}

	messagesTmp := append([]Message{message}, m.Messages[i:]...)
	messagesTmp2 := m.Messages[:i]
	m.Messages = append(messagesTmp2, messagesTmp...)
	return i
}

// MessageQueue 消息队列主体
type MessageQueue struct {
	MessageStack       MessageStack
	MaximumConcurrency int // 最大并发执行消息数
}

var i int

// publish 生产一个消息
func (m *MessageQueue) publish(message Message) {
	m.MessageStack.Lock.Lock()
	defer m.MessageStack.Lock.Unlock()
	i++
	message.ID = i
	m.MessageStack.insert(message)
}

// consume 消费一个消息
func (m *MessageQueue) consume() (Message, error) {
	m.MessageStack.Lock.Lock()
	defer m.MessageStack.Lock.Unlock()
	if len(m.MessageStack.Messages) == 0 {
		return Message{}, errors.New("no consume")
	}
	job := m.MessageStack.Messages[0]
	m.MessageStack.Messages = m.MessageStack.Messages[1:]
	return job, nil
}

// execJob 执行单个消息任务
func (m *MessageQueue) execJob(job Message) {
	fmt.Println("execute job:", job)
	time.Sleep(time.Duration(3000+rand.Intn(1000)) * time.Millisecond)
	fmt.Println("executed job:", job)
}

// run 消息队列持续的运行
func (m *MessageQueue) Run() {
	fmt.Println("MessageQueue running...")

	jobsChan := make(chan Message, m.MaximumConcurrency)

	go func() {
		for {
			job, err := m.consume()
			if err == nil {
				jobsChan <- job
			}
		}
	}()

	for i := 0; i < m.MaximumConcurrency; i++ {
		go func() {
			for job := range jobsChan {
				m.execJob(job)
			}
		}()
	}
}

var mqOnce sync.Once
var mq *MessageQueue

// 获取MessageQueue单例
func MQSingleton() *MessageQueue {
	mqOnce.Do(func() {
		mq = &MessageQueue{MaximumConcurrency: 3}
	})
	return mq
}
