package ali

import (
	"github.com/alibabacloud-go/tea/tea"
	"strconv"
	"strings"

	gre "github.com/alibabacloud-go/green-20220302/v2/client"

	"github.com/muxi-Infra/auditor-Backend/api/response"
	"github.com/muxi-Infra/auditor-Backend/langchain/model"
	"github.com/muxi-Infra/auditor-Backend/pkg/stringx"
)

const (
	ContentTag     = "content: "
	LastCommentTag = "评论1: "
	NextCommentTag = "评论2: "
)

/*
	写这里时突然感觉不太好办，我们记录评论本意是作为上下文辅助理解帖子的，但这种审核模型似乎会因为评论的违规而将帖子

驳回，感觉并不合理；图片的审核ollama做不了这里可以，但对于整体内容的审核感觉还是语义大模型比较好，后续可以改为用这个
做一层基本层次的审核，在每个部分都合法合规的前提下，通过ollama之类的语义模型去做宏观的判断，比如引战之类的，
*/
func parseContent(content response.Contents) (string, []string) {
	var texts []string
	var pictures []string

	if content.Topic.Title != "" {
		texts = append(texts, content.Topic.Title)
	}
	if content.Topic.Content != "" {
		texts = append(texts, ContentTag)
		texts = append(texts, content.Topic.Content)
	}
	pictures = append(pictures, content.Topic.Pictures...)

	if content.LastComment.Content != "" {
		texts = append(texts, LastCommentTag)
		texts = append(texts, content.LastComment.Content)
	}
	pictures = append(pictures, content.LastComment.Pictures...)

	if content.NextComment.Content != "" {
		texts = append(texts, NextCommentTag)
		texts = append(texts, content.NextComment.Content)
	}
	pictures = append(pictures, content.NextComment.Pictures...)

	// 拼接所有文字为一条字符串（可用于一次性送审）
	fullText := ""
	if len(texts) > 0 {
		fullText = strings.Join(texts, "\n") // 每条文字换行分隔
	}

	return fullText, pictures
}

func parseImageResponse(i []*gre.ImageModerationResponseBodyData) []model.ImageAuditResult {
	var results []model.ImageAuditResult
	for _, img := range i {
		re := model.ImageAuditResult{
			ID:         tea.StringValue(img.DataId),
			Result:     model.PassBeforeHook,
			Confidence: 100,
		}

		for _, v := range img.Result {
			if v.Confidence != nil {
				c := tea.Float32Value(v.Confidence)
				re.Reason = append(re.Reason, stringx.Build("Label:", tea.StringValue(v.Label), "--",
					"confidence:", strconv.FormatFloat(float64(c), 'f', 2, 32)))

				if c < 60 {
					continue
				}
				re.Result = model.RejectBeforeHook
				if c < re.Confidence {
					re.Confidence = c
				}
			}
		}

		results = append(results, re)
	}

	return results
}

func parseTextResponse(t *gre.TextModerationPlusResponseBodyData) *model.TextAuditResult {
	var result = model.TextAuditResult{
		Result:     model.PassBeforeHook,
		Confidence: 100,
	}
	for _, v := range t.Result {
		if v.Confidence != nil {
			c := tea.Float32Value(v.Confidence)
			result.Reason = append(result.Reason, stringx.Build("Label:", tea.StringValue(v.Label), "--",
				"riskWords:", tea.StringValue(v.RiskWords), "--",
				"confidence:", strconv.FormatFloat(float64(c), 'f', 2, 32)))

			if c < 60 {
				continue
			}
			result.Result = model.RejectBeforeHook
			if tea.Float32Value(v.Confidence) < result.Confidence {
				result.Confidence = tea.Float32Value(v.Confidence)
			}
		}
	}
	return &result
}

func transformPics(pics []string) []model.ImageParameters {
	var results []model.ImageParameters
	for _, pic := range pics {
		results = append(results, model.ImageParameters{
			DataId:   pic,
			ImageUrl: pic,
		})
	}
	return results
}

func merge(t *model.TextAuditResult, imgs []model.ImageAuditResult) model.AuditResult {
	if t == nil {
		return model.AuditResult{}
	}
	var result = model.AuditResult{
		Result:     t.Result,
		Confidence: t.Confidence,
	}
	b := stringx.Acquire()
	defer stringx.Release(b)
	for _, v := range t.Reason {
		b.WriteString(v)
		b.WriteByte('\n')
	}

	if len(imgs) > 0 {
		for _, img := range imgs {
			// pass是3而reject是4
			if img.Result == model.RejectBeforeHook {
				if img.Result > result.Result {
					result.Result = img.Result
					if result.Confidence > img.Confidence {
						result.Confidence = img.Confidence
					}
				}

				if result.Confidence < img.Confidence {
					result.Confidence = img.Confidence
				}
				// 图片审核通过的理由不做记录，减小reason长度
				for _, v := range img.Reason {
					b.WriteString(v)
					b.WriteByte('\n')
				}
			}
		}
	}

	result.Reason = b.String()
	return result
}
