package statements

//go:generate mockgen -package mocks -destination mocks/mocks.go github.com/bicycolet/bicycolet/internal/db/statements/statement Preparer,Hasher
