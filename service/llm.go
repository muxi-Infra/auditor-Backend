package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/errorx"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/request"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/config"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/events"
	confluentincllm "github.com/cqhasy/2025-Muxi-Team-auditor-Backend/events/confluentinc-llm"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/client"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/pkg/logger"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/cache"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/repository/dao"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/service/pool"
)

const (
	LLMWorkersNum   = "muxi-auditor-llm-worker-num"
	LLMJobsChanSize = "muxi-auditor-llm-jobs-chan-size"
)

const (
	retry         = 3
	GroupID       = "LLM_1"
	TopicPending  = "PendingToAuditItems"
	TopicFinished = "FinishedAuditItems"
	TimeOut       = 500 * time.Millisecond
)

type AuditState int

const (
	StateAIResult    AuditState = 0
	StateHookSuccess            = 1
	StateHookFail               = 2
)

type LLMService struct {
	log      logger.Logger
	userDAO  dao.UserDAOInterface
	itemDAO  dao.ItemDaoInterface
	proDAO   dao.ProjectDAOInterface
	pcache   cache.ProjectCacheInterface
	client   client.AuditAIClient
	jobsPool *pool.Pool
	pdu      events.Producer
	csu      events.Consumer
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelCauseFunc
}

func NewLLMService(userDAO *dao.UserDAO, itemDAO *dao.ItemDao, proDAO *dao.ProjectDAO,
	c client.AuditAIClient, lo logger.Logger, pc *cache.ProjectCache,
	pd *confluentincllm.LlmProducer, conf *config.KafkaConfig) *LLMService {
	l := LLMService{
		userDAO:  userDAO,
		itemDAO:  itemDAO,
		proDAO:   proDAO,
		log:      lo,
		client:   c,
		jobsPool: pool.NewPool(envInt(LLMWorkersNum, pool.DefaultWorkers), envInt(LLMJobsChanSize, pool.DefaultTaskNum)),
		pdu:      pd,
		csu:      confluentincllm.NewLlmConsumer(conf, GroupID),
		pcache:   pc,
	}

	c.WrapLogger(lo)
	pd.WrapLogger(lo)
	l.csu.WrapLogger(lo)
	l.csu.Subscribe([]string{TopicPending, TopicFinished})

	ctx, cancel := context.WithCancelCause(context.Background())
	l.ctx = ctx
	l.cancel = cancel

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go l.start()
	return &l
}

func (l *LLMService) start() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.worker()
	}()

	for i := 0; i < l.jobsPool.GetWorkerNums(); i++ {
		l.wg.Add(1)

		go func() {
			defer l.wg.Done()
			for v := range l.jobsPool.PendingJobs {
				l.processData(v.Topic, v.Data)
			}
		}()
	}

}

func (l *LLMService) worker() {
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			topic, data := l.csu.Consume()
			if data == nil {
				time.Sleep(50 * time.Millisecond) // 避免过快消费造成空轮询
				continue
			}
			l.jobsPool.Submit(topic, data)
		}
	}
}

func (l *LLMService) sendToLLM(id uint, role string, c response.Contents) {
	td, imd := l.client.Transform(role, c)

	for i := 0; i < retry; i++ {
		resp, err := l.client.SendMessage(td, imd)

		if errors.Is(err, errorx.ErrUnSupportImage) {
			l.log.Error(err.Error())
			break
		}

		if err == nil {
			resp.ID = id

			if resp.Confidence > 50 {
				key := []byte(strconv.FormatInt(int64(resp.ID), 10))

				data, err := json.Marshal(resp)
				if err != nil {
					l.log.Error(err.Error())
					continue
				}

				l.pdu.Produce(TopicFinished, key, data)
			}
			return
		}

		l.log.Warn("retrying sendToLLM", logger.Int("retry", retry), logger.Error(err))
		time.Sleep(time.Second * time.Duration(retry+1))
	}
	l.log.Error("sendToLLM failed after retries", logger.Int("ItemID", int(id)))
}

func (l *LLMService) Audit(Data []request.AuditItem) {
	for _, item := range Data {
		data, err := json.Marshal(item)
		if err != nil {
			l.log.Error(err.Error())
			continue
		}

		key := []byte(strconv.FormatInt(int64(item.ID), 10))

		l.pdu.Produce(TopicPending, key, data)
	}
}

func (l *LLMService) tryHook(result model.AuditResult) bool {
	for i := 0; i < retry; i++ {
		item, err := l.itemDAO.FindItemByID(context.Background(), result.ID)
		if err != nil {
			l.log.Error(err.Error())
			continue
		}

		data := request.WebHookData{
			Id:     item.HookId,
			Status: auditStatusForHook(result.Result),
			Msg:    item.Reason,
		}

		_, err = hookBack(item.HookUrl, request.HookPayload{
			Event: "audit result back",
			Data:  data,
			Try:   retry,
		}, "")

		if err != nil {
			l.log.Error(err.Error())
			continue
		}

		return true
	}
	return false
}

func (l *LLMService) persistState(result model.AuditResult, state AuditState) error {
	switch state {

	case StateHookSuccess:
		return l.itemDAO.AuditItem(result.ID, model.Pass, result.Reason)

	case StateHookFail:
		return l.itemDAO.AuditItem(result.ID, model.Pending, result.Reason)
	default:
		l.log.Error("意外的state状态")
	}

	return nil
}

func (l *LLMService) processData(topic string, data []byte) {
	switch topic {
	case TopicPending:
		{
			var item request.AuditItem

			err := json.Unmarshal(data, &item)
			if err != nil {
				l.log.Error("解析Item失败", logger.Error(err))
			}

			role, err := l.pcache.GetAuditRole(context.Background(), item.ProjectID)
			if err != nil {
				if errors.Is(err, redis.Nil) {
					role, er := l.proDAO.GetProjectRole(context.Background(), item.ProjectID)
					if er != nil {
						l.log.Error("获取project_role失败", logger.Error(er))
						return
					}

					go func(pid uint, role string) {
						if er = l.pcache.SetAuditRole(context.Background(), pid, role); er != nil {
							l.log.Error(er.Error())
						}
					}(item.ProjectID, role)
				} else {
					l.log.Error(err.Error())
				}
			}

			l.sendToLLM(item.ID, role, item.Contents)
		}

	case TopicFinished:
		{
			state := StateAIResult
			var result model.AuditResult
			err := json.Unmarshal(data, &result)
			if err != nil {
				l.log.Error(err.Error())
			}

			ok := l.tryHook(result)
			if ok {
				state = StateHookSuccess
			} else {
				state = StateHookFail
			}

			err = l.persistState(result, state)
			if err != nil {
				l.log.Error("failed to persist audit state", logger.Error(err))
			}
		}

	default:
		l.log.Error("unexpected topic", logger.String("topic", topic))
	}
}

func (l *LLMService) Close() {
	l.cancel(nil)
	close(l.jobsPool.PendingJobs)
	l.wg.Wait()

	l.pdu.Close(TimeOut)
	l.csu.Close()
}
