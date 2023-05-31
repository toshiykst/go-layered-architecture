package database

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/testutil"
)

func TestDatabase_dbGroupRepository_Find(t *testing.T) {
	tests := []struct {
		name           string
		gID            model.GroupID
		want           *model.Group
		wantErr        error
		dbGroupErr     error
		dbGroupUserErr error
	}{
		{
			name: "Returns a group",
			gID:  model.GroupID("TEST_GROUP_ID"),
			want: model.NewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				},
			),
			wantErr:        nil,
			dbGroupErr:     nil,
			dbGroupUserErr: nil,
		},
		{
			name:           "Not found",
			gID:            model.GroupID("TEST_GROUP_ID"),
			want:           nil,
			wantErr:        nil,
			dbGroupErr:     gorm.ErrRecordNotFound,
			dbGroupUserErr: nil,
		},
		{
			name:           "DB group error",
			gID:            model.GroupID("TEST_GROUP_ID"),
			want:           nil,
			wantErr:        errors.New("an error occurred"),
			dbGroupErr:     errors.New("an error occurred"),
			dbGroupUserErr: nil,
		},
		{
			name:           "DB groupuser error",
			gID:            model.GroupID("TEST_GROUP_ID"),
			want:           nil,
			wantErr:        errors.New("an error occurred"),
			dbGroupErr:     nil,
			dbGroupUserErr: errors.New("an error occurred"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, db := testutil.DBMock(t)
			sqlDB, err := db.DB()
			if err != nil {
				t.Fatalf("want no err, but has error %v", err)
			}
			defer sqlDB.Close()

			groupsSQL := "SELECT * FROM `groups` WHERE `groups`.`id` = ? ORDER BY `groups`.`id` LIMIT 1"
			groupsExpectQuery := mock.
				ExpectQuery(regexp.QuoteMeta(groupsSQL)).
				WithArgs(tt.gID)

			now := time.Now()

			if tt.dbGroupErr != nil {
				groupsExpectQuery.WillReturnError(tt.dbGroupErr)
			} else {
				groupRows := sqlmock.
					NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(tt.want.ID(), tt.want.Name(), now, now)
				groupsExpectQuery.WillReturnRows(groupRows)

				groupUsersSQL := "SELECT * FROM `group_users` WHERE group_id = ?"
				groupUsersExpectQuery := mock.
					ExpectQuery(regexp.QuoteMeta(groupUsersSQL)).
					WithArgs(tt.gID)
				if tt.dbGroupUserErr != nil {
					groupUsersExpectQuery.WillReturnError(tt.dbGroupUserErr)
				} else {
					groupUserRows := sqlmock.
						NewRows([]string{"group_id", "user_id", "created_at"})
					for _, uID := range tt.want.UserIDs() {
						groupUserRows.AddRow(tt.gID, uID, now)
					}
					groupUsersExpectQuery.WillReturnRows(groupUserRows)
				}
			}

			r := &dbGroupRepository{db: db}
			got, err := r.Find(tt.gID)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("r.Find(%s)=_, %v; want _, %v", tt.gID, got, tt.want)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.Group{})); diff != "" {
					t.Errorf(
						"r.Find(%s)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.gID, got, tt.want, diff,
					)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
