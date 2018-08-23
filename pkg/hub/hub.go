// Package hub implements the methods required to manage subscribers by mapping
// topic strings to subscriber chans.
package hub

import (
	"sync"
)

// Hub instances map topic strings to subscriber chans and provide methods for
// managing subscribers and publishing.
type Hub struct {
	sync.RWMutex
	topicChans map[string]map[chan []byte]struct{}
	chanTopics map[chan []byte]map[string]struct{}
}

// New returns a new Hub instance.
func New() *Hub {
	h := &Hub{}
	h.topicChans = make(map[string]map[chan []byte]struct{})
	h.chanTopics = make(map[chan []byte]map[string]struct{})
	return h
}

// Subscribe chan c to topic t.
func (h *Hub) Subscribe(t string, c chan []byte) {
	h.Lock()
	defer h.Unlock()
	chans, ok := h.topicChans[t]
	if !ok {
		chans = make(map[chan []byte]struct{})
		h.topicChans[t] = chans
	}
	chans[c] = struct{}{}
	topics, ok := h.chanTopics[c]
	if !ok {
		topics = make(map[string]struct{})
		h.chanTopics[c] = topics
	}
	topics[t] = struct{}{}
}

// Unsubscribe chan c from topic t.
func (h *Hub) Unsubscribe(t string, ch chan []byte) {
	h.Lock()
	defer h.Unlock()
	chans, ok := h.topicChans[t]
	if ok {
		delete(chans, ch)
	}
	topics, ok := h.chanTopics[ch]
	if ok {
		delete(topics, t)
	}
}

// UnsubscribeAll unsubscribes chan c from all topics.
func (h *Hub) UnsubscribeAll(ch chan []byte) {
	h.Lock()
	defer h.Unlock()
	topics, ok := h.chanTopics[ch]
	if !ok {
		return
	}
	for t := range topics {
		chans, ok := h.topicChans[t]
		if ok {
			delete(chans, ch)
		}
	}
	delete(h.chanTopics, ch)
}

// Publish bytes b to topic t.
func (h *Hub) Publish(t string, b []byte) {
	h.RLock()
	defer h.RUnlock()
	chans, ok := h.topicChans[t]
	if !ok {
		return
	}
	for ch := range chans {
		select {
		case ch <- b:
		default:
			continue
		}
	}
}
