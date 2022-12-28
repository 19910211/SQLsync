package logic

const (
	DatetimeLayout = "20060102150405" // 年月日 时分秒
	YearsLayout    = "2006"           // 年
	MonthLayout    = "01"             // 月
	TimeLayout     = "150405"         // 时分秒
)

const (
	NewVersion  = "new"  // 新增sql版本
	SyncVersion = "sync" // 同步sql版本
)

var EmptyVersion = Version{} // 空版本
