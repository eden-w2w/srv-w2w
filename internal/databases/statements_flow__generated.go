package databases

import (
	fmt "fmt"
	time "time"

	github_com_eden_framework_sqlx "github.com/eden-framework/sqlx"
	github_com_eden_framework_sqlx_builder "github.com/eden-framework/sqlx/builder"
	github_com_eden_framework_sqlx_datatypes "github.com/eden-framework/sqlx/datatypes"
)

func (StatementsFlow) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (StatementsFlow) UniqueIndexUInterval() string {
	return "U_interval"
}

func (StatementsFlow) UniqueIndexUStatementsID() string {
	return "U_statements_id"
}

func (StatementsFlow) UniqueIndexes() github_com_eden_framework_sqlx_builder.Indexes {
	return github_com_eden_framework_sqlx_builder.Indexes{
		"U_interval": []string{
			"UserID",
			"Name",
			"DeletedAt",
		},
		"U_statements_id": []string{
			"StatementsID",
			"DeletedAt",
		},
	}
}

func (StatementsFlow) Comments() map[string]string {
	return map[string]string{
		"Amount":       "结算金额",
		"Name":         "名称",
		"Proportion":   "计算比例",
		"StatementsID": "结算单ID",
		"Status":       "结算状态",
		"TotalSales":   "销售总额",
		"UserID":       "用户ID",
	}
}

var StatementsFlowTable *github_com_eden_framework_sqlx_builder.Table

func init() {
	StatementsFlowTable = Config.DB.Register(&StatementsFlow{})
}

type StatementsFlowIterator struct {
}

func (StatementsFlowIterator) New() interface{} {
	return &StatementsFlow{}
}

func (StatementsFlowIterator) Resolve(v interface{}) *StatementsFlow {
	return v.(*StatementsFlow)
}

func (StatementsFlow) TableName() string {
	return "t_statements_flow"
}

func (StatementsFlow) ColDescriptions() map[string][]string {
	return map[string][]string{
		"Amount": []string{
			"结算金额",
		},
		"Name": []string{
			"名称",
		},
		"Proportion": []string{
			"计算比例",
		},
		"StatementsID": []string{
			"结算单ID",
		},
		"Status": []string{
			"结算状态",
		},
		"TotalSales": []string{
			"销售总额",
		},
		"UserID": []string{
			"用户ID",
		},
	}
}

func (StatementsFlow) FieldKeyID() string {
	return "ID"
}

func (m *StatementsFlow) FieldID() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyID())
}

func (StatementsFlow) FieldKeyStatementsID() string {
	return "StatementsID"
}

func (m *StatementsFlow) FieldStatementsID() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyStatementsID())
}

func (StatementsFlow) FieldKeyUserID() string {
	return "UserID"
}

func (m *StatementsFlow) FieldUserID() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyUserID())
}

func (StatementsFlow) FieldKeyName() string {
	return "Name"
}

func (m *StatementsFlow) FieldName() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyName())
}

func (StatementsFlow) FieldKeyTotalSales() string {
	return "TotalSales"
}

func (m *StatementsFlow) FieldTotalSales() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyTotalSales())
}

func (StatementsFlow) FieldKeyProportion() string {
	return "Proportion"
}

func (m *StatementsFlow) FieldProportion() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyProportion())
}

func (StatementsFlow) FieldKeyAmount() string {
	return "Amount"
}

func (m *StatementsFlow) FieldAmount() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyAmount())
}

func (StatementsFlow) FieldKeyStatus() string {
	return "Status"
}

func (m *StatementsFlow) FieldStatus() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyStatus())
}

func (StatementsFlow) FieldKeyCreatedAt() string {
	return "CreatedAt"
}

func (m *StatementsFlow) FieldCreatedAt() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyCreatedAt())
}

func (StatementsFlow) FieldKeyUpdatedAt() string {
	return "UpdatedAt"
}

func (m *StatementsFlow) FieldUpdatedAt() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyUpdatedAt())
}

func (StatementsFlow) FieldKeyDeletedAt() string {
	return "DeletedAt"
}

func (m *StatementsFlow) FieldDeletedAt() *github_com_eden_framework_sqlx_builder.Column {
	return StatementsFlowTable.F(m.FieldKeyDeletedAt())
}

func (StatementsFlow) ColRelations() map[string][]string {
	return map[string][]string{}
}

func (m *StatementsFlow) IndexFieldNames() []string {
	return []string{
		"ID",
		"Name",
		"StatementsID",
		"UserID",
	}
}

func (m *StatementsFlow) ConditionByStruct(db github_com_eden_framework_sqlx.DBExecutor) github_com_eden_framework_sqlx_builder.SqlCondition {
	table := db.T(m)
	fieldValues := github_com_eden_framework_sqlx_builder.FieldValuesFromStructByNonZero(m)

	conditions := make([]github_com_eden_framework_sqlx_builder.SqlCondition, 0)

	for _, fieldName := range m.IndexFieldNames() {
		if v, exists := fieldValues[fieldName]; exists {
			conditions = append(conditions, table.F(fieldName).Eq(v))
			delete(fieldValues, fieldName)
		}
	}

	if len(conditions) == 0 {
		panic(fmt.Errorf("at least one of field for indexes has value"))
	}

	for fieldName, v := range fieldValues {
		conditions = append(conditions, table.F(fieldName).Eq(v))
	}

	condition := github_com_eden_framework_sqlx_builder.And(conditions...)

	condition = github_com_eden_framework_sqlx_builder.And(condition, table.F("DeletedAt").Eq(0))
	return condition
}

func (m *StatementsFlow) Create(db github_com_eden_framework_sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	_, err := db.ExecExpr(github_com_eden_framework_sqlx.InsertToDB(db, m, nil))
	return err

}

func (m *StatementsFlow) CreateOnDuplicateWithUpdateFields(db github_com_eden_framework_sqlx.DBExecutor, updateFields []string) error {

	if len(updateFields) == 0 {
		panic(fmt.Errorf("must have update fields"))
	}

	if m.CreatedAt.IsZero() {
		m.CreatedAt = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	fieldValues := github_com_eden_framework_sqlx_builder.FieldValuesFromStructByNonZero(m, updateFields...)

	delete(fieldValues, "ID")

	table := db.T(m)

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	fields := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		fields[field] = true
	}

	for _, fieldNames := range m.UniqueIndexes() {
		for _, field := range fieldNames {
			delete(fields, field)
		}
	}

	if len(fields) == 0 {
		panic(fmt.Errorf("no fields for updates"))
	}

	for field := range fieldValues {
		if !fields[field] {
			delete(fieldValues, field)
		}
	}

	additions := github_com_eden_framework_sqlx_builder.Additions{}

	switch db.Dialect().DriverName() {
	case "mysql":
		additions = append(additions, github_com_eden_framework_sqlx_builder.OnDuplicateKeyUpdate(table.AssignmentsByFieldValues(fieldValues)...))
	case "postgres":
		indexes := m.UniqueIndexes()
		fields := make([]string, 0)
		for _, fs := range indexes {
			fields = append(fields, fs...)
		}
		indexFields, _ := db.T(m).Fields(fields...)

		additions = append(additions,
			github_com_eden_framework_sqlx_builder.OnConflict(indexFields).
				DoUpdateSet(table.AssignmentsByFieldValues(fieldValues)...))
	}

	additions = append(additions, github_com_eden_framework_sqlx_builder.Comment("User.CreateOnDuplicateWithUpdateFields"))

	expr := github_com_eden_framework_sqlx_builder.Insert().Into(table, additions...).Values(cols, vals...)

	_, err := db.ExecExpr(expr)
	return err

}

func (m *StatementsFlow) DeleteByStruct(db github_com_eden_framework_sqlx.DBExecutor) error {

	_, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Delete().
			From(
				db.T(m),
				github_com_eden_framework_sqlx_builder.Where(m.ConditionByStruct(db)),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.DeleteByStruct"),
			),
	)

	return err
}

func (m *StatementsFlow) FetchByID(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	err := db.QueryExprAndScan(
		github_com_eden_framework_sqlx_builder.Select(nil).
			From(
				db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("ID").Eq(m.ID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.FetchByID"),
			),
		m,
	)

	return err
}

func (m *StatementsFlow) UpdateByIDWithMap(db github_com_eden_framework_sqlx.DBExecutor, fieldValues github_com_eden_framework_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdatedAt"]; !ok {
		fieldValues["UpdatedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	table := db.T(m)

	result, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Update(db.T(m)).
			Where(
				github_com_eden_framework_sqlx_builder.And(
					table.F("ID").Eq(m.ID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.UpdateByIDWithMap"),
			).
			Set(table.AssignmentsByFieldValues(fieldValues)...),
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return m.FetchByID(db)
	}

	return nil

}

func (m *StatementsFlow) UpdateByIDWithStruct(db github_com_eden_framework_sqlx.DBExecutor, zeroFields ...string) error {

	fieldValues := github_com_eden_framework_sqlx_builder.FieldValuesFromStructByNonZero(m, zeroFields...)
	return m.UpdateByIDWithMap(db, fieldValues)

}

func (m *StatementsFlow) FetchByIDForUpdate(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	err := db.QueryExprAndScan(
		github_com_eden_framework_sqlx_builder.Select(nil).
			From(
				db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("ID").Eq(m.ID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.ForUpdate(),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.FetchByIDForUpdate"),
			),
		m,
	)

	return err
}

func (m *StatementsFlow) DeleteByID(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	_, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Delete().
			From(db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("ID").Eq(m.ID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.DeleteByID"),
			))

	return err
}

func (m *StatementsFlow) SoftDeleteByID(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	fieldValues := github_com_eden_framework_sqlx_builder.FieldValues{}
	if _, ok := fieldValues["DeletedAt"]; !ok {
		fieldValues["DeletedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	if _, ok := fieldValues["UpdatedAt"]; !ok {
		fieldValues["UpdatedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	_, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Update(db.T(m)).
			Where(
				github_com_eden_framework_sqlx_builder.And(
					table.F("ID").Eq(m.ID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.SoftDeleteByID"),
			).
			Set(table.AssignmentsByFieldValues(fieldValues)...),
	)

	return err

}

func (m *StatementsFlow) FetchByStatementsID(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	err := db.QueryExprAndScan(
		github_com_eden_framework_sqlx_builder.Select(nil).
			From(
				db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("StatementsID").Eq(m.StatementsID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.FetchByStatementsID"),
			),
		m,
	)

	return err
}

func (m *StatementsFlow) UpdateByStatementsIDWithMap(db github_com_eden_framework_sqlx.DBExecutor, fieldValues github_com_eden_framework_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdatedAt"]; !ok {
		fieldValues["UpdatedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	table := db.T(m)

	result, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Update(db.T(m)).
			Where(
				github_com_eden_framework_sqlx_builder.And(
					table.F("StatementsID").Eq(m.StatementsID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.UpdateByStatementsIDWithMap"),
			).
			Set(table.AssignmentsByFieldValues(fieldValues)...),
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return m.FetchByStatementsID(db)
	}

	return nil

}

func (m *StatementsFlow) UpdateByStatementsIDWithStruct(db github_com_eden_framework_sqlx.DBExecutor, zeroFields ...string) error {

	fieldValues := github_com_eden_framework_sqlx_builder.FieldValuesFromStructByNonZero(m, zeroFields...)
	return m.UpdateByStatementsIDWithMap(db, fieldValues)

}

func (m *StatementsFlow) FetchByStatementsIDForUpdate(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	err := db.QueryExprAndScan(
		github_com_eden_framework_sqlx_builder.Select(nil).
			From(
				db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("StatementsID").Eq(m.StatementsID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.ForUpdate(),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.FetchByStatementsIDForUpdate"),
			),
		m,
	)

	return err
}

func (m *StatementsFlow) DeleteByStatementsID(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	_, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Delete().
			From(db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("StatementsID").Eq(m.StatementsID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.DeleteByStatementsID"),
			))

	return err
}

func (m *StatementsFlow) SoftDeleteByStatementsID(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	fieldValues := github_com_eden_framework_sqlx_builder.FieldValues{}
	if _, ok := fieldValues["DeletedAt"]; !ok {
		fieldValues["DeletedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	if _, ok := fieldValues["UpdatedAt"]; !ok {
		fieldValues["UpdatedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	_, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Update(db.T(m)).
			Where(
				github_com_eden_framework_sqlx_builder.And(
					table.F("StatementsID").Eq(m.StatementsID),
					table.F("DeletedAt").Eq(m.DeletedAt),
				),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.SoftDeleteByStatementsID"),
			).
			Set(table.AssignmentsByFieldValues(fieldValues)...),
	)

	return err

}

func (m *StatementsFlow) FetchByUserIDAndName(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	err := db.QueryExprAndScan(
		github_com_eden_framework_sqlx_builder.Select(nil).
			From(
				db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("UserID").Eq(m.UserID),
					table.F("Name").Eq(m.Name),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.FetchByUserIDAndName"),
			),
		m,
	)

	return err
}

func (m *StatementsFlow) UpdateByUserIDAndNameWithMap(db github_com_eden_framework_sqlx.DBExecutor, fieldValues github_com_eden_framework_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdatedAt"]; !ok {
		fieldValues["UpdatedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	table := db.T(m)

	result, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Update(db.T(m)).
			Where(
				github_com_eden_framework_sqlx_builder.And(
					table.F("UserID").Eq(m.UserID),
					table.F("Name").Eq(m.Name),
					table.F("DeletedAt").Eq(m.DeletedAt),
				),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.UpdateByUserIDAndNameWithMap"),
			).
			Set(table.AssignmentsByFieldValues(fieldValues)...),
	)

	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return m.FetchByUserIDAndName(db)
	}

	return nil

}

func (m *StatementsFlow) UpdateByUserIDAndNameWithStruct(db github_com_eden_framework_sqlx.DBExecutor, zeroFields ...string) error {

	fieldValues := github_com_eden_framework_sqlx_builder.FieldValuesFromStructByNonZero(m, zeroFields...)
	return m.UpdateByUserIDAndNameWithMap(db, fieldValues)

}

func (m *StatementsFlow) FetchByUserIDAndNameForUpdate(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	err := db.QueryExprAndScan(
		github_com_eden_framework_sqlx_builder.Select(nil).
			From(
				db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("UserID").Eq(m.UserID),
					table.F("Name").Eq(m.Name),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.ForUpdate(),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.FetchByUserIDAndNameForUpdate"),
			),
		m,
	)

	return err
}

func (m *StatementsFlow) DeleteByUserIDAndName(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	_, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Delete().
			From(db.T(m),
				github_com_eden_framework_sqlx_builder.Where(github_com_eden_framework_sqlx_builder.And(
					table.F("UserID").Eq(m.UserID),
					table.F("Name").Eq(m.Name),
					table.F("DeletedAt").Eq(m.DeletedAt),
				)),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.DeleteByUserIDAndName"),
			))

	return err
}

func (m *StatementsFlow) SoftDeleteByUserIDAndName(db github_com_eden_framework_sqlx.DBExecutor) error {

	table := db.T(m)

	fieldValues := github_com_eden_framework_sqlx_builder.FieldValues{}
	if _, ok := fieldValues["DeletedAt"]; !ok {
		fieldValues["DeletedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	if _, ok := fieldValues["UpdatedAt"]; !ok {
		fieldValues["UpdatedAt"] = github_com_eden_framework_sqlx_datatypes.Timestamp(time.Now())
	}

	_, err := db.ExecExpr(
		github_com_eden_framework_sqlx_builder.Update(db.T(m)).
			Where(
				github_com_eden_framework_sqlx_builder.And(
					table.F("UserID").Eq(m.UserID),
					table.F("Name").Eq(m.Name),
					table.F("DeletedAt").Eq(m.DeletedAt),
				),
				github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.SoftDeleteByUserIDAndName"),
			).
			Set(table.AssignmentsByFieldValues(fieldValues)...),
	)

	return err

}

func (m *StatementsFlow) List(db github_com_eden_framework_sqlx.DBExecutor, condition github_com_eden_framework_sqlx_builder.SqlCondition, additions ...github_com_eden_framework_sqlx_builder.Addition) ([]StatementsFlow, error) {

	list := make([]StatementsFlow, 0)

	table := db.T(m)
	_ = table

	condition = github_com_eden_framework_sqlx_builder.And(condition, table.F("DeletedAt").Eq(0))

	finalAdditions := []github_com_eden_framework_sqlx_builder.Addition{
		github_com_eden_framework_sqlx_builder.Where(condition),
		github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.List"),
	}

	if len(additions) > 0 {
		finalAdditions = append(finalAdditions, additions...)
	}

	err := db.QueryExprAndScan(
		github_com_eden_framework_sqlx_builder.Select(nil).
			From(db.T(m), finalAdditions...),
		&list,
	)

	return list, err

}

func (m *StatementsFlow) Count(db github_com_eden_framework_sqlx.DBExecutor, condition github_com_eden_framework_sqlx_builder.SqlCondition, additions ...github_com_eden_framework_sqlx_builder.Addition) (int, error) {

	count := -1

	table := db.T(m)
	_ = table

	condition = github_com_eden_framework_sqlx_builder.And(condition, table.F("DeletedAt").Eq(0))

	finalAdditions := []github_com_eden_framework_sqlx_builder.Addition{
		github_com_eden_framework_sqlx_builder.Where(condition),
		github_com_eden_framework_sqlx_builder.Comment("StatementsFlow.Count"),
	}

	if len(additions) > 0 {
		finalAdditions = append(finalAdditions, additions...)
	}

	err := db.QueryExprAndScan(
		github_com_eden_framework_sqlx_builder.Select(
			github_com_eden_framework_sqlx_builder.Count(),
		).
			From(db.T(m), finalAdditions...),
		&count,
	)

	return count, err

}

func (m *StatementsFlow) BatchFetchByIDList(db github_com_eden_framework_sqlx.DBExecutor, values []uint64) ([]StatementsFlow, error) {

	if len(values) == 0 {
		return nil, nil
	}

	table := db.T(m)

	condition := table.F("ID").In(values)

	return m.List(db, condition)

}

func (m *StatementsFlow) BatchFetchByNameList(db github_com_eden_framework_sqlx.DBExecutor, values []string) ([]StatementsFlow, error) {

	if len(values) == 0 {
		return nil, nil
	}

	table := db.T(m)

	condition := table.F("Name").In(values)

	return m.List(db, condition)

}

func (m *StatementsFlow) BatchFetchByStatementsIDList(db github_com_eden_framework_sqlx.DBExecutor, values []uint64) ([]StatementsFlow, error) {

	if len(values) == 0 {
		return nil, nil
	}

	table := db.T(m)

	condition := table.F("StatementsID").In(values)

	return m.List(db, condition)

}

func (m *StatementsFlow) BatchFetchByUserIDList(db github_com_eden_framework_sqlx.DBExecutor, values []uint64) ([]StatementsFlow, error) {

	if len(values) == 0 {
		return nil, nil
	}

	table := db.T(m)

	condition := table.F("UserID").In(values)

	return m.List(db, condition)

}
