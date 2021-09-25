package enums

import (
	"bytes"
	"encoding"
	"errors"

	github_com_eden_framework_enumeration "github.com/eden-framework/enumeration"
)

var InvalidStatementType = errors.New("invalid StatementType")

func init() {
	github_com_eden_framework_enumeration.RegisterEnums("StatementType", map[string]string{
		"MONTH": "月",
		"WEEK":  "周",
	})
}

func ParseStatementTypeFromString(s string) (StatementType, error) {
	switch s {
	case "":
		return STATEMENT_TYPE_UNKNOWN, nil
	case "MONTH":
		return STATEMENT_TYPE__MONTH, nil
	case "WEEK":
		return STATEMENT_TYPE__WEEK, nil
	}
	return STATEMENT_TYPE_UNKNOWN, InvalidStatementType
}

func ParseStatementTypeFromLabelString(s string) (StatementType, error) {
	switch s {
	case "":
		return STATEMENT_TYPE_UNKNOWN, nil
	case "月":
		return STATEMENT_TYPE__MONTH, nil
	case "周":
		return STATEMENT_TYPE__WEEK, nil
	}
	return STATEMENT_TYPE_UNKNOWN, InvalidStatementType
}

func (StatementType) EnumType() string {
	return "StatementType"
}

func (StatementType) Enums() map[int][]string {
	return map[int][]string{
		int(STATEMENT_TYPE__MONTH): {"MONTH", "月"},
		int(STATEMENT_TYPE__WEEK):  {"WEEK", "周"},
	}
}

func (v StatementType) String() string {
	switch v {
	case STATEMENT_TYPE_UNKNOWN:
		return ""
	case STATEMENT_TYPE__MONTH:
		return "MONTH"
	case STATEMENT_TYPE__WEEK:
		return "WEEK"
	}
	return "UNKNOWN"
}

func (v StatementType) Label() string {
	switch v {
	case STATEMENT_TYPE_UNKNOWN:
		return ""
	case STATEMENT_TYPE__MONTH:
		return "月"
	case STATEMENT_TYPE__WEEK:
		return "周"
	}
	return "UNKNOWN"
}

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*StatementType)(nil)

func (v StatementType) MarshalText() ([]byte, error) {
	str := v.String()
	if str == "UNKNOWN" {
		return nil, InvalidStatementType
	}
	return []byte(str), nil
}

func (v *StatementType) UnmarshalText(data []byte) (err error) {
	*v, err = ParseStatementTypeFromString(string(bytes.ToUpper(data)))
	return
}
