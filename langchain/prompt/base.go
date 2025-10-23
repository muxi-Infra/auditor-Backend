package prompt

import "strings"

func GetToolsPrompt(task string, tools string) string {
	prompt := `
你是一个智能助手，需要完成我指定的任务：
请按照以下规则操作：
1. 只返回你认为需要调用的工具列表，不要执行工具。
2. 请根据任务生成 JSON 输出，严格按照以下格式：
{
  "data": [ /* 需要调用的工具名称数组，例如 ["TextAnalyzer","SentimentChecker"] ,只返回必须调用的工具，对此次任务没有帮助的工具不要返回*/ ],
  "tokens": "你本次任务消耗的token数，如果是本地模型不需要消耗的返回0，如果不知道返回-1",
  "error": "" // 如无法确定工具列表，填 "DesideToolListErr;如果此次任务无需外部工具或者外部工具为空，请跟就任务要求和内容把data字段赋值为对应结果"
}
3.只返回 JSON，不要添加任何解释、<think> 标签或者额外文本，不要任何注释或 Markdown 标记.
本次任务: {{TASK}}
你可以使用以下工具：
{{TOOLS}}`
	replacer := strings.NewReplacer(
		"{{TASK}}", string(task),
		"{{TOOLS}}", tools,
	)

	return replacer.Replace(prompt)
}

func ExecPrompt(task string, tools string) string {

	template := `
只返回一个标准 JSON，不能有其他说明文字。
我是一个外部服务，我会给你布置任务并提供相应工具，
你是一个智能助手，需要完成以下任务：
任务: {{TASK}}

你已经得到了工具调用结果(tool_results):
{{TOOL_RESULTS_JSON}}

请基于这些工具结果完成任务，并严格按照以下 JSON 格式返回：
{
  "data": { /* 基于工具结果的最终任务输出 */ },
  "tokens": "本次任务消耗的 token 数",
  "error": "如执行任务出错填 AISelfErr，否则留空"
}

注意：
- 只返回 JSON，不要返回其他文字。
`
	replacer := strings.NewReplacer(
		"{{INPUT}}", string(task),
		"{{TOOLS}}", tools,
	)

	return replacer.Replace(template)
}
