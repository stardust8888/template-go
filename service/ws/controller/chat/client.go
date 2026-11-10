package chat

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"template-go/library/ginx"
	"template-go/logs"
	"template-go/model"
	"template-go/utils/constant"
	"template-go/utils/x/slicex"
	"template-go/utils/x/timex"
)

const (
	PongTimeout  = 60 * time.Second
	MaxReadLimit = 3 * constant.KB
	SendChannel  = 128
)

type Client struct {
	conn      *websocket.Conn
	user      *constant.SimpleUser
	send      chan []byte
	done      chan struct{}
	closeOnce sync.Once
}

func (c *Client) cleanup() {
	c.closeOnce.Do(func() {
		close(c.done)
		_ = c.conn.Close()
		WriteTimeout(mediator.unregister, c, WriteChannelTime)
	})
}

func (c *Client) Read() error {
	c.conn.SetReadLimit(MaxReadLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(PongTimeout))
	defer c.cleanup()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return err
			}
			logs.Error(ginx.Context, "ConnRead error:", zap.Error(err))
			continue
		}

		if string(msg) == "pong" {
			_ = c.conn.SetReadDeadline(time.Now().Add(PongTimeout))
			continue
		}
		logs.Debug(ginx.Context, "Read:", zap.String("msg", string(msg)))

		message := &model.Message{}
		err = message.Unmarshal(msg)
		if err != nil {
			m := model.NewTipMessage(constant.CodeMessageContentError)
			_ = c.conn.WriteMessage(websocket.TextMessage, m.Marshal())
			logs.Error(ginx.Context, "ConnRead Unmarshal error:", zap.Error(err), zap.String("msg", string(msg)))
			continue
		}

		if message.MsgID == "" {
			continue
		}
		if message.ChatID == "" {
			continue
		}

		if message.SeqID == 0 && slicex.ContainInSlice([]int{model.MessageTypeSingle, model.MessageTypeGroup}, message.MessageType) {
			err = message.SetSeqID(message.ChatID)
			if err != nil {
				m := model.NewTipMessage(constant.CodeMessageSendError)
				_ = c.conn.WriteMessage(websocket.TextMessage, m.Marshal())
				logs.Error(ginx.Context, "ConnRead SetSeqID error:", zap.Error(err))
				continue
			}
		}

		message.SendTime = timex.GetBeijingTime()

		WriteTimeout(mediator.broadcast, message.Marshal(), WriteChannelTime)
	}
}

func (c *Client) Write() error {
	heartbeatTW.AddTask(HeartBeatInterval, -1, c.user.ID, c)
	defer heartbeatTW.RemoveTimer(c.user.ID)

	for {
		select {
		case msg := <-c.send:
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logs.Error(ginx.Context, "ConnWrite error:", zap.Error(err))
				c.cleanup()
				return err
			}
		case <-c.done:
			return nil
		}
	}
}
