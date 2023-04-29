package mycache

import "sync"

type call struct {
	wg   sync.WaitGroup
	data interface{}
	err  error
}

type CallGroup struct {
	callMap map[string]*call
	mu      sync.Mutex
}

func (g *CallGroup) Call(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.callMap == nil {
		g.callMap = make(map[string]*call)
	}
	if call, ok := g.callMap[key]; ok {
		g.mu.Unlock()
		call.wg.Wait()
		return call.data, call.err
	}
	newCall := new(call)
	newCall.wg.Add(1)
	g.callMap[key] = newCall
	g.mu.Unlock()

	newCall.data, newCall.err = fn()
	newCall.wg.Done()

	g.mu.Lock()
	delete(g.callMap, key)
	g.mu.Unlock()

	return newCall.data, newCall.err
}
