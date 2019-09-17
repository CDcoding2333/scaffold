package middleware

import (
	"CDcoding2333/scaffold/errs"
	"CDcoding2333/scaffold/utils/trace"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code    uint        `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	TraceID interface{} `json:"trace_id,omitempty"`
}

const (
	succ = iota
)

// Response ...
func Response(ctx *gin.Context, result interface{}) {

	// handler error
	switch result.(type) {
	case *errs.MError:
		err := result.(*errs.MError)
		resp := &response{
			TraceID: getTraceID(ctx),
			Code:    uint(err.Code),
			Msg:     err.Error(),
		}

		response, _ := json.Marshal(resp)
		ctx.Set("response", string(response))
		ctx.AbortWithStatusJSON(200, resp)
		return
	case error:
		err := result.(error)
		resp := &response{
			TraceID: getTraceID(ctx),
			Code:    uint(errs.ErrInterior),
			Msg:     err.Error(),
		}

		response, _ := json.Marshal(resp)
		ctx.Set("response", string(response))
		ctx.AbortWithStatusJSON(200, resp)
		return
	}

	resp := &response{
		Code: uint(succ),
		Data: result,
	}
	ctx.JSON(200, resp)
}

func getTraceID(ctx *gin.Context) string {
	t, _ := ctx.Get("trace")
	traceContext, _ := t.(*trace.Context)
	if traceContext != nil {
		return traceContext.TraceID
	}
	return ""
}
