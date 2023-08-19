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

func TestGroupFactory_Create(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		setup   func()
		want    *model.Group
		wantErr error
	}{
		{
			name: "Returns group",
			args: args{
				name: "TEST_GROUP_NAME",
			},
			setup: func() {
				uuid.SetRand(strings.NewReader("abcdefgh12345678"))
			},
			want: model.MustNewGroup(
				"61626364-6566-4768-b132-333435363738",
				"TEST_GROUP_NAME",
				[]model.UserID{},
			),
			wantErr: nil,
		},
		{
			name: "Error creating uuid",
			args: args{
				name: "TEST_GROUP_NAME",
			},
			setup: func() {
				uuid.SetRand(strings.NewReader("0"))
			},
			want:    nil,
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name: "Error invalid group name",
			args: args{
				name: "TEST_GROUP_NAME_XXXXXXXXXXXXXXXXXXXXXXXX",
			},
			setup: func() {
				uuid.SetRand(strings.NewReader("abcdefgh12345678"))
			},
			want:    nil,
			wantErr: model.ErrInvalidGroup,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			f := NewGroupFactory()
			got, err := f.Create(tt.args.name)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("f.Create(%s)=_, nil; want _, nil; want %s", tt.args.name, tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("f.Create(%s)=_, %t; want _, nil; want %#v", tt.args.name, err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Fatalf("want no error, but has error %v", err)
				}
				if diff := cmp.Diff(got, tt.want, cmp.AllowUnexported(model.Group{})); diff != "" {
					t.Errorf(
						"f.Create(%s)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
						tt.args.name, got, tt.want, diff,
					)
				}
			}
		})
	}
}
