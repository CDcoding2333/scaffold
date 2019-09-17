package wrap

import (
	"sync"
)

// WaitGroupWrapper ...
type WaitGroupWrapper struct {
	sync.WaitGroup
	ExitChan chan int
}

// Wrap ...
func (w *WaitGroupWrapper) Wrap(cb func(chan int, ...interface{}), params ...interface{}) {
	w.Add(1)
	go func() {
		cb(w.ExitChan, params...)
		w.Done()
	}()
}

// example cb:
// func cb(exitChan chan int, params ...interface{}) {
// 	// parse params
// 	for {
// 		select {
// 		case <-exitChan:
// 			return
// 		}
// 	}
// }
