package core

import (
	"CDcoding2333/scaffold/types"
	"CDcoding2333/scaffold/utils/wrap"
)

// Core ...
type Core struct {
	handlers  map[int]Handler
	msgChan   chan *types.WsMsg
	waitGroup *wrap.WaitGroupWrapper
}

// Handler ...
type Handler interface {
	ReceiveMsg(*types.RecvWsMsg)
}

// New ...
func New() *Core {
	return &Core{
		handlers: make(map[int]Handler),
		msgChan:  make(chan *types.WsMsg, 10),
		waitGroup: &wrap.WaitGroupWrapper{
			ExitChan: make(chan int),
		},
	}
}

// RegistHandlers ...
func (c *Core) RegistHandlers(handlerName int, h Handler) {
	if _, ok := c.handlers[handlerName]; ok {
		panic("handlerName has been registed")
	}
	c.handlers[handlerName] = h
}

// GetHandler ...
func (c *Core) GetHandler(handlerName int) (Handler, bool) {
	v, ok := c.handlers[handlerName]
	return v, ok
}

// GetMsgChan ...
func (c *Core) GetMsgChan() chan *types.WsMsg {
	return c.msgChan
}

// GetWg ...
func (c *Core) GetWg() *wrap.WaitGroupWrapper {
	return c.waitGroup
}

// SendMsg ...
func (c *Core) SendMsg(msg *types.WsMsg) {
	c.msgChan <- msg
}
