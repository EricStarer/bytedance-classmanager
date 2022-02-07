package request
/**
	排课求解器，使老师绑定课程的最优解， 老师有且只能绑定一个课程
 	Method： Post
 */
type ScheduleCourseRequest struct {
	TeacherCourseRelationShip map[string][]string // key 为 teacherID , val 为老师期望绑定的课程 courseID 数组
}
