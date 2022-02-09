package response

import "Register/types"

/**
	课程已满返回 CourseNotAvailable
 */

type BookCourseResponse struct {
	Code types.ErrNo
}
