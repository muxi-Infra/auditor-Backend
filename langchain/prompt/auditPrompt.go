package prompt

import (
	"encoding/json"
	"fmt"
	"github.com/cqhasy/2025-Muxi-Team-auditor-Backend/api/response"
)

func BuildPrompt(auditRole string, contents response.Contents) string {
	contentJSON, _ := json.MarshalIndent(contents, "", "  ")

	prompt := fmt.Sprintf(`
你是一个专业内容审核系统。请严格根据以下规则审核内容，并**只返回 JSON**，不要添加任何解释或其他文本。
【最基本审核规则】(若违反，直接拒绝)
1. 违法犯罪：恐怖主义、暴力、毒品、诈骗等
2. 色情低俗：色情、露骨性暗示
3. 辱骂歧视：人身攻击、辱骂、歧视性言论
4. 广告骚扰：垃圾广告、营销内容
5. 其他违规：违反社区规则的其他内容

【特定审核规则】
%s

【输出要求】
请严格返回 JSON，包含以下字段：
{
  "result": "string", 
  "reason": "string",                     
  "confidence": "number"                  
}
【字段说明】
result: 没有违反基本规则和特定审核规则就pass(4), 不确定返回review(0), 认为不合适返回reject(5)
reason: 简短说明拒绝或审核原因
confidence: 0~1 之间，表示审核结果的可信度
【示例】
{
  "result": 5,
  "reason": "傻逼是辱骂用语",
  "confidence": 0.95
}
{
  "result": 4,
  "reason": "你好是正常对话",
  "confidence": 0.95
}

【待审核内容】
%s

严格遵守 JSON 格式，不要输出任何额外文本。
`, auditRole, string(contentJSON))

	return prompt
}
