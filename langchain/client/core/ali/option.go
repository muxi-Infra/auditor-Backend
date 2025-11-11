package ali

import (
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func WithConnetTimeLimit(t int) AlClientOpt {
	return func(c *AlClient) {
		r := &util.RuntimeOptions{
			ConnectTimeout: tea.Int(t),
		}
		c.runtime = r
	}
}

func WithReadTimeLimit(t int) AlClientOpt {
	return func(c *AlClient) {
		r := &util.RuntimeOptions{
			ReadTimeout: tea.Int(t),
		}
		c.runtime = r
	}
}

func WithTimeLimit(t int) AlClientOpt {
	return func(c *AlClient) {
		r := &util.RuntimeOptions{
			ReadTimeout:    tea.Int(t),
			ConnectTimeout: tea.Int(t),
		}
		c.runtime = r
	}
}
