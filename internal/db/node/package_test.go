package node_test

import (
	"fmt"
	"reflect"

	"github.com/bicycolet/bicycolet/internal/db/database"
	"github.com/golang/mock/gomock"
)

//go:generate mockgen -package mocks -destination mocks/database_mock.go github.com/bicycolet/bicycolet/internal/db/node DatabaseRegistrar,DatabaseOpener,DatabaseIO
//go:generate mockgen -package mocks -destination mocks/filesystem_mock.go github.com/bicycolet/bicycolet/internal/fsys FileSystem
//go:generate mockgen -package mocks -destination mocks/schema_mock.go github.com/bicycolet/bicycolet/internal/db/node Schema,SchemaProvider
//go:generate mockgen -package mocks -destination mocks/db_mock.go github.com/bicycolet/bicycolet/internal/db/database DB,Tx,Rows
//go:generate mockgen -package mocks -destination mocks/driver_mock.go database/sql/driver Driver
//go:generate mockgen -package mocks -destination mocks/result_mock.go database/sql Result

func connectionInfo() database.ConnectionInfo {
	return database.ConnectionInfo{
		Host:     "localhost",
		Port:     5435,
		User:     "postgres",
		Password: "postgres",
		DBName:   "test",
	}
}

type scanMatcher struct {
	x int
}

func ScanMatcher(v int) gomock.Matcher {
	return scanMatcher{
		x: v,
	}
}

func (m scanMatcher) Matches(x interface{}) bool {
	ref := reflect.ValueOf(x).Elem()
	ref.Set(reflect.ValueOf(m.x))
	return true
}

func (m scanMatcher) String() string {
	return fmt.Sprintf("%v", m.x)
}

type typeMatcher struct {
	x reflect.Kind
}

func TypeMatcher(v reflect.Kind) gomock.Matcher {
	return typeMatcher{
		x: v,
	}
}

func (m typeMatcher) Matches(x interface{}) bool {
	v := reflect.ValueOf(x)
	switch v.Kind() {
	case m.x:
		return true
	}
	return false
}

func (m typeMatcher) String() string {
	return fmt.Sprintf("%v", m.x)
}
