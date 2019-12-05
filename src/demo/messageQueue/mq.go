package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
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

// 优先级消息队列
type MessageStack struct {
	Messages []Message
	Lock     sync.RWMutex
}

// push 插入消息到队列中正确的位置
func (m *MessageStack) push(message Message) int {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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

// pop 弹出队列中第一个消息
func (m *MessageStack) pop() (Message, bool) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if len(m.Messages) == 0 {
		return Message{}, false
	}
	message := m.Messages[0]
	m.Messages = m.Messages[1:]
	return message, true
}

// 延迟队列消息
type DelayMessageStack struct {
	MessageStack
}

func (m *DelayMessageStack) push(message Message) int {
	m.Lock.Lock()
	defer m.Lock.Unlock()

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
	MessageStack       MessageStack      // 优先级消息队列
	DelayMessageStack  DelayMessageStack // 延迟消息队列
	MaximumConcurrency int               // 最大并发执行消息数
}

var i int

// publish 生产一个消息
func (m *MessageQueue) publish(message Message) {
	i++
	message.ID = i
	if message.ExecuteTime > time.Now().Unix() {
		m.DelayMessageStack.push(message)
	} else {
		m.MessageStack.push(message)
	}
}

// consume 消费一个消息
func (m *MessageQueue) consume() (Message, bool) {
	return m.MessageStack.pop()
}

// execJob 执行单个消息任务
func (m *MessageQueue) execJob(job Message) {
	fmt.Println("Execute job:", job, ", Goroutine num:", runtime.NumGoroutine())
	time.Sleep(time.Duration(3000+rand.Intn(1000)) * time.Millisecond)
	fmt.Println("Executed job:", job, ", Goroutine num:", runtime.NumGoroutine())
}

// run 消息队列持续的运行
func (m *MessageQueue) Run() {
	fmt.Println("MessageQueue(maximumConcurrency:" + strconv.Itoa(m.MaximumConcurrency) + ") running...")

	jobsChan := make(chan Message, m.MaximumConcurrency)

	go func() {
		for {
			message, ok := m.DelayMessageStack.pop()
			if ok {
				m.publish(message)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	go func() {
		for {
			job, ok := m.consume()
			if ok {
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
		mq = &MessageQueue{MaximumConcurrency: conf.MaximumConcurrency}
	})
	return mq
}
