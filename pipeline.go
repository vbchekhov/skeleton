package skeleton

import (
	"sync"
	"time"
)

// Pipeline
type Pipeline struct {
	sync.Mutex
	chatId int64
	data   []string
	rule   *Rule
	exec   map[int64]*Rule
}

// newPipeline
func newPipeline() *Pipeline {
	return &Pipeline{
		exec: make(map[int64]*Rule),
	}
}

// set()
func (p *Pipeline) set() {
	p.Lock()
	defer p.Unlock()

	if r, ok := p.exec[p.chatId]; ok {
		r.timer.Stop()
	}

	timeout := *p.rule.timeout
	if timeout == 0 {
		timeout = time.Second * 180
	}

	p.exec[p.chatId] = p.rule
	p.exec[p.chatId].timer = time.AfterFunc(timeout, func() {
		p.del()
	})
}

// get()
func (p *Pipeline) get(chatId int64) (*Rule, bool) {
	p.Lock()
	defer p.Unlock()

	rule, ok := p.exec[chatId]

	return rule, ok
}

// del()
func (p *Pipeline) del() {
	p.Lock()
	defer p.Unlock()

	p.rule.timer.Stop()
	p.data = []string{}

	delete(p.exec, p.chatId)
}

// Data get saving data
func (p *Pipeline) Data() []string {
	return p.data
}

// Save add string data in storage pipeline
func (p *Pipeline) Save(s string) {
	p.data = append(p.data, s)
}

// Timeout delete func on pipeline
func (p *Pipeline) Timeout() float64 {
	if p.rule.timeout == nil {
		return 0
	}
	return p.rule.timeout.Seconds()
}

// Prev command
func (p *Pipeline) Prev() {
	p.rule = p.rule.prev
	p.set()
}

// Next command
func (p *Pipeline) Next() {
	p.rule = p.rule.next
	p.set()
}

// Repeat command
func (p *Pipeline) Repeat() {
	p.rule = p.rule
	p.set()
}

// Stop pipeline
func (p *Pipeline) Stop() {
	p.del()
}
