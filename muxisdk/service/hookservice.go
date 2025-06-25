package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"muxisdk/dao"
	"muxisdk/webhook"
	"net/http"
	"time"
)

//todo 可以尝试模仿gorm的链式调用

type HookService struct {
	Oc *dao.OrmClient
	Wb *webhook.Webhook //这是本质是缓存，真实数据始终存储在数据库中。
}
type HookPayload struct {
	Event string `json:"event"`
	Data  any    `json:"data,omitempty"`
	Try   int    `json:"try"` // 重试次数，最大不超过五
}
type FromHook struct {
	gorm.Model
	From    string       `gorm:"column:from"`
	Targets []HookTarget `gorm:"foreignKey:HookID;references:ID"`
}
type HookTarget struct {
	gorm.Model
	HookID uint   `gorm:"index"`
	Target string `gorm:"index"`
}

//new-->init-->register_from-->register_target-->curd

func NewService(dsn string) *HookService {
	o := dao.NewOrmClient(dsn)
	wb := webhook.NewWebhook()
	return &HookService{
		Oc: o,
		Wb: wb,
	}
}

func (svc *HookService) Init() {
	var f = FromHook{}
	var hook = HookTarget{}
	svc.Oc.InitTable(&hook, &f)
}
func (svc *HookService) RegisterFrom(c context.Context, f string) error {
	var from FromHook
	err := svc.Oc.Get(c, &from, dao.Key{"from", f})
	if err == nil {
		return errors.New("already registered")
	}
	from.From = f
	err = svc.Oc.Create(c, &from)
	if err != nil {
		return err
	}
	return nil
}
func (svc *HookService) GetFromId(c context.Context, f string) (uint, error) {
	var from FromHook
	err := svc.Oc.Get(c, &from, dao.Key{"from", f})
	if err != nil {
		return 0, err
	}
	return from.ID, nil
}
func (svc *HookService) RegisterTar(c context.Context, f uint, to string) error {
	var from FromHook
	err := svc.Oc.Get(c, &from, dao.Key{"id", f})
	if err != nil {
		//防止from未注册
		return err
	}
	// 新增 HookTarget
	target := HookTarget{
		HookID: from.ID,
		Target: to,
	}
	err = svc.Oc.Create(c, &target)
	if err != nil {
		log.Println("Create HookTarget failed:", err)
		return err
	}
	svc.Wb.Register(f, to)
	return nil
}

//方法不支持泛型，这很坏了，为保留性能，不写interface了，强制id查询

func (svc *HookService) GetHook(c context.Context, f uint) ([]string, error) {
	t := svc.Wb.Get(f)
	if t == nil {
		var hook FromHook
		var k = dao.Key{"id", f}
		err := svc.Oc.PreGet(c, hook, []string{"targets"}, k)
		if err != nil {
			return nil, err
		}
		var s []string
		for _, v := range hook.Targets {
			s = append(s, v.Target)
			svc.Wb.Register(f, v.Target)
		}
		return s, nil
	}
	return t, nil
}

func (svc *HookService) Remove(c context.Context, f uint) error {
	var hook FromHook
	var tar HookTarget
	var k1 = dao.Key{"id", f}
	var k2 = dao.Key{"HookId", f}
	err := svc.Oc.Delete(c, &tar, k2)
	if err != nil {
		return err
	}
	err = svc.Oc.Delete(c, &hook, k1)
	if err != nil {
		return err
	}
	svc.Wb.Remove(f)
	return nil
}
func (svc *HookService) Change(c context.Context, f uint, old string, new string) error {
	var hook = HookTarget{
		HookID: f,
		Target: new,
	}
	var k = []dao.Key{{"HookId", f}, {"Target", old}}
	err := svc.Oc.Update(c, &hook, k[0], k[1])
	if err != nil {
		return err
	}
	return nil
}
func (svc *HookService) HookBack(t string, data HookPayload, authorization string) ([]byte, error) {
	if data.Try > 5 {
		return nil, errors.New("too many hooks")
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal hook payload: %w", err)
	}
	var lasterr error
	for i := 0; i < data.Try; i++ {
		reqs, err := http.NewRequest("POST", t, bytes.NewBuffer(jsonBytes))
		if err != nil {
			lasterr = err
			time.Sleep(time.Second)
			continue
		}
		reqs.Header.Set("Content-Type", "application/json")
		if authorization != "" {
			reqs.Header.Set("Authorization", authorization)
		}
		client := &http.Client{}
		resp, err := client.Do(reqs)
		if err != nil {
			lasterr = err
			time.Sleep(time.Second)
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lasterr = readErr
			break
		}
		if resp.StatusCode == http.StatusOK {
			return body, nil
		}
	}

	return nil, lasterr
}
