package request
/**
	登陆请求
 */
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
