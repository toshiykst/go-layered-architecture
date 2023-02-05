package factory

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

func TestUserFactory_Create(t *testing.T) {
	uID := model.UserID("61626364-6566-4768-b132-333435363738")
	name := "TEST_USER_NAME"
	email := "TEST_USER_EMAIL"
	uuid.SetRand(strings.NewReader("abcdefgh12345678"))

	want := model.NewUser(uID, name, email)
	f := NewUserFactory()

	got, err := f.Create(name, email)
	if err != nil {
		t.Fatalf("want no error, but has error %v", err)
	}
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(model.User{})); diff != "" {
		t.Errorf(
			"f.Create(%s, %s)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
			name, email, got, want, diff,
		)
	}
}

func TestUserFactory_Create_Error_UUID(t *testing.T) {
	wantErr := errors.New("unexpected EOF")

	name := "TEST_USER_NAME"
	email := "TEST_USER_EMAIL"
	uuid.SetRand(strings.NewReader("0"))

	f := NewUserFactory()

	_, err := f.Create(name, email)
	if err == nil {
		t.Fatalf("f.Create(%s, %s)=_, nil; want _, nil; want %s", name, email, wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("f.Create(%s, %s)=_, %v; want _, nil; want %s", name, email, err, wantErr)
	}
}
