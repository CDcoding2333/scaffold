package constant

const (
	// ContextUserID ...
	ContextUserID = "__user_id__"
)

const (
	// UserStateActive ...
	UserStateActive = iota + 1
	// UserStateDisable ...
	UserStateDisable
	// UserStateDelete ...
	UserStateDelete
)

const (
	// UserHandler ...
	UserHandler = iota + 1
)

const (
	// RcvMsgPing ...
	RcvMsgPing = iota
)

const (
	// SendMsgPing ...
	SendMsgPing = iota
)
