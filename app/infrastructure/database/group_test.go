package database

import (
	"errors"
	"regexp"
	"strings"
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

func TestDatabase_dbGroupRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		group   *model.Group
		dbErr   error
		want    *model.Group
		wantErr error
	}{
		{
			name: "Creates a new group",
			group: model.NewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{},
			),
			dbErr: nil,
			want: model.NewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				[]model.UserID{},
			),
			wantErr: nil,
		},
		{
			name: "DB error",
			group: model.NewGroup(
				"TEST_GROUP_ID",
				"TEST_GROUP_NAME",
				nil,
			),
			dbErr:   errors.New("an error occurred"),
			want:    nil,
			wantErr: errors.New("an error occurred"),
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

			expectExec := mock.
				ExpectExec(regexp.QuoteMeta("INSERT INTO `groups` (`id`,`name`) VALUES (?,?)")).
				WithArgs(tt.group.ID(), tt.group.Name())

			if tt.wantErr != nil {
				expectExec.WillReturnError(tt.wantErr)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			r := &dbGroupRepository{db: db}
			got, err := r.Create(tt.group)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("r.Create(%v)=_, %v; want _, %v", tt.group, got, tt.want)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.Group{})); diff != "" {
					t.Errorf(
						"r.Create(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.group, got, tt.want, diff,
					)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_dbGroupRepository_AddUsers(t *testing.T) {
	type args struct {
		gID  model.GroupID
		uIDs []model.UserID
	}

	tests := []struct {
		name    string
		args    args
		dbErr   error
		wantErr error
	}{
		{
			name: "Add users to the group",
			args: args{
				gID: "TEST_GROUP_ID",
				uIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				},
			},
			dbErr:   nil,
			wantErr: nil,
		},
		{
			name: "DB group error",
			args: args{
				gID: "TEST_GROUP_ID",
				uIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_2",
					"TEST_USER_ID_3",
				},
			},
			dbErr:   errors.New("an error occurred"),
			wantErr: errors.New("an error occurred"),
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

			var (
				sqlArgs      []any
				placeHolders []string
			)
			for _, uID := range tt.args.uIDs {
				placeHolders = append(placeHolders, "(?,?)")
				sqlArgs = append(sqlArgs, tt.args.gID, uID)
			}

			expectExec := mock.
				ExpectExec(regexp.QuoteMeta("INSERT INTO `group_users` (`group_id`,`user_id`) VALUES " + strings.Join(placeHolders, ","))).
				WithArgs(testutil.ToDriverValues(t, sqlArgs...)...)

			if tt.wantErr != nil {
				expectExec.WillReturnError(tt.wantErr)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			r := &dbGroupRepository{db: db}
			err = r.AddUsers(tt.args.gID, tt.args.uIDs)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("r.AddUsers(%s, %v)=%v; want %v", tt.args.gID, tt.args.uIDs, err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
