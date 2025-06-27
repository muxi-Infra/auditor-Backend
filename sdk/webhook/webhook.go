package webhook

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/gin-gonic/gin"

	"net/http"
)

// HandlerFunc 用户定义的回调处理函数
type HandlerFunc func(event string, data request.HookPayload)

type Listener struct {
	Engine  *gin.Engine
	Addr    string
	Path    string
	Handler HandlerFunc
}

func NewListener(engine *gin.Engine, addr string, path string, handler HandlerFunc) *Listener {
	l := &Listener{
		Engine:  engine,
		Addr:    addr,
		Path:    path,
		Handler: handler,
	}
	return l
}

func (l *Listener) RegisterRoutes() {
	l.Engine.POST(l.Path, func(c *gin.Context) {
		var payload request.HookPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Code: 400,
				Msg:  "invalid payload",
				Data: nil,
			})
			return
		}

		l.Handler(payload.Event, payload)

		c.JSON(http.StatusOK, response.Response{
			Code: 200,
			Msg:  "success",
			Data: payload,
		})
	})
}

// 启动监听器
func (l *Listener) Start() error {
	return l.Engine.Run(l.Addr)
}
func (l *Listener) RegisterRouteWithKa(kafkaProducer sarama.SyncProducer, topic string) {
	l.Engine.POST(l.Path, func(c *gin.Context) {
		var payload request.HookPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Code: 400,
				Msg:  "invalid payload",
				Data: nil,
			})
			return
		}
		// 序列化为 JSON
		bytes, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Response{
				Code: 500,
				Msg:  "序列化失败",
				Data: err,
			})
			return
		}
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(bytes),
		}

		// 发送到 Kafka
		_, _, err = kafkaProducer.SendMessage(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Response{
				Code: 500,
				Msg:  "fail to send message to kafka",
				Data: err,
			})
			return
		}

		c.JSON(http.StatusOK, response.Response{
			Code: 200,
			Msg:  "success to send message to kafka",
		})
	})
}
