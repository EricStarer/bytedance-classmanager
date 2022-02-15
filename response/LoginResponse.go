package response

import "mergeVersion1/types"

/**
登录成功后需要 Set-Cookie("camp-session", ${value})
密码错误返回密码错误状态码
*/
type LoginResponse struct {
	Code types.ErrNo
	Data struct {
		UserID string
	}
}
