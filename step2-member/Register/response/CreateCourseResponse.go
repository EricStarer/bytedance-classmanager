package response

import "Register/types"

/**
	创建课程的返回值
 */
type CreateCourseResponse struct {
	Code types.ErrNo
	Data struct {
		CourseID string
	}
}

