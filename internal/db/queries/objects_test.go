package queries

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/queries/mocks"
	"github.com/bicycolet/bicycolet/internal/db/queries/query"
	"github.com/golang/mock/gomock"
)

func TestSelectObjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Query("SELECT * FROM schema").Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).Do(func(a *int, b *int) error {
			*a = 1
			*b = 2
			return nil
		}).Return(nil),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Err().Return(nil),
		mockRows.EXPECT().Close().Return(nil),
	)

	queries := &Query{
		statements: mockStatements,
	}

	record := make([]int, 2)
	err := queries.SelectObjects(mockTx, func(i int) []interface{} {
		return []interface{}{
			&record[0],
			&record[1],
		}
	}, "SELECT * FROM schema")
	if err != nil {
		t.Errorf("expected err to be nil")
	}
	if expected, actual := []int{1, 2}, record; !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestSelectObjectsWithQueryFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Query("SELECT * FROM schema").Return(mockRows, errors.New("bad")),
	)

	queries := &Query{
		statements: mockStatements,
	}
	err := queries.SelectObjects(mockTx, func(i int) []interface{} {
		return []interface{}{}
	}, "SELECT * FROM schema")
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestSelectObjectsWithScanFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Query("SELECT * FROM schema").Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).Do(func(a *int, b *int) error {
			*a = 1
			*b = 2
			return nil
		}).Return(errors.New("bad")),
		mockRows.EXPECT().Close().Return(nil),
	)

	queries := &Query{
		statements: mockStatements,
	}

	record := make([]int, 2)
	err := queries.SelectObjects(mockTx, func(i int) []interface{} {
		return []interface{}{
			&record[0],
			&record[1],
		}
	}, "SELECT * FROM schema")
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestUpsertObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockResult := mocks.NewMockResult(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockStatements.EXPECT().UpsertObject(query.Table("schema"), []string{"id", "name"}).Return(mockQuery),
		mockQuery.EXPECT().Exec(mockTx, []interface{}{1, "fred"}).Return(mockResult, nil),
		mockResult.EXPECT().RowsAffected().Return(int64(1), nil),
	)

	queries := &Query{
		statements: mockStatements,
	}

	err := queries.UpsertObject(mockTx, "schema", []string{"id", "name"}, []interface{}{1, "fred"})
	if err != nil {
		t.Errorf("expected err to be nil")
	}
}

/*
func TestUpsertObjectWithExecFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockResult := mocks.NewMockResult(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Exec("INSERT OR REPLACE INTO schema (id, name) VALUES (?, ?)", []interface{}{1, "fred"}).Return(mockResult, errors.New("bad")),
	)

	queries := &Query{
		statements: mockStatements,
	}
	err := queries.UpsertObject(mockTx, "schema", []string{"id", "name"}, []interface{}{1, "fred"})
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestUpsertObjectWithLastInsertedIdFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockResult := mocks.NewMockResult(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Exec("INSERT OR REPLACE INTO schema (id, name) VALUES (?, ?)", []interface{}{1, "fred"}).Return(mockResult, nil),
		mockResult.EXPECT().LastInsertId().Return(int64(5), errors.New("bad")),
	)

	queries := &Query{
		statements: mockStatements,
	}
	err := queries.UpsertObject(mockTx, "schema", []string{"id", "name"}, []interface{}{1, "fred"})
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestUpsertObjectWithNoColumns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	queries := &Query{
		statements: mockStatements,
	}
	err := queries.UpsertObject(mockTx, "schema", []string{}, []interface{}{1, "fred"})
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestUpsertObjectWithColumnsValuesMissmatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	queries := &Query{
		statements: mockStatements,
	}
	err := queries.UpsertObject(mockTx, "schema", []string{"id"}, []interface{}{1, "fred"})
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestDeleteObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockResult := mocks.NewMockResult(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Exec("DELETE FROM schema WHERE id=?", int64(1)).Return(mockResult, nil),
		mockResult.EXPECT().RowsAffected().Return(int64(1), nil),
	)

	queries := &Query{
		statements: mockStatements,
	}
	ok, err := queries.DeleteObject(mockTx, "schema", int64(1))
	if err != nil {
		t.Errorf("expected err to be nil")
	}
	if expected, actual := true, ok; expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestDeleteObjectWithExecFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockResult := mocks.NewMockResult(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Exec("DELETE FROM schema WHERE id=?", int64(1)).Return(mockResult, errors.New("bad")),
	)

	queries := &Query{
		statements: mockStatements,
	}
	_, err := queries.DeleteObject(mockTx, "schema", int64(1))
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestDeleteObjectWithRowsAffectedFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockResult := mocks.NewMockResult(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Exec("DELETE FROM schema WHERE id=?", int64(1)).Return(mockResult, nil),
		mockResult.EXPECT().RowsAffected().Return(int64(1), errors.New("bad")),
	)

	queries := &Query{
		statements: mockStatements,
	}
	_, err := queries.DeleteObject(mockTx, "schema", int64(1))
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestDeleteObjectWithTooManyRowsAffected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockResult := mocks.NewMockResult(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockTx.EXPECT().Exec("DELETE FROM schema WHERE id=?", int64(1)).Return(mockResult, nil),
		mockResult.EXPECT().RowsAffected().Return(int64(2), nil),
	)

	queries := &Query{
		statements: mockStatements,
	}
	_, err := queries.DeleteObject(mockTx, "schema", int64(1))
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}
*/
