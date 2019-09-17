package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"CDcoding2333/scaffold/utils/trace"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RequestInLog 请求进入日志
func RequestInLog(c *gin.Context) {
	traceContext := trace.NewTrace()
	if traceID := c.Request.Header.Get("com-header-rid"); traceID != "" {
		traceContext.TraceID = traceID
	}
	if spanID := c.Request.Header.Get("com-header-spanid"); spanID != "" {
		traceContext.SpanID = spanID
	}

	c.Set("startExecTime", time.Now())
	c.Set("trace", traceContext)

	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back

	m := trace.GetLogFileds(traceContext, trace.DLTagRequestIn)
	m["uri"] = c.Request.RequestURI
	m["method"] = c.Request.Method
	m["args"] = c.Request.PostForm
	m["body"] = string(bodyBytes)
	m["from"] = c.ClientIP()

	log.WithFields(m).Info()
}

// RequestOutLog 请求输出日志
func RequestOutLog(c *gin.Context) {
	// after request
	endExecTime := time.Now()
	response, _ := c.Get("response")
	st, _ := c.Get("startExecTime")

	startExecTime, _ := st.(time.Time)

	m := trace.GetLogFileds(trace.GetGinTraceContext(c), trace.DLTagRequestOut)
	m["uri"] = c.Request.RequestURI
	m["method"] = c.Request.Method
	m["args"] = c.Request.PostForm
	m["response"] = response
	m["from"] = c.ClientIP()
	m["proc_time"] = endExecTime.Sub(startExecTime).Seconds()

	log.WithFields(m).Info()
}

// RequestLog ...
func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		RequestInLog(c)
		defer RequestOutLog(c)
		c.Next()
	}
}
