package request

/**
老师绑定课程请求
 Method： Post
 注：这里的 teacherID 不需要做已落库校验
 一个老师可以绑定多个课程 , 不过，一个课程只能绑定在一个老师下面
*/

type BindCourseRequest struct {
	CourseID  string
	TeacherID string
}
