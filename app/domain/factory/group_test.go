package factory

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/toshiykst/go-layerd-architecture/app/domain/model"
)

func TestGroupFactory_Create(t *testing.T) {
	gID := model.GroupID("61626364-6566-4768-b132-333435363738")
	name := "TEST_GROUP_NAME"

	want := model.NewGroup(gID, name, []model.UserID{})

	uuid.SetRand(strings.NewReader("abcdefgh12345678"))
	f := NewGroupFactory()

	got, err := f.Create(name)
	if err != nil {
		t.Fatalf("want no error, but has error %v", err)
	}
	if diff := cmp.Diff(got, want, cmp.AllowUnexported(model.Group{})); diff != "" {
		t.Errorf(
			"f.Create(%s)=%v, nil; want %v, nil\ndiffers: (-got +want)\n%s",
			name, got, want, diff,
		)
	}
}

func TestGroupFactory_Create_Error_UUID(t *testing.T) {
	wantErr := errors.New("unexpected EOF")

	name := "TEST_GROUP_NAME"
	uuid.SetRand(strings.NewReader("0"))

	f := NewGroupFactory()

	_, err := f.Create(name)
	if err == nil {
		t.Fatalf("f.Create(%s)=_, nil; want _, nil; want %s", name, wantErr)
	}
	if err.Error() != wantErr.Error() {
		t.Errorf("f.Create(%s)=_, %v; want _, nil; want %s", name, err, wantErr)
	}
}
