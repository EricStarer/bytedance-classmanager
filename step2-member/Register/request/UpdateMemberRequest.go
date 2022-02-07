package request
/**
	更新成员的请求
 */
type UpdateMemberRequest struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
}
