package sqlite

import "errors"

var (
	// ErrEmptyCond 防止在update、delete操作时，漏掉条件造成的严重后果
	// 如果确实不需要条件，请将条件设置为 1=1
	ErrEmptyCond               = errors.New("条件不能为空")
	ErrEmptyValue              = errors.New("值不能为空")
	ErrInvalidArguments        = errors.New("参数错误")
	ErrInvalidCondArguments    = errors.New("条件参数错误")
	ErrInvalidColumnsArguments = errors.New("字段名类型错误")
	ErrInvalidResultPtr        = errors.New("无效的结果指针")
	ErrEmptyTableName          = errors.New("没有指定表名称")
	ErrDatabaseNotExists       = errors.New("指定数据库不存在，或未初始化")
)
