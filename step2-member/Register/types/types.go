package types

// 说明：
// 1. 所提到的「位数」均以字节长度为准
// 2. 所有的 ID 均为 int64（以 string 方式表现）

// 通用结构
type Test struct {
	NickName string `json:"nick_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserType string `json:"type_user"`
}

type ErrNo int

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

type TMember struct {
	UserID   string `gorm:"uniqueIndex"`
	Nickname string `gorm:"column:nick_name"`
	Username string	`gorm:"column:user_name"`
	UserType UserType
}

type TAdmin struct {
	ID uint64 `gorm:"primaryKey"`
	TMember TMember `gorm:"embedded"`
	IsDel int `gorm:"default:0"`
	Password string
}

type TStudent struct {
	ID uint64 `gorm:"primaryKey"`
	TMember TMember `gorm:"embedded"`
	IsDel int  `gorm:"default:0"`
	Password string
	CourseRecordId string
}

type TTeacher struct {
	ID uint64 `gorm:"primaryKey"`
	TMember TMember `gorm:"embedded"`
	IsDel int `gorm:"default:0"`
	Password string
	TeachRecordId string
}

type GenerateId struct {
	ID uint64 `gorm:"primaryKey"`
	UserType UserType
	IsDel int `gorm:"default:0"`
}

func (GenerateId) TableName() string {
	return "generate_id"
}

type ResponseMeta struct {
	Code ErrNo
}



type TCourse struct {
	CourseID  string
	Name      string
	TeacherID string
}

// -----------------------------------

// 成员管理

type UserType int

const (
	Admin   UserType = 1
	Student UserType = 2
	Teacher UserType = 3
)

// 系统内置管理员账号
// 账号名：JudgeAdmin 密码：JudgePassword2022
