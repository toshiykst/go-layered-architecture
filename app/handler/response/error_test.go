package response_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo"

	"github.com/toshiykst/go-layerd-architecture/app/handler/response"
)

func TestError(t *testing.T) {
	type args struct {
		code   response.ErrorCode
		status int
		err    error
	}
	tests := []struct {
		name     string
		args     args
		wantBody string
	}{
		{
			name: "Respond with http status in arg and error response body",
			args: args{
				code:   response.ErrorCodeInvalidArguments,
				status: http.StatusBadRequest,
				err:    errors.New("an error occurred"),
			},
			wantBody: `{"code":"INVALID_ARGUMENTS","status":400,"message":"an error occurred"}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(
				http.MethodGet,
				"https://example.com:8080/test",
				nil,
			)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			if err := response.Error(c, tt.args.code, tt.args.status, tt.args.err); err != nil {
				t.Fatalf("want no error, but has error: %s", err.Error())
			}

			res := rec.Result()

			wantStatusCode := tt.args.status
			gotStatusCode := res.StatusCode
			if gotStatusCode != wantStatusCode {
				t.Errorf("statusCode got = %d, want = %d", gotStatusCode, wantStatusCode)
			}

			resBody, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("Failed to read body: %s", err.Error())
			}
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Fatalf("Failed to close body: %s", err.Error())
				}
			}(res.Body)

			gotBody := string(resBody)
			if diff := cmp.Diff(gotBody, tt.wantBody); diff != "" {
				t.Errorf(
					"response body: got = %s, want = %s\ndiffers: (-got +want)\n%s",
					gotBody, tt.wantBody, diff,
				)
			}

		})
	}
}

func TestErrorInternal(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wantBody string
	}{
		{
			name: "Respond with http status internal server error and error response body",
			args: args{
				err: errors.New("an error occurred"),
			},
			wantBody: `{"code":"INTERNAL_SERVER_ERROR","status":500,"message":"an error occurred"}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(
				http.MethodGet,
				"https://example.com:8080/test",
				nil,
			)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			if err := response.ErrorInternal(c, tt.args.err); err != nil {
				t.Fatalf("want no error, but has error: %s", err.Error())
			}

			res := rec.Result()

			wantStatusCode := http.StatusInternalServerError
			gotStatusCode := res.StatusCode
			if gotStatusCode != wantStatusCode {
				t.Errorf("statusCode got = %d, want = %d", gotStatusCode, wantStatusCode)
			}

			resBody, err := io.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("Failed to read body: %s", err.Error())
			}
			defer func(Body io.ReadCloser) {
				if err := Body.Close(); err != nil {
					t.Fatalf("Failed to close body: %s", err.Error())
				}
			}(res.Body)

			gotBody := string(resBody)
			if diff := cmp.Diff(gotBody, tt.wantBody); diff != "" {
				t.Errorf(
					"response body: got = %s, want = %s\ndiffers: (-got +want)\n%s",
					gotBody, tt.wantBody, diff,
				)
			}

		})
	}
}
