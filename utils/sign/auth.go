package sign

import (
	"CDcoding2333/scaffold/constant"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

// Auth ...
type Auth struct {
	iss    string
	secret []byte
}

type jwtClaims struct {
	UserID uint `json:"user_id"`
}

// NewAuth ...
func NewAuth(iss, secret string) *Auth {
	return &Auth{
		iss:    iss,
		secret: []byte(secret),
	}
}

// GenAuth 生成jwt的token
func (a *Auth) GenAuth(id uint) string {
	exp := time.Now().UTC().Add(12 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"iss":     a.iss,
		"exp":     exp,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(a.secret)
	return "Bearer " + tokenString
}

//HandleAuth ...
func (a *Auth) HandleAuth(ctx *gin.Context) {
	id, err := a.decryptAuth(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set(constant.ContextUserID, id)
	ctx.Next()
}

//BuildGRPCAUTH grpc authc
func (a *Auth) BuildGRPCAUTH() func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("without Authorization")
		}

		Authorization := md.Get("Authorization")
		if len(Authorization) == 0 {
			return nil, errors.New("Authorization error")
		}

		id, err := a.decryptAuth(Authorization[0])
		if err != nil {
			return nil, err
		}

		return context.WithValue(ctx, constant.ContextUserID, id), nil
	}
}

func (a *Auth) decryptAuth(authInfo string) (uint, error) {
	if authInfo == "" {
		return 0, errors.New("auth is null")
	}

	authArray := strings.Split(authInfo, " ")
	if len(authArray) != 2 {
		return 0, errors.New("auth illegal")
	}

	// RFC 6750
	if authArray[0] != "Bearer" {
		return 0, errors.New("auth illegal")
	}

	info := struct {
		jwtClaims
		jwt.StandardClaims
	}{}

	token, err := jwt.ParseWithClaims(authArray[1], &info, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})

	if !token.Valid || err != nil {
		return 0, errors.New("param error")
	}

	return info.UserID, nil
}
