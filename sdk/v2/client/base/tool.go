package base

import (
	"strconv"
	"strings"
)

// 工具函数：将 []uint 转换为逗号分隔的字符串
func uintSliceToString(ids []int) string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = strconv.FormatUint(uint64(id), 10)
	}
	return strings.Join(strs, ",")
}
