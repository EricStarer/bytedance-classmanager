package response

import "mergeVersion1/types"

/**
登出成功需要删除 Cookie
*/

type LogoutResponse struct {
	Code types.ErrNo
}
