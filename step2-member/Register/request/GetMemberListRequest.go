package request

/**
	批量获取成员信息
 */
type GetMemberListRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
