package request
/**
	创建课程请求
	Method: Post
 */

type CreateCourseRequest struct {
	Name string `json:"name"`
	Cap  int `json:"cap"`
}
