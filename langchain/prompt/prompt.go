package prompt

import "fmt"

func buildPrompt(auditRole, content string) string {
	return fmt.Sprintf(`
	你是一个专业的内容审核系统。请根据以下要求审核如下内容：

	【最基本审核维度】(如有违法直接拒绝)
	1. 违法犯罪：恐怖主义、暴力、毒品、诈骗等
	2. 色情低俗：色情、露骨性暗示
	3. 辱骂歧视：人身攻击、辱骂、歧视性言论
	4. 广告骚扰：垃圾广告、营销内容
	5. 其他违规：违反社区规则的其他内容

	【特定审核规则】
	%s

	【输出要求】
	- 只返回 JSON 格式
	- 字段包括：
	  - result: "pass" | "review" | "reject"
	  - reasons: 数组，违规类别
	  - confidence: 0~1 之间的置信度
	  - suggestion: "通过" | "人工复核" | "删除或屏蔽"
	
	【待审核内容】
	%s
	`, auditRole, content)
}
