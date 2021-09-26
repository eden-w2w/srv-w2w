package enums

//go:generate eden generate enum --type-name=StatementStatus
// api:enum
type StatementStatus uint8

// 结算状态
const (
	STATEMENT_STATUS_UNKNOWN   StatementStatus = iota
	STATEMENT_STATUS__CREATED                  // 待结算
	STATEMENT_STATUS__COMPLETE                 // 已结算
)
