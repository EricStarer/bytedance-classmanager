package request
/**
	老师解绑课程
 	Method： Post
 */

type UnbindCourseRequest struct {
	CourseID  string
	TeacherID string
}
