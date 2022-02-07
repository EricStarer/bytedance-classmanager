package response

import "Register/types"

/**
	登录成功后需要 Set-Cookie("camp-session", ${value})
	密码错误范围密码错误状态码
 */
type LoginResponse struct {
	Code types.ErrNo
	Data struct {
		UserID string
	}
}
