package removesingletonarrays

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveSingletonArrays(t *testing.T) {
	cases := []struct {
		test string
		want string
	}{
		{
			`{"a":"b"}`,
			`{"a":"b"}`,
		},
		{
			`{"a":"1","b":["2"]}`,
			`{"a":"1","b":"2"}`,
		},
		{
			`{"path":[{"secret/foo":[{"capabilities":["read"]}]}]}`,
			`{"path":{"secret/foo":{"capabilities":"read"}}}`,
		},
		{
			`{"arr": ["some val", {"k": ["v"]}],"map": {"nested": {"k": ["v"]}}}`,
			`{"arr":["some val",{"k":"v"}],"map":{"nested":{"k":"v"}}}`,
		},
	}

	for i, test := range cases {
		returnedValue, err := RemoveSingletonArrays(test.test)

		raw := json.RawMessage(returnedValue)

		if err != nil {
			t.Errorf("%d: failed to remove singletons: %v", i+1, err)
		}

		assert.Equal(t, string(raw), test.want, "The JSON isnt as expected")
	}
}

func TestRemoveSingletonArrayErrors(t *testing.T) {
	cases := []struct {
		got  string
		want error
	}{
		{
			got:  `{:}`,
			want: errors.New("invalid character ':' looking for beginning of object key string"),
		},
	}

	for i, test := range cases {
		returnedValue, err := RemoveSingletonArrays(test.got)

		if err == nil {
			t.Errorf("%d: expected error, did not get error: %v", i+1, returnedValue)
		}

		assert.EqualError(t, err, test.want.Error())
	}
}

func TestWithIgnores(t *testing.T) {
	cases := []struct {
		got       string
		want      string
		keyIgnore []string
	}{
		{
			got:       `{"a":"b"}`,
			want:      `{"a":"b"}`,
			keyIgnore: []string{"1"},
		},
		{
			got:       `{"ignored-singleton":["b"]}`,
			want:      `{"ignored-singleton":"b"}`,
			keyIgnore: []string{""},
		},
		{
			got:       `{"ignored-singleton":["b"]}`,
			want:      `{"ignored-singleton":["b"]}`,
			keyIgnore: []string{"ignored-singleton"},
		},
		{
			got:       `{"ignored-singleton":["b"]}`,
			want:      `{"ignored-singleton":"b"}`,
			keyIgnore: []string{"not-ignored-singleton"},
		},
		{
			got:       `{"path":[{"secret/foo":[{"capabilities":["read"]}]}]}`,
			want:      `{"path":{"secret/foo":{"capabilities":"read"}}}`,
			keyIgnore: []string{},
		},
		{
			got:       `{"path":[{"secret/foo":[{"capabilities":["read"]}]}]}`,
			want:      `{"path":{"secret/foo":{"capabilities":["read"]}}}`,
			keyIgnore: []string{"capabilities"},
		},
	}

	for i, test := range cases {
		returnedValue, err := WithIgnores(test.got, test.keyIgnore)

		raw := json.RawMessage(returnedValue)

		if err != nil {
			t.Errorf("%d: failed to remove singletons: %v", i+1, err)
		}

		assert.Equal(t, string(raw), test.want, "The JSON isnt as expected")
	}
}

func TestWithIgnoresErrors(t *testing.T) {
	cases := []struct {
		got       string
		want      error
		keyIgnore []string
	}{
		{
			got:       `{:}`,
			want:      errors.New("invalid character ':' looking for beginning of object key string"),
			keyIgnore: []string{"1"},
		},
	}

	for i, test := range cases {
		returnedValue, err := WithIgnores(test.got, test.keyIgnore)

		if err == nil {
			t.Errorf("%d: expected error, did not get error: %v", i+1, returnedValue)
		}

		assert.EqualError(t, err, test.want.Error())
	}
}
