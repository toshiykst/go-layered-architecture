package response

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo"
)

func TestOK(t *testing.T) {
	tests := []struct {
		name     string
		body     any
		wantBody string
	}{
		{
			name: "Respond with http status OK and arg body",
			body: struct {
				Message string `json:"message"`
			}{
				Message: "TEST_MESSAGE",
			},
			wantBody: `{"message":"TEST_MESSAGE"}
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

			if err := OK(c, tt.body); err != nil {
				t.Fatalf("want no err, but has error: %s", err.Error())
			}

			res := rec.Result()

			wantStatusCode := http.StatusOK
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

func TestCreated(t *testing.T) {
	tests := []struct {
		name     string
		body     any
		wantBody string
	}{
		{
			name: "Respond with http status created and arg body",
			body: struct {
				Message string `json:"message"`
			}{
				Message: "TEST_MESSAGE",
			},
			wantBody: `{"message":"TEST_MESSAGE"}
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

			if err := Created(c, tt.body); err != nil {
				t.Fatalf("want no err, but has error: %s", err.Error())
			}

			res := rec.Result()

			wantStatusCode := http.StatusCreated
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
