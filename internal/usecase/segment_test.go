package usecase_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"id-maker/internal/entity"
	"testing"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	res  interface{}
	err  error
}

func segment(t *testing.T) (*MockSegment, *MockSegmentRepo) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	s := NewMockSegment(mockCtrl)
	repo := NewMockSegmentRepo(mockCtrl)

	return s, repo
}

func TestGetid(t *testing.T) {
	t.Parallel()
	segment, _ := segment(t)

	tests := []test{
		{
			name: "right result",
			mock: func() {
				segment.EXPECT().GetId("test").Return(int64(100), nil)
				segment.EXPECT().SnowFlakeGetId().Return(int64(100))
			},
			res: int64(100),
			err: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()
			res, err := segment.GetId("test")

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestCreateTag(t *testing.T) {
	t.Parallel()

	segment, _ := segment(t)

	tests := []test{
		{
			name: "create tag",
			mock: func() {
				segment.EXPECT().CreateTag(&entity.Segments{}).Return(nil)
			},
			res: nil,
			err: nil,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()
			err := segment.CreateTag(&entity.Segments{})

			require.ErrorIs(t, err, tc.err)
		})
	}
}

func TestGetList(t *testing.T) {
	t.Parallel()

	_, repo := segment(t)

	tests := []test{
		{
			name: "empty result",
			mock: func() {
				repo.EXPECT().GetList().Return(nil, nil)
			},
			res: []entity.Segments(nil),
			err: nil,
		},
	}
	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()
			res, err := repo.GetList()
			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
