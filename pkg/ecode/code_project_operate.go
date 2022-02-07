package ecode

var (
	ProjectStopError = New(13001001) // 项目状态不可终止

	UndoneMilepostTask = New(13002001) // 存在未完成的里程碑关键结果
)

var ProjectOperateCodeMsg = map[int]string{
	ProjectStopError.Code():   "项目状态不可终止",
	UndoneMilepostTask.Code(): "存在未完成的里程碑关键结果",
}
