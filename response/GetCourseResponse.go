package response

import "mergeVersion1/types"

/**
获取课程的返回值
*/
type GetCourseResponse struct {
	Code types.ErrNo
	Data types.TCourse
}
