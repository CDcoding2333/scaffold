package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors ...
func Cors(allowOrigins []string) gin.HandlerFunc {
	cnf := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "PATCH", "OPTION", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-TOKEN", "Tus-Extension", "Tus-Resumable", "Tus-Version", "Upload-Length", "Upload-Metadata", "Upload-Offset", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		AllowOrigins:     allowOrigins,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(cnf)
}
