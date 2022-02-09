package response

import "Register/types"

/**
	批量获取成员信息的返回
 */
type GetMemberListResponse struct {
	Code types.ErrNo
	Data struct {
		MemberList []types.TMember
	}
}
