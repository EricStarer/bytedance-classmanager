package response

import "Register/types"

/**
	登出成功需要删除 Cookie
 */

type LogoutResponse struct {
	Code types.ErrNo
}

