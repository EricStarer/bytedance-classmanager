package response

import "Register/types"

/**
	排课求解器的响应
 */
type ScheduleCourseResponse struct {
	Code types.ErrNo
	Data map[string]string // key 为 teacherID , val 为老师最终绑定的课程 courseID
}

//按照目前意思来看,其好像是说,课程和老师是一一对应的,如果该课程没找到老师,则找一个未分配的老师来教