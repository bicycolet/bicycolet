package queries

import (
	"errors"
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/queries/mocks"
	"github.com/bicycolet/bicycolet/internal/db/queries/query"
	"github.com/golang/mock/gomock"
)

func TestCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockStatements.EXPECT().Op(query.Equal).Return("="),
		mockStatements.EXPECT().Params(1).Return("(?)"),
		mockStatements.EXPECT().Count(query.Table("schema"), equals("a", "b")).Return(mockQuery),
		mockQuery.EXPECT().Run(mockTx, "b").Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).Do(func(x *int) error {
			*x = 1
			return nil
		}).Return(nil),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Err().Return(nil),
		mockRows.EXPECT().Close().Return(nil),
	)

	queries := &Query{
		statements: mockStatements,
	}
	count, err := queries.Count(mockTx, "schema", query.Equals("a", "b"))
	if err != nil {
		t.Errorf("expected err to be nil")
	}
	if expected, actual := 1, count; expected != actual {
		t.Errorf("expected: %d, actual: %d", expected, actual)
	}
}

func TestCountWithQueryFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockStatements.EXPECT().Op(query.Equal).Return("="),
		mockStatements.EXPECT().Params(1).Return("(?)"),
		mockStatements.EXPECT().Count(query.Table("schema"), equals("a", "b")).Return(mockQuery),
		mockQuery.EXPECT().Run(mockTx, "b").Return(mockRows, errors.New("bad")),
	)

	queries := &Query{
		statements: mockStatements,
	}
	_, err := queries.Count(mockTx, "schema", query.Equals("a", "b"))
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestCountWithNoNext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockStatements.EXPECT().Op(query.Equal).Return("="),
		mockStatements.EXPECT().Params(1).Return("(?)"),
		mockStatements.EXPECT().Count(query.Table("schema"), equals("a", "b")).Return(mockQuery),
		mockQuery.EXPECT().Run(mockTx, "b").Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Close().Return(nil),
	)

	queries := &Query{
		statements: mockStatements,
	}
	_, err := queries.Count(mockTx, "schema", query.Equals("a", "b"))
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestCountWithScanFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockStatements.EXPECT().Op(query.Equal).Return("="),
		mockStatements.EXPECT().Params(1).Return("(?)"),
		mockStatements.EXPECT().Count(query.Table("schema"), equals("a", "b")).Return(mockQuery),
		mockQuery.EXPECT().Run(mockTx, "b").Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).Do(func(x *int) error {
			*x = 1
			return nil
		}).Return(errors.New("bad")),
		mockRows.EXPECT().Close().Return(nil),
	)

	queries := &Query{
		statements: mockStatements,
	}
	_, err := queries.Count(mockTx, "schema", query.Equals("a", "b"))
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestCountWithMoreRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockStatements.EXPECT().Op(query.Equal).Return("="),
		mockStatements.EXPECT().Params(1).Return("(?)"),
		mockStatements.EXPECT().Count(query.Table("schema"), equals("a", "b")).Return(mockQuery),
		mockQuery.EXPECT().Run(mockTx, "b").Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).Do(func(x *int) error {
			*x = 1
			return nil
		}).Return(nil),
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Close().Return(nil),
	)

	queries := &Query{
		statements: mockStatements,
	}
	_, err := queries.Count(mockTx, "schema", query.Equals("a", "b"))
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func TestCountWithErrFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTx := mocks.NewMockTxn(ctrl)
	mockRows := mocks.NewMockRows(ctrl)
	mockQuery := mocks.NewMockQuery(ctrl)
	mockStatements := mocks.NewMockStatements(ctrl)

	gomock.InOrder(
		mockStatements.EXPECT().Op(query.Equal).Return("="),
		mockStatements.EXPECT().Params(1).Return("(?)"),
		mockStatements.EXPECT().Count(query.Table("schema"), equals("a", "b")).Return(mockQuery),
		mockQuery.EXPECT().Run(mockTx, "b").Return(mockRows, nil),
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any()).Do(func(x *int) error {
			*x = 1
			return nil
		}).Return(nil),
		mockRows.EXPECT().Next().Return(false),
		mockRows.EXPECT().Err().Return(errors.New("bad")),
		mockRows.EXPECT().Close().Return(nil),
	)

	queries := &Query{
		statements: mockStatements,
	}
	_, err := queries.Count(mockTx, "schema", query.Equals("a", "b"))
	if err == nil {
		t.Errorf("expected err not to be nil")
	}
}

func equals(a, b string) query.Expression {
	return query.Statement{
		Name:  a,
		Value: b,
		Op:    query.Equal,
	}
}
