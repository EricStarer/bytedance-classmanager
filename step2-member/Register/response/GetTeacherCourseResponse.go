package response

import "Register/types"

/**
	获取老师下所有课程的返回值
 */
type GetTeacherCourseResponse struct {
	Code types.ErrNo
	Data struct {
		CourseList []*types.TCourse
	}
}