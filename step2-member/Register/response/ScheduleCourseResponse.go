package response

import "Register/types"

/**
	排课求解器的响应
 */
type ScheduleCourseResponse struct {
	Code types.ErrNo
	Data map[string]string // key 为 teacherID , val 为老师最终绑定的课程 courseID
}
