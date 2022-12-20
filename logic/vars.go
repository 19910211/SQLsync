package logic

const (
	_datetimeLayout = "20060102150405" // 年月日 时分秒
	_yearsLayout    = "2006"           // 年
	_monthLayout    = "01"             // 月
	_timeLayout     = "150405"         // 时分秒
)

const (
	NewVersion  = "new"  // 新增sql版本
	SyncVersion = "sync" // 同步sql版本
)

var EmptyVersion = Version{} // 空版本
