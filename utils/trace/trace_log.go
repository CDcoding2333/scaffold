package trace

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 通用DLTag常量定义
const (
	DLTagRequestIn  = "_com_request_in"
	DLTagRequestOut = "_com_request_out"
)

const (
	_dlTag       = "dltag"
	_traceID     = "traceid"
	_spanID      = "spanid"
	_childSpanID = "cspanid"
)

// GetLogFileds ...
func GetLogFileds(trace *Context, dltag string) log.Fields {
	m := make(map[string]interface{}, 0)
	m[_dlTag] = dltag
	m[_traceID] = trace.TraceID
	m[_childSpanID] = trace.CSpanID
	m[_spanID] = trace.SpanID
	return m
}

// GetGinTraceContext ...从gin的Context中获取数据
func GetGinTraceContext(c *gin.Context) *Context {
	// 防御
	if c == nil {
		return NewTrace()
	}
	traceContext, exists := c.Get("trace")
	if exists {
		if tc, ok := traceContext.(*Context); ok {
			return tc
		}
	}
	return NewTrace()
}
