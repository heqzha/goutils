package container

import (
	"sync"
)

type Queue []interface{}

func (q *Queue) Clear() {
	*q = []interface{}{}
}

func (q *Queue) Push(n interface{}) {
	*q = append(*q, n)
}

func (q *Queue) Pop() interface{} {
	if len(*q) > 0 {
		n := (*q)[0]
		*q = (*q)[1:]
		return n
	}
	return nil
}

func (q *Queue) Len() int {
	return len(*q)
}

type MtxGroupQueue struct {
	q     map[string]*Queue
	mutex *sync.RWMutex
}

func (m *MtxGroupQueue) Init() {
	m.q = make(map[string]*Queue)
	m.mutex = &sync.RWMutex{}
}

func (m *MtxGroupQueue) Push(grp string, ctnt interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.q[grp]
	if !ok {
		m.q[grp] = new(Queue)
	}
	m.q[grp].Push(ctnt)
}

func (m *MtxGroupQueue) Pop(grp string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	c, ok := m.q[grp]
	if ok && c != nil {
		return c.Pop()
	}
	return nil
}

func (m *MtxGroupQueue) Groups() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	grps := []string{}
	for k, _ := range m.q {
		grps = append(grps, k)
	}
	return grps
}

func (m *MtxGroupQueue) GroupsLen() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.q)
}

func (m *MtxGroupQueue) Len(grp string) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.q[grp].Len()
}

func (m *MtxGroupQueue) Clear(grp string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.q[grp].Clear()
}

func (m *MtxGroupQueue) ClearAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.q = make(map[string]*Queue)
}
