package user

import (
	"CDcoding2333/scaffold/errs"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v9"
)

// NewUserReq ...
type NewUserReq struct {
	Name     string `json:"name" validate:"gte=2,lte=10"`
	Alias    string `json:"alias" validate:"gte=2,lte=10"`
	Password string `json:"password" validate:"gte=6"`
	Age      int    `json:"age" validate:"number"`
	Sex      int    `json:"sex" validate:""`
	Email    string `json:"email" validate:"required,email"`
	Tel      string `json:"tel" validate:"required,tel"`
	Addr     string `json:"addr" validate:""`
	Birth    int64  `json:"birth" validate:"eq=0|gte=1000000000"`
}

// LoginReq ...
type LoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// BindingValidParams ...
func (h *Handler) BindingValidParams(ctx *gin.Context, req interface{}) error {
	if err := ctx.BindJSON(req); err != nil {
		return errs.New(errs.ErrParameterInvalied, err.Error())
	}

	if err := h.validate.Struct(req); err != nil {
		sliceErrs := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			sliceErrs = append(sliceErrs, e.StructField())
		}
		return errs.New(errs.ErrParameterInvalied, fmt.Sprintf("Invalid parameter: %s", strings.Join(sliceErrs, ",")))
	}

	return nil
}

func customValiedTel(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile("^\\d{5,13}$")
	return reg.MatchString(fl.Field().String())
}
