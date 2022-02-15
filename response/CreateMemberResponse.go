package response

import "mergeVersion1/types"

/**
创建成员的返回
*/
type CreateMemberResponse struct {
	Code types.ErrNo
	Data struct {
		UserID string // int64 范围
	}
}
