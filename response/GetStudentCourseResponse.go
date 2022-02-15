package response

import "mergeVersion1/types"

/**
返回学生的排课情况
*/
type GetStudentCourseResponse struct {
	Code types.ErrNo
	Data struct {
		CourseList []types.TCourse
	}
}
