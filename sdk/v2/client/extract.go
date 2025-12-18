package client

import (
	"encoding/json"
	"fmt"

	"github.com/muxi-Infra/auditor-Backend/sdk/v2/api/errorx"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/api/response"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/dto"
)

func extractInt(row response.Resp) (int, error) {
	v, ok := row.Data.(float64)
	if !ok {
		return 0, errorx.SeverDataIllegal(&errorx.TypeMismatchCause{
			ActualType:   fmt.Sprintf("%T", row.Data),
			ExpectedType: "float64",
			ActualValue:  row.Data,
		}, row.Code)
	}

	return int(v), nil
}

func extractItemsStatus(row response.Resp) (dto.ItemsStatus, error) {
	b, err := marshalData(row.Data)
	if err != nil {
		return dto.ItemsStatus{}, errorx.MarshalErr(err)
	}

	var result dto.ItemsStatus
	if err := json.Unmarshal(b, &result); err != nil {
		return dto.ItemsStatus{}, errorx.SeverDataIllegal(&errorx.TypeMismatchCause{
			ActualType:   fmt.Sprintf("%T", row.Data),
			ExpectedType: "dto.Contents",
			ActualValue:  row.Data,
		}, row.Code)
	}

	return result, nil
}

func marshalData(data any) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Data: %w", err)
	}

	return b, nil
}
