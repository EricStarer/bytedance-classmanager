package response

import "mergeVersion1/types"

/**
获取成员的返回结构
如果用户已删除请返回已删除状态码，不存在请返回不存在状态码
*/
type GetMemberResponse struct {
	Code types.ErrNo
	Data types.TMember
}
