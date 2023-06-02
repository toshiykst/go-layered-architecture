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
	"github.com/toshiykst/go-layerd-architecture/app/domain/repository"
	"github.com/toshiykst/go-layerd-architecture/app/testutil"
)

func TestDatabase_dbUserRepository_Find(t *testing.T) {
	tests := []struct {
		name    string
		uID     model.UserID
		want    *model.User
		wantErr error
		dbErr   error
	}{
		{
			name:    "Returns a user",
			uID:     model.UserID("TEST_USER_ID"),
			want:    model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"),
			wantErr: nil,
			dbErr:   nil,
		},
		{
			name:    "Not found",
			uID:     model.UserID("TEST_USER_ID"),
			want:    nil,
			wantErr: nil,
			dbErr:   gorm.ErrRecordNotFound,
		},
		{
			name:    "Error",
			uID:     model.UserID("TEST_USER_ID"),
			want:    nil,
			wantErr: errors.New("an error occurred"),
			dbErr:   errors.New("an error occurred"),
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

			sql := "SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1"
			expectQuery := mock.
				ExpectQuery(regexp.QuoteMeta(sql)).
				WithArgs(tt.uID)

			if tt.dbErr != nil {
				expectQuery.WillReturnError(tt.dbErr)
			} else {
				now := time.Now()
				rows := sqlmock.
					NewRows([]string{"id", "name", "email", "created_at", "updated_at"}).
					AddRow(tt.want.ID(), tt.want.Name(), tt.want.Email(), now, now)
				expectQuery.WillReturnRows(rows)
			}

			r := &dbUserRepository{db: db}
			got, err := r.Find(tt.uID)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("r.Find(%s)=_, %v; want _, %v", tt.uID, got, tt.want)
				}
			} else {
				if err != nil {
					t.Fatalf("want no errpr, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.User{})); diff != "" {
					t.Errorf(
						"r.Find(%s)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.uID, got, tt.want, diff,
					)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_dbUserRepository_List(t *testing.T) {
	tests := []struct {
		name    string
		filter  repository.UserListFilter
		want    model.Users
		wantSQL string
		wantErr error
		dbErr   error
	}{
		{
			name:   "Returns users",
			filter: repository.UserListFilter{},
			want: model.Users{
				model.NewUser(
					"TEST_USER_ID_1",
					"TEST_USER_NAME_1",
					"TEST_USER_EMAIL_1",
				),
				model.NewUser(
					"TEST_USER_ID_2",
					"TEST_USER_NAME_2",
					"TEST_USER_EMAIL_2",
				),
				model.NewUser(
					"TEST_USER_ID_3",
					"TEST_USER_NAME_3",
					"TEST_USER_EMAIL_3",
				),
			},
			wantSQL: "SELECT * FROM `users`",
			wantErr: nil,
			dbErr:   nil,
		},
		{
			name: "Returns users filtered by user_ids",
			filter: repository.UserListFilter{
				UserIDs: []model.UserID{
					"TEST_USER_ID_1",
					"TEST_USER_ID_3",
				},
			},
			want: model.Users{
				model.NewUser(
					"TEST_USER_ID_1",
					"TEST_USER_NAME_1",
					"TEST_USER_EMAIL_1",
				),
				model.NewUser(
					"TEST_USER_ID_3",
					"TEST_USER_NAME_3",
					"TEST_USER_EMAIL_3",
				),
			},
			wantSQL: "SELECT * FROM `users` WHERE id IN (?,?)",
			wantErr: nil,
			dbErr:   nil,
		},
		{
			name:    "Error",
			want:    nil,
			wantSQL: "SELECT * FROM `users`",
			wantErr: errors.New("an error occurred"),
			dbErr:   errors.New("an error occurred"),
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

			expectQuery := mock.
				ExpectQuery(regexp.QuoteMeta(tt.wantSQL))

			if len(tt.filter.UserIDs) > 0 {
				expectQuery.WithArgs(testutil.ToDriverValues[model.UserID](t, tt.filter.UserIDs...)...)
			}

			if tt.dbErr != nil {
				expectQuery.WillReturnError(tt.dbErr)
			} else {
				now := time.Now()
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at"})
				for _, u := range tt.want {
					rows.AddRow(u.ID(), u.Name(), u.Email(), now, now)
				}
				expectQuery.WillReturnRows(rows)
			}

			r := &dbUserRepository{db: db}
			got, err := r.List(tt.filter)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("r.List(%v)=_, %v; want _, %v", tt.filter, got, tt.want)
				}
			} else {
				if err != nil {
					t.Fatalf("want no errpr, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.User{})); diff != "" {
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

func TestDatabase_dbUserRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		want    *model.User
		wantErr error
	}{
		{
			name:    "Creates a new user",
			user:    model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"),
			want:    model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"),
			wantErr: nil,
		},
		{
			name:    "Error",
			user:    model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"),
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
				ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`id`,`name`,`email`) VALUES (?,?,?)")).
				WithArgs(tt.user.ID(), tt.user.Name(), tt.user.Email())

			if tt.wantErr != nil {
				expectExec.WillReturnError(tt.wantErr)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			r := &dbUserRepository{db: db}
			got, err := r.Create(tt.user)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("r.Create(%v)=_, %v; want _, %v", tt.user, got, tt.want)
				}
			} else {
				if err != nil {
					t.Fatalf("want no errpr, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.User{})); diff != "" {
					t.Errorf(
						"r.Create(%v)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.user, got, tt.want, diff,
					)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_dbUserRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		wantErr error
	}{
		{
			name:    "Creates a new user",
			user:    model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"),
			wantErr: nil,
		},
		{
			name:    "Error",
			user:    model.NewUser("TEST_USER_ID", "TEST_USER_NAME", "TEST_USER_EMAIL"),
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
				ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `name`=?,`email`=? WHERE `id`=?")).
				WithArgs(tt.user.Name(), tt.user.Email(), tt.user.ID())

			if tt.wantErr != nil {
				expectExec.WillReturnError(tt.wantErr)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			r := &dbUserRepository{db: db}
			err = r.Update(tt.user)
			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
			} else {
				if err != nil {
					t.Fatalf("want no errpr, but has error %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDatabase_dbUserRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		userID  model.UserID
		wantErr error
	}{
		{
			name:    "Delete a user",
			userID:  model.UserID("TEST_USER_ID"),
			wantErr: nil,
		},
		{
			name:    "Error",
			userID:  model.UserID("TEST_USER_ID"),
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
				ExpectExec(regexp.QuoteMeta("DELETE FROM `users` WHERE `users`.`id` = ?")).
				WithArgs(tt.userID)

			if tt.wantErr != nil {
				expectExec.WillReturnError(tt.wantErr)
			} else {
				expectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			r := &dbUserRepository{db: db}

			err = r.Delete(tt.userID)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("want an error, but has no error")
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("r.Delete(%s)=%v; want %v", tt.userID, err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Fatalf("want no errpr, but has error %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
