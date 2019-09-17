package errs

// ErrorCode ...
type ErrorCode uint

//MError defined error with code
type MError struct {
	Code ErrorCode
	Msg  string
}

// Error ...
func (e *MError) Error() string {
	return e.Msg
}

const (
	// ErrInterior ... 1000
	ErrInterior ErrorCode = iota + 1000
	// ErrParameterInvalied 1001 请求参数异常
	ErrParameterInvalied
	// ErrDatabase ... 1002
	ErrDatabase
	// ErrUserName ... 1003
	ErrUserName
	// ErrPassword ... 1004
	ErrPassword
)

// New ...
func New(code ErrorCode, msg string) *MError {
	return &MError{code, msg}
}
