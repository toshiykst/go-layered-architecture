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

func TestDatabase_dbGroupRepository_List(t *testing.T) {
	tests := []struct {
		name              string
		groups            model.Groups
		want              model.Groups
		wantErr           error
		wantGroupsSQL     string
		wantGroupUsersSQL string
		dbGroupsErr       error
		dbGroupUsersErr   error
	}{
		{
			name: "Returns groups",
			groups: model.Groups{
				model.NewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{"TEST_USER_ID_1"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{"TEST_USER_ID_2"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]model.UserID{"TEST_USER_ID_3"},
				),
			},
			want: model.Groups{
				model.NewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{"TEST_USER_ID_1"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{"TEST_USER_ID_2"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]model.UserID{"TEST_USER_ID_3"},
				),
			},
			wantGroupsSQL:     "SELECT * FROM `groups`",
			wantGroupUsersSQL: "SELECT * FROM `group_users` WHERE group_id IN (?,?,?)",
			wantErr:           nil,
			dbGroupsErr:       nil,
			dbGroupUsersErr:   nil,
		},
		{
			name:              "Groups not found",
			groups:            nil,
			want:              nil,
			wantGroupsSQL:     "SELECT * FROM `groups`",
			wantGroupUsersSQL: "",
			wantErr:           nil,
			dbGroupsErr:       nil,
			dbGroupUsersErr:   nil,
		},
		{
			name:              "DB group error",
			groups:            nil,
			want:              nil,
			wantGroupsSQL:     "SELECT * FROM `groups`",
			wantGroupUsersSQL: "",
			wantErr:           errors.New("an error occurred"),
			dbGroupsErr:       errors.New("an error occurred"),
			dbGroupUsersErr:   nil,
		},
		{
			name: "DB groupusers error",
			groups: model.Groups{
				model.NewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{"TEST_USER_ID_1"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{"TEST_USER_ID_2"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]model.UserID{"TEST_USER_ID_3"},
				),
			},
			want: model.Groups{
				model.NewGroup(
					"TEST_GROUP_ID_1",
					"TEST_GROUP_NAME_1",
					[]model.UserID{"TEST_USER_ID_1"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_2",
					"TEST_GROUP_NAME_2",
					[]model.UserID{"TEST_USER_ID_2"},
				),
				model.NewGroup(
					"TEST_GROUP_ID_3",
					"TEST_GROUP_NAME_3",
					[]model.UserID{"TEST_USER_ID_3"},
				),
			},
			wantGroupsSQL:     "SELECT * FROM `groups`",
			wantGroupUsersSQL: "SELECT * FROM `group_users` WHERE group_id IN (?,?,?)",
			wantErr:           errors.New("an error occurred"),
			dbGroupsErr:       nil,
			dbGroupUsersErr:   errors.New("an error occurred"),
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

			groupsExpectQuery := mock.ExpectQuery(regexp.QuoteMeta(tt.wantGroupsSQL))

			if tt.dbGroupsErr != nil {
				groupsExpectQuery.WillReturnError(tt.dbGroupsErr)
			} else {
				now := time.Now()
				groupRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"})
				for _, g := range tt.groups {
					groupRows.AddRow(g.ID(), g.Name(), now, now)
				}
				groupsExpectQuery.WillReturnRows(groupRows)

				if tt.wantGroupUsersSQL != "" {
					groupUsersExpectQuery := mock.
						ExpectQuery(regexp.QuoteMeta(tt.wantGroupUsersSQL)).
						WithArgs(testutil.ToDriverValues[model.GroupID](t, tt.want.IDs()...)...)
					if tt.dbGroupUsersErr != nil {
						groupUsersExpectQuery.WillReturnError(tt.dbGroupUsersErr)
					} else {
						groupUserRows := sqlmock.
							NewRows([]string{"group_id", "user_id", "created_at"})
						for _, g := range tt.groups {
							for _, uID := range g.UserIDs() {
								groupUserRows.AddRow(g.ID(), uID, now)
							}
						}
						groupUsersExpectQuery.WillReturnRows(groupUserRows)
					}
				}
			}

			r := &dbGroupRepository{db: db}
			got, err := r.List()
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("want an error, but has no error")
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("r.List()=_, %v; want _, %v", got, tt.want)
				}
			} else {
				if err != nil {
					t.Fatalf("want no err, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.Group{})); diff != "" {
					t.Errorf(
						"r.List()=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						got, tt.want, diff,
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

func TestDatabase_dbGroupRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		group   *model.Group
		wantErr error
	}{
		{
			name:    "Updates a group",
			group:   model.NewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{}),
			wantErr: nil,
		},
		{
			name:    "Error",
			group:   model.NewGroup("TEST_GROUP_ID", "TEST_GROUP_NAME", []model.UserID{}),
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
				ExpectExec(regexp.QuoteMeta("UPDATE `groups` SET `name`=? WHERE `id` = ?")).
				WithArgs(tt.group.Name(), tt.group.ID())

			if tt.wantErr != nil {
				expectExec.WillReturnError(tt.wantErr)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			r := &dbGroupRepository{db: db}
			err = r.Update(tt.group)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, but has error %v", err)
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
