package sql

import (
	"errors"
	"strings"
	"testing"

	"github.com/bicycolet/bicycolet/internal/db/statements/mocks"
	"github.com/golang/mock/gomock"
)

func TestRegistry(t *testing.T) {
	t.Parallel()

	t.Run("get or store", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sql := "SELECT * FROM table;"

		hasher := mocks.NewMockHasher(ctrl)
		hasher.EXPECT().Hash(sql).Return("aaa")

		prep := mocks.NewMockPreparer(ctrl)
		prep.EXPECT().Prepare(sql).Return(nil, nil)

		reg := New(prep, hasher)
		_, err := reg.Create(sql)
		if expected, actual := true, err == nil; expected != actual {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})

	t.Run("get or store fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sql := "SELECT * FROM table;"
		fail := errors.New("fail")

		hasher := mocks.NewMockHasher(ctrl)
		hasher.EXPECT().Hash(sql).Return("aaa")

		prep := mocks.NewMockPreparer(ctrl)
		prep.EXPECT().Prepare(sql).Return(nil, fail)

		reg := New(prep, hasher)
		_, err := reg.Create(sql)
		if expected, actual := fail.Error(), err.Error(); !strings.Contains(actual, expected) {
			t.Errorf("expected: %v, actual: %v, err: %v", expected, actual, err)
		}
	})
}
