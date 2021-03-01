package queries

//go:generate mockgen -package mocks -destination mocks/db_mock.go github.com/bicycolet/bicycolet/internal/db/queries/query Txn,Rows,Result
//go:generate mockgen -package mocks -destination mocks/statements_mock.go github.com/bicycolet/bicycolet/internal/db/queries/query Statements,Query
