package request
/**
	课程请求
 */
type BookCourseRequest struct {
	StudentID string `json:"student_id"`
	CourseID  string `json:"course_id"`
}
