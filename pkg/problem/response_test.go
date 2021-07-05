package problem

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"
)

func TestIsHTTPResponse(t *testing.T) {
	var testCases = []struct {
		name   string
		header http.Header
		want   bool
	}{
		{
			name:   "nil header",
			header: nil,
			want:   false,
		},
		{
			name:   "empty header",
			header: nil,
			want:   false,
		},
		{
			name:   "problem content-type",
			header: newHeader("Content-Type", HTTPContentType),
			want:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				Header: tc.header,
			}
			got := IsHTTPResponse(resp)
			if got != tc.want {
				t.Errorf("wanted %t, got: %t", tc.want, got)
			}
		})
	}
}

func newHeader(key, value string) http.Header {
	h := make(http.Header)
	h.Add(key, value)
	return h
}

func TestParseHTTPResponse_fail(t *testing.T) {
	var testErr = errors.New("test err")
	var testCases = []struct {
		name   string
		reader io.ReadCloser
		errIs  error
		errAs  error
	}{
		{
			name:   "read",
			reader: io.NopCloser(iotest.ErrReader(testErr)),
			errIs:  testErr,
		},
		{
			name:   "close",
			reader: errCloser{strings.NewReader(""), testErr},
			errIs:  testErr,
		},
		{
			name:   "parse",
			reader: io.NopCloser(strings.NewReader("???")),
			errAs:  &json.SyntaxError{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var response = &http.Response{
				Body: tc.reader,
			}
			_, err := ParseHTTPResponse(response)
			if err == nil {
				t.Error("wanted error, got nil")
			}
			if tc.errIs != nil {
				if !errors.Is(err, tc.errIs) {
					t.Errorf("wanted: %s; got: %s", tc.errIs, err)
				}
			} else {
				if !errors.As(err, &tc.errAs) {
					t.Errorf("wanted: %T; got: %s", tc.errAs, err)
				}
			}
		})
	}
}

type errCloser struct {
	reader io.Reader
	err    error
}

func (e errCloser) Read(p []byte) (n int, err error) { return e.reader.Read(p) }
func (e errCloser) Close() error                     { return e.err }
