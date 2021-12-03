package display

import (
	"runtime"
)

type Rasterizer struct {
	Workers int

	workChan chan chunk
	started  bool
}

type chunk interface {
	Process()
}

func (r *Rasterizer) Run() {
	if r.started {
		return
	}
	if r.Workers == 0 {
		r.Workers = runtime.NumCPU()
	}
	r.workChan = make(chan chunk, r.Workers)

	for i := 0; i < r.Workers; i++ {
		go r.runWorker(i)
	}
	r.started = true
}

func (r *Rasterizer) runWorker(num int) {
	for {
		c, ok := <-r.workChan
		if !ok {
			return
		}
		c.Process()
	}
}

func (r *Rasterizer) renderChunk(c chunk) {
	r.workChan <- c
}
