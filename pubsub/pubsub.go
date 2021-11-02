package pubsub

import (
	"sync"
	"time"
)

type (
	subscriber chan interface{}         // 订阅者为一个通道
	topicFunc  func(v interface{}) bool // 主题即是过滤器
)

// Publisher 发布者对象
type Publisher struct {
	m           sync.RWMutex             // 读写锁
	buffer      int                      // 订阅者缓冲长度
	timeout     time.Duration            // 发布允许延迟
	subscribers map[subscriber]topicFunc // 订阅者
}

// NewPublisher 构建一个发布者
func NewPublisher(buffer int, timeout time.Duration) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     timeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// Subscribe 添加一个订阅者，订阅全部主题
func (p *Publisher) Subscribe() subscriber {
	return p.SubscribeTopic(nil)
}

// SubscribeTopic 添加一个订阅者，订阅过滤后主题
func (p *Publisher) SubscribeTopic(topic topicFunc) subscriber {
	sub := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[sub] = topic
	p.m.Unlock()
	return sub
}

// Publish 发布一个主题
func (p *Publisher) Publish(v interface{}) {
	// 用到读写锁
	p.m.RLock()
	defer p.m.RUnlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, &wg, v)
	}
	wg.Wait()
}

// 往订阅者通道中发送通过过滤器过滤的主题, 发送过程允许一定的延迟
func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, wg *sync.WaitGroup, v interface{}) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}
	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

// Unsubscribe 退订
func (p *Publisher) Unsubscribe(sub subscriber) {
	p.m.Lock()
	defer p.m.Unlock()
	delete(p.subscribers, sub)
	close(sub)
}

// Close 关闭发布者通道
func (p *Publisher) Close() {

	p.m.Lock()
	defer p.m.Unlock()

	// 关闭所有的订阅者
	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}
