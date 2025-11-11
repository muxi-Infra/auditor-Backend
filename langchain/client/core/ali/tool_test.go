package ali

import (
	"fmt"
	gre "github.com/alibabacloud-go/green-20220302/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/langchain/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	var contents = response.Contents{
		Topic: response.Topics{
			Title:    "hello china",
			Content:  "中国真好啊",
			Pictures: []string{"111", "222"},
		},
		LastComment: response.Comment{
			Content:  "你逗我呢",
			Pictures: []string{"333", "444"},
		},
		NextComment: response.Comment{
			Content:  "乐子别叫",
			Pictures: []string{"555", "666"},
		},
	}
	text, pic := parseContent(contents)
	assert.Equal(t, text, "hello china\ncontent: \n中国真好啊\n评论1: \n你逗我呢\n评论2: \n乐子别叫")
	assert.Equal(t, pic, []string{"111", "222", "333", "444", "555", "666"})
	fmt.Println(text, pic)
}

func TestParseImageResponse(t *testing.T) {
	data := []*gre.ImageModerationResponseBodyData{
		{
			DataId: tea.String("2"),
			Result: []*gre.ImageModerationResponseBodyDataResult{
				{
					Description: tea.String("未检测出风险"),
					Label:       tea.String("nonLabel"),
					RiskLevel:   tea.String("none"),
				},
			},
			RiskLevel: tea.String("none"),
		},
		{
			DataId: tea.String("1"),
			Result: []*gre.ImageModerationResponseBodyDataResult{
				{
					Confidence:  tea.Float32(94.92),
					Description: tea.String("成人色情"),
					Label:       tea.String("pornographic_adultContent"),
					RiskLevel:   tea.String("high"),
				},
			},
			RiskLevel: tea.String("high"),
		},
	}
	re := parseImageResponse(data)
	expected := []model.ImageAuditResult{
		{
			ID:         "2",
			Result:     model.PassBeforeHook,
			Confidence: 100,
		},
		{
			ID:         "1",
			Result:     model.RejectBeforeHook,
			Confidence: 94.92,
			Reason:     []string{"Label:pornographic_adultContent--confidence:94.92"},
		},
	}
	assert.Equal(t, expected, re)
}

func TestParseTextResponse(t *testing.T) {
	tar1 := &gre.TextModerationPlusResponseBodyData{
		Result: []*gre.TextModerationPlusResponseBodyDataResult{
			{
				Confidence:  tea.Float32(100),
				Description: tea.String("疑似敏感政治内容"),
				Label:       tea.String("political_n"),
				RiskWords:   tea.String("推翻政府"),
			},
		},
		RiskLevel: tea.String("high"),
	}
	tar2 := &gre.TextModerationPlusResponseBodyData{
		Result: []*gre.TextModerationPlusResponseBodyDataResult{
			{
				Description: tea.String("未检测出风险"),
				Label:       tea.String("nonLabel"),
			},
		},
		RiskLevel: tea.String("none"),
	}
	re1 := parseTextResponse(tar1)
	re2 := parseTextResponse(tar2)
	assert.Equal(t, &model.TextAuditResult{
		Result:     model.RejectBeforeHook,
		Confidence: 100.0,
		Reason:     []string{"Label:political_n--riskWords:推翻政府--confidence:100.00"},
	}, re1)
	assert.Equal(t, &model.TextAuditResult{
		Result:     model.PassBeforeHook,
		Confidence: 100,
	}, re2)
}

func TestMerge(t *testing.T) {
	tr := &model.TextAuditResult{
		Result:     model.RejectBeforeHook,
		Confidence: 95.0,
		Reason:     []string{"Label:political_n--riskWords:推翻政府--confidence:100.00"},
	}

	ir := []model.ImageAuditResult{
		{
			ID:         "2",
			Result:     model.PassBeforeHook,
			Confidence: 100,
		},
		{
			ID:         "1",
			Result:     model.RejectBeforeHook,
			Confidence: 94.92,
			Reason:     []string{"Label:pornographic_adultContent--confidence:94.92"},
		},
	}
	re := merge(tr, ir)
	expected := model.AuditResult{
		Result:     model.RejectBeforeHook,
		Confidence: 95.0,
		Reason:     "Label:political_n--riskWords:推翻政府--confidence:100.00\nLabel:pornographic_adultContent--confidence:94.92\n",
	}
	assert.Equal(t, expected, re)
}
