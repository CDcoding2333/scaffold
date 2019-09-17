package ws

import (
	"CDcoding2333/scaffold/constant"
	"CDcoding2333/scaffold/core"
	"CDcoding2333/scaffold/types"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	melody "gopkg.in/olahol/melody.v1"
)

// Handler ...
type Handler struct {
	m *melody.Melody
	c *core.Core
}

// NewWebsocketHandler ...
func NewWebsocketHandler(c *core.Core) *Handler {

	h := &Handler{
		m: melody.New(),
		c: c,
	}

	h.m.Config = &melody.Config{
		WriteWait:         10 * time.Second,
		PongWait:          60 * time.Second,
		PingPeriod:        (60 * time.Second * 9) / 10,
		MaxMessageSize:    1024 * 1024,
		MessageBufferSize: 512 * 1024,
	}

	c.GetWg().Wrap(h.boardcastMsg)

	h.m.HandleConnect(h.clientConnected)
	h.m.HandleDisconnect(h.clientDisConnected)
	h.m.HandleMessage(h.dispatchMsg)
	return h
}

// WebsocketConnect ...
func (h *Handler) WebsocketConnect(ctx *gin.Context) {
	id, _ := ctx.Get(constant.ContextUserID)

	keys := map[string]interface{}{constant.ContextUserID: id}
	h.m.HandleRequestWithKeys(ctx.Writer, ctx.Request, keys)
}

func (h *Handler) clientConnected(s *melody.Session) {
	id, _ := s.Get(constant.ContextUserID)
	log.Infof("[ws] user %d connected", id)
}

func (h *Handler) clientDisConnected(s *melody.Session) {
	id, _ := s.Get(constant.ContextUserID)
	log.Infof("[ws] user %d disconnected", id)
}

// dispatchMsg receive msg from client
func (h *Handler) dispatchMsg(s *melody.Session, bs []byte) {
	var msg types.RecvWsMsg
	if err := json.Unmarshal(bs, &msg); err != nil {
		log.Infof("dispatch msg %s error:%s", string(bs), err)
		return
	}
	id, _ := s.Get(constant.ContextUserID)
	msg.FromUser = cast.ToUint(id)

	switch msg.MsgType {
	case constant.RcvMsgPing:
		log.Infof("recv [%s] from user %d", cast.ToString(msg.Data), msg.FromUser)
		h.sendMsgByFilter(&types.WsMsg{
			ToUsers: []uint{msg.FromUser},
			WsMsgContent: types.WsMsgContent{
				MsgType: constant.SendMsgPing,
				Data:    "reply: " + cast.ToString(msg.Data),
			},
		})
	}

	// userHandler, _ := h.c.GetHandler(constant.UserHandler)
	// userHandler.ReceiveMsg(&msg)

}

func (h *Handler) boardcastMsg(exitChan chan int, params ...interface{}) {
	log.Infoln("[ws] boardcastMsg start")
	msgChan := h.c.GetMsgChan()
	for {
		select {
		case msg := <-msgChan:
			h.sendMsgByFilter(msg)
		case <-exitChan:
			if len(msgChan) == 0 {
				log.Infoln("[ws] boardcastMsg stopped")
				return
			}

			for {
				msg, ok := <-msgChan
				if !ok {
					log.Infoln("[ws] boardcastMsg stopped")
					return
				}
				h.sendMsgByFilter(msg)
			}
		}
	}
}

func (h *Handler) sendMsgByFilter(msg *types.WsMsg) {
	bs, err := json.Marshal(msg.WsMsgContent)
	if err != nil {
		log.Errorf("[boardcastMsg] Marshal error :%v", err)
		return
	}

	userMap := make(map[uint]struct{})
	for _, id := range msg.ToUsers {
		userMap[id] = struct{}{}
	}

	h.m.BroadcastFilter(bs, func(s *melody.Session) bool {
		id, _ := s.Get(constant.ContextUserID)
		_, ok := userMap[cast.ToUint(id)]
		return ok
	})
}
