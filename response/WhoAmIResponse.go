package response

import "mergeVersion1/types"

/**
用户未登录请返回用户未登录状态码
*/

type WhoAmIResponse struct {
	Code types.ErrNo
	Data types.TMember
}
