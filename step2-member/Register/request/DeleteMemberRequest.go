package request

/**
删除成员的请求
成员删除后，该成员不能够被登录且不应该不可见，ID 不可复用
*/
type DeleteMemberRequest struct {
	UserID string
}
