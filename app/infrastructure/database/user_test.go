package database

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
	"github.com/toshiykst/go-layerd-architecture/app/testutil"
)

func TestDatabase_dbUserRepository_CreatePost(t *testing.T) {
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
