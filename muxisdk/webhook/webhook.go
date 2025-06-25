package webhook

import "sync"

type Webhook struct {
	Targets sync.Map
}

func NewWebhook() *Webhook {
	return &Webhook{
		Targets: sync.Map{},
	}
}
func (w *Webhook) Register(f uint, t string) {
	li := w.Get(f)
	if li == nil {
		var l = []string{t}
		w.Targets.Store(f, l)
	}
	li = append(li, t)
	w.Targets.Store(f, li)
}
func (w *Webhook) Remove(f uint) {
	w.Targets.Delete(f)
}

func (w *Webhook) RemoveAll() {
	w.Targets = sync.Map{}
}
func (w *Webhook) Get(f uint) []string {
	if catch, ok := w.Targets.Load(f); ok {
		if res, ok := catch.([]string); ok {
			return res
		}
	}
	return nil
}
