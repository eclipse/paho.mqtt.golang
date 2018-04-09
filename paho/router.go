package paho

import (
	"strings"
	"sync"

	p "github.com/eclipse/paho.mqtt.golang/packets"
)

// MessageHandler is a type for a function that is invoked
// by a Router when it has received a Publish.
type MessageHandler func(Message)

// Router is an interface of the functions for a struct that is
// used to handle invoking MessageHandlers depending on the
// the topic the message was published on.
// RegisterHandler() takes a string of the topic, and a MessageHandler
// to be invoked when Publishes are received that match that topic
// UnregisterHandler() takes a string of the topic to remove
// MessageHandlers for
// Route() takes a Publish message and determines which MessageHandlers
// should be invoked
type Router interface {
	RegisterHandler(string, MessageHandler)
	UnregisterHandler(string)
	Route(*p.Publish)
}

// StandardRouter is a library provided implementation of a Router that
// allows for unique and multiple MessageHandlers per topic
type StandardRouter struct {
	sync.RWMutex
	subscriptions map[string][]MessageHandler
}

// RegisterHandler is the library provided StandardRouter's
// implementation of the required interface function()
func (r *StandardRouter) RegisterHandler(topic string, h MessageHandler) {
	r.Lock()
	defer r.Unlock()
	debug.Println("Registering handler for:", topic)

	r.subscriptions[topic] = append(r.subscriptions[topic], h)
}

// UnregisterHandler is the library provided StandardRouter's
// implementation of the required interface function()
func (r *StandardRouter) UnregisterHandler(topic string) {
	r.Lock()
	defer r.Unlock()
	debug.Println("Unregistering handler for:", topic)

	delete(r.subscriptions, topic)
}

// Route is the library provided StandardRouter's implementation
// of the required interface function()
func (r *StandardRouter) Route(pb *p.Publish) {
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

// SingleHandlerRouter is a library provided implementation of a Router
// that stores only a single MessageHandler and invokes this MessageHandler
// for all received Publishes
type SingleHandlerRouter struct {
	messageHandler MessageHandler
}

// RegisterHandler is the library provided SingleHandlerRouter's
// implementation of the required interface function()
func (s *SingleHandlerRouter) RegisterHandler(topic string, h MessageHandler) {}

// UnregisterHandler is the library provided SingleHandlerRouter's
// implementation of the required interface function()
func (s *SingleHandlerRouter) UnregisterHandler(topic string) {}

// Route is the library provided SingleHandlerRouter's
// implementation of the required interface function()
func (s *SingleHandlerRouter) Route(pb *p.Publish) {
	s.messageHandler(MessageFromPublish(pb))
}
