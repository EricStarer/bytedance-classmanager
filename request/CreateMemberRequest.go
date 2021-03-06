package request

import "mergeVersion1/types"

/**
创建成员
参数不合法返回 ParamInvalid
只有管理员才能添加
*/
type CreateMemberRequest struct {
	Nickname string         // required，不小于 4 位 不超过 20 位
	Username string         // required，只支持大小写，长度不小于 8 位 不超过 20 位
	Password string         // required，同时包括大小写、数字，长度不少于 8 位 不超过 20 位
	UserType types.UserType // required, 枚举值
}
