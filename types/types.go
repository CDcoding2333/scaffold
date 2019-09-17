package types

// WsMsgContent ...
type WsMsgContent struct {
	MsgType int         `json:"msg"`
	Data    interface{} `json:"data"`
}

// WsMsg ...
type WsMsg struct {
	WsMsgContent
	ToUsers []uint
}

// RecvWsMsg ...
type RecvWsMsg struct {
	FromUser uint        `json:"-"`
	MsgType  int         `json:"msg_type,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}
