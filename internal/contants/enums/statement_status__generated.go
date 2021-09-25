package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_eden_framework_enumeration "github.com/eden-framework/enumeration"
)

var InvalidStatementStatus = errors.New("invalid StatementStatus")

func init() {
	github_com_eden_framework_enumeration.RegisterEnums("StatementStatus", map[string]string{
		"COMPLETE": "已结算",
		"CREATED":  "未结算",
	})
}

func ParseStatementStatusFromString(s string) (StatementStatus, error) {
	switch s {
	case "":
		return STATEMENT_STATUS_UNKNOWN, nil
	case "COMPLETE":
		return STATEMENT_STATUS__COMPLETE, nil
	case "CREATED":
		return STATEMENT_STATUS__CREATED, nil
	}
	return STATEMENT_STATUS_UNKNOWN, InvalidStatementStatus
}

func ParseStatementStatusFromLabelString(s string) (StatementStatus, error) {
	switch s {
	case "":
		return STATEMENT_STATUS_UNKNOWN, nil
	case "已结算":
		return STATEMENT_STATUS__COMPLETE, nil
	case "未结算":
		return STATEMENT_STATUS__CREATED, nil
	}
	return STATEMENT_STATUS_UNKNOWN, InvalidStatementStatus
}

func (StatementStatus) EnumType() string {
	return "StatementStatus"
}

func (StatementStatus) Enums() map[int][]string {
	return map[int][]string{
		int(STATEMENT_STATUS__COMPLETE): {"COMPLETE", "已结算"},
		int(STATEMENT_STATUS__CREATED):  {"CREATED", "未结算"},
	}
}

func (v StatementStatus) String() string {
	switch v {
	case STATEMENT_STATUS_UNKNOWN:
		return ""
	case STATEMENT_STATUS__COMPLETE:
		return "COMPLETE"
	case STATEMENT_STATUS__CREATED:
		return "CREATED"
	}
	return "UNKNOWN"
}

func (v StatementStatus) Label() string {
	switch v {
	case STATEMENT_STATUS_UNKNOWN:
		return ""
	case STATEMENT_STATUS__COMPLETE:
		return "已结算"
	case STATEMENT_STATUS__CREATED:
		return "未结算"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*StatementStatus)(nil)

func (v StatementStatus) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidStatementStatus
	}
	return []byte(str), nil
}

func (v *StatementStatus) UnmarshalText(data []byte) (err error) {
	*v, err = ParseStatementStatusFromString(string(bytes.ToUpper(data)))
	return
}
