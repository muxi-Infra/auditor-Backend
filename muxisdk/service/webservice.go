package service

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"muxisdk/model"
	"net/http"
)

// HandlerFunc 用户定义的回调处理函数
type HandlerFunc func(event string, data interface{})

type Listener struct {
	engine  *gin.Engine
	addr    string
	path    string
	handler HandlerFunc
}

func NewListener(engine *gin.Engine, addr string, path string, handler HandlerFunc) *Listener {
	l := &Listener{
		engine:  engine,
		addr:    addr,
		path:    path,
		handler: handler,
	}
	l.registerRoutes()
	return l
}

func (l *Listener) registerRoutes() {
	l.engine.POST(l.path, func(c *gin.Context) {
		var payload HookPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, model.Response{
				Code: 400,
				Msg:  "invalid payload",
				Data: nil,
			})
			return
		}

		l.handler(payload.Event, payload.Data)

		c.JSON(http.StatusOK, model.Response{
			Code: 200,
			Msg:  "success",
			Data: payload,
		})
	})
}

// 启动监听器
func (l *Listener) Start() error {
	return l.engine.Run(l.addr)
}
func (l *Listener) RegisterRouteWithKa(kafkaProducer sarama.SyncProducer, topic string) {
	l.engine.POST(l.path, func(c *gin.Context) {
		var payload HookPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, model.Response{
				Code: 400,
				Msg:  "invalid payload",
				Data: nil,
			})
			return
		}
		// 序列化为 JSON
		bytes, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{
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
			c.JSON(http.StatusInternalServerError, model.Response{
				Code: 500,
				Msg:  "fail to send message to kafka",
				Data: err,
			})
			return
		}

		c.JSON(http.StatusOK, model.Response{
			Code: 200,
			Msg:  "success to send message to kafka",
		})
	})
}
