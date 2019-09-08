package query_test

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

//go:generate mockgen -package mocks -destination mocks/db_mock.go github.com/bicycolet/bicycolet/internal/db/database DB,Tx,Rows,ColumnType
//go:generate mockgen -package mocks -destination mocks/result_mock.go database/sql Result

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	os.Exit(m.Run())
}

type intScanMatcher struct {
	x int
}

func IntScanMatcher(v int) gomock.Matcher {
	return intScanMatcher{
		x: v,
	}
}

func (m intScanMatcher) Matches(x interface{}) bool {
	ref := reflect.ValueOf(x).Elem()
	ref.Set(reflect.ValueOf(m.x))
	return true
}

func (m intScanMatcher) String() string {
	return fmt.Sprintf("%v", m.x)
}

type stringScanMatcher struct {
	x string
}

func StringScanMatcher(v string) gomock.Matcher {
	return stringScanMatcher{
		x: v,
	}
}

func (m stringScanMatcher) Matches(x interface{}) bool {
	ref := reflect.ValueOf(x).Elem()
	ref.Set(reflect.ValueOf(m.x))
	return true
}

func (m stringScanMatcher) String() string {
	return m.x
}
