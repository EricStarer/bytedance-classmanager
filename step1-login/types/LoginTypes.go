package types

//状态码
type ErrNo int

//状态码枚举
const (
	OK                 ErrNo = 0
	ParamInvalid       ErrNo = 1  // 参数不合法
	UserHasExisted     ErrNo = 2  // 该 Username 已存在
	UserHasDeleted     ErrNo = 3  // 用户已删除
	UserNotExisted     ErrNo = 4  // 用户不存在
	WrongPassword      ErrNo = 5  // 密码错误
	LoginRequired      ErrNo = 6  // 用户未登录
	CourseNotAvailable ErrNo = 7  // 课程已满
	CourseHasBound     ErrNo = 8  // 课程已绑定过
	CourseNotBind      ErrNo = 9  // 课程未绑定过
	PermDenied         ErrNo = 10 // 没有操作权限
	StudentNotExisted  ErrNo = 11 // 学生不存在
	CourseNotExisted   ErrNo = 12 // 课程不存在
	StudentHasNoCourse ErrNo = 13 // 学生没有课程
	StudentHasCourse   ErrNo = 14 // 学生有课程

	UnknownError ErrNo = 255 // 未知错误
)

//通用响应状态码
type ResponseMeta struct {
	Code ErrNo
}

//用户类型
type UserType int

//用户类型枚举
const (
	Admin   UserType = 1
	Student UserType = 2
	Teacher UserType = 3
)

// 系统内置管理员账号
// 账号名：JudgeAdmin 密码：JudgePassword2022

//用户信息
type TMember struct {
	UserID   string   `db:"user_id"`
	Nickname string   `db:"nick_name"`
	Username string   `db:"user_name"`
	UserType UserType `db:"user_type"`
}

// ----------------------------------------
// 登录
//登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录成功后需要 Set-Cookie("camp-session", ${value})
// 密码错误范围密码错误状态码
//登录响应
type LoginResponse struct {
	Code ErrNo
	Data struct {
		UserID string
	}
}

// 登出
//登出请求
type LogoutRequest struct{}

// 登出成功需要删除 Cookie
//登出响应
type LogoutResponse struct {
	Code ErrNo
}

// WhoAmI 接口，用来测试是否登录成功，只有此接口需要带上 Cookie
//whoami请求
type WhoAmIRequest struct {
}

// 用户未登录请返回用户未登录状态码
//whoami响应
type WhoAmIResponse struct {
	Code ErrNo
	Data TMember
}
