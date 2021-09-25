package enums

//go:generate eden generate enum --type-name=StatementType
// api:enum
type StatementType uint8

// 结算周期类型
const (
	STATEMENT_TYPE_UNKNOWN StatementType = iota
	STATEMENT_TYPE__WEEK                 // 周
	STATEMENT_TYPE__MONTH                // 月
)
