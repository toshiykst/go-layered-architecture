package factory

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

func TestUserFactory_Create(t *testing.T) {
	type args struct {
		name  string
		email string
	}
	tests := []struct {
		name    string
		args    args
		setup   func()
		want    *model.User
		wantErr error
	}{
		{
			name: "Returns group",
			args: args{
				name:  "TEST_GROUP_NAME",
				email: "TEST_USER_EMAIL",
			},
			setup: func() {
				uuid.SetRand(strings.NewReader("abcdefgh12345678"))
			},
			want: model.MustNewUser(
				"61626364-6566-4768-b132-333435363738",
				"TEST_GROUP_NAME",
				"TEST_USER_EMAIL",
			),
			wantErr: nil,
		},
		{
			name: "Error creating uuid",
			args: args{
				name:  "TEST_USER_NAME",
				email: "TEST_USER_EMAIL",
			},
			setup: func() {
				uuid.SetRand(strings.NewReader("0"))
			},
			want:    nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name: "Error invalid user input",
			args: args{
				name:  "TEST_USER_NAME_XXXXXXXXXXXXXXXXXXXXXXXXX",
				email: "TEST_USER_EMAIL",
			},
			setup: func() {
				uuid.SetRand(strings.NewReader("abcdefgh12345678"))
			},
			want:    nil,
			wantErr: model.ErrInvalidUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			f := NewUserFactory()
			got, err := f.Create(tt.args.name, tt.args.email)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf(
						"f.Create(%s, %s)=_, nil; want _, %v",
						tt.args.name, tt.args.email, tt.wantErr,
					)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"f.Create(%s, %s)=_, %v; want _, %v",
						tt.args.name, tt.args.email, err, tt.wantErr,
					)
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.User{})); diff != "" {
					t.Errorf(
						"f.Create(%s, %s)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.args.name, tt.args.email, got, tt.want, diff,
					)
				}
			}
		})
	}
}
