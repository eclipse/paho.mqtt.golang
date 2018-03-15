package paho

import (
	"strings"
	"sync"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

type MessageHandler func(Message)

type Router interface {
	RegisterHandler(string, MessageHandler)
	UnregisterHandler(string)
	Route(*p.Publish)
}

type router struct {
	sync.RWMutex
	subscriptions map[string][]MessageHandler
}

func (r *router) RegisterHandler(topic string, h MessageHandler) {
	r.Lock()
	defer r.Unlock()
	debug.Println("Registering handler for:", topic)

	r.subscriptions[topic] = append(r.subscriptions[topic], h)
}

func (r *router) UnregisterHandler(topic string) {
	r.Lock()
	defer r.Unlock()
	debug.Println("Unregistering handler for:", topic)

	delete(r.subscriptions, topic)
}

func (r *router) Route(pb *p.Publish) {
	r.RLock()
	defer r.RUnlock()
	debug.Println("Routing message for:", pb.Topic)

	m := MessageFromPublish(pb)
	for route, handlers := range r.subscriptions {
		if match(route, m.Topic) {
			for _, handler := range handlers {
				handler(m)
			}
		}
	}
}

func match(route, topic string) bool {
	return route == topic || routeIncludesTopic(route, topic)
}

func matchDeep(route []string, topic []string) bool {
	if len(route) == 0 {
		if len(topic) == 0 {
			return true
		}
		return false
	}

	if len(topic) == 0 {
		if route[0] == "#" {
			return true
		}
		return false
	}

	if route[0] == "#" {
		return true
	}

	if (route[0] == "+") || (route[0] == topic[0]) {
		return matchDeep(route[1:], topic[1:])
	}
	return false
}

func routeIncludesTopic(route, topic string) bool {
	return matchDeep(routeSplit(route), topicSplit(topic))
}

func routeSplit(route string) []string {
	if len(route) == 0 {
		return nil
	}
	var result []string
	if strings.HasPrefix(route, "$share") {
		result = strings.Split(route, "/")[1:]
	} else {
		result = strings.Split(route, "/")
	}
	return result
}

func topicSplit(topic string) []string {
	if len(topic) == 0 {
		return nil
	}
	return strings.Split(topic, "/")
}

type SingleHandlerRouter struct {
	messageHandler MessageHandler
}

func (s *SingleHandlerRouter) RegisterHandler(topic string, h MessageHandler) {}

func (s *SingleHandlerRouter) UnregisterHandler(topic string) {}

func (s *SingleHandlerRouter) Route(pb *p.Publish) {
	s.messageHandler(MessageFromPublish(pb))
}
