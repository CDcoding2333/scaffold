package user

import (
	"CDcoding2333/scaffold/constant"
	"CDcoding2333/scaffold/core"
	"CDcoding2333/scaffold/dao"
	"CDcoding2333/scaffold/errs"
	"CDcoding2333/scaffold/pkg/uuid"
	"CDcoding2333/scaffold/types"
	"CDcoding2333/scaffold/utils/middleware"
	"CDcoding2333/scaffold/utils/misc"
	"CDcoding2333/scaffold/utils/sign"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	validator "gopkg.in/go-playground/validator.v9"
)

// Handler ...
type Handler struct {
	Auth     *sign.Auth
	c        *core.Core
	validate *validator.Validate
}

// NewUserHandler ...
func NewUserHandler(auth *sign.Auth, c *core.Core) *Handler {
	validate := validator.New()
	validate.RegisterValidation("tel", customValiedTel)
	return &Handler{
		Auth:     auth,
		c:        c,
		validate: validate,
	}
}

// UsersRegist ...
func (h *Handler) UsersRegist(ctx *gin.Context) {
	req := &NewUserReq{}
	if err := h.BindingValidParams(ctx, req); err != nil {
		middleware.Response(ctx, err)
		return
	}

	uuid, err := uuid.GenUUID()
	if err != nil {
		log.WithError(err).Error("NewUser:uuid error")
		middleware.Response(ctx, err)
		return
	}

	salt := misc.RandString(4)
	user := dao.User{
		UID:      uuid,
		Name:     req.Name,
		Alias:    req.Alias,
		Age:      req.Age,
		Sex:      req.Sex,
		Email:    req.Email,
		Tel:      req.Tel,
		Addr:     req.Addr,
		State:    constant.UserStateActive,
		Salt:     salt,
		Password: misc.CryptWithMd5(req.Password, salt),
	}

	if req.Birth != 0 {
		t := time.Unix(req.Birth, 0)
		user.Birth = &t
	}

	if err := user.Create(); err != nil {
		middleware.Response(ctx, err)
		return
	}

	ctx.Writer.Header().Set("Authorization", h.Auth.GenAuth(user.ID))
	middleware.Response(ctx, user)
}

// UsersLogin ...
func (h *Handler) UsersLogin(ctx *gin.Context) {
	req := &LoginReq{}
	if err := h.BindingValidParams(ctx, req); err != nil {
		middleware.Response(ctx, err)
		return
	}

	user := dao.User{Name: req.Name}
	if err := user.GetUserByName(); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			middleware.Response(ctx, errs.New(errs.ErrUserName, "invalied name"))
			return
		}
		middleware.Response(ctx, err)
		return
	}

	if user.State != constant.UserStateActive {
		middleware.Response(ctx, errs.New(errs.ErrUserName, "invalied user state"))
		return
	}

	if !misc.VerifyHashWithMd5(req.Password, user.Password, user.Salt) {
		middleware.Response(ctx, errs.New(errs.ErrPassword, "invalied password"))
		return
	}

	ctx.Writer.Header().Set("Authorization", h.Auth.GenAuth(user.ID))
	middleware.Response(ctx, user)
}

// GetUserInfo ....
func (h *Handler) GetUserInfo(ctx *gin.Context) {
	id, _ := ctx.Get(constant.ContextUserID)

	user := dao.User{}
	users, err := user.GetUsers(cast.ToUint(id))
	if err != nil {
		middleware.Response(ctx, err)
		return
	}

	if len(users) == 0 {
		middleware.Response(ctx, errs.New(errs.ErrParameterInvalied, "invalied parameter id"))
		return
	}

	middleware.Response(ctx, users[0])
}

// DelUsers ...
func (h *Handler) DelUsers(ctx *gin.Context) {

}

// ReceiveMsg ...
func (h Handler) ReceiveMsg(msg *types.RecvWsMsg) {

}
