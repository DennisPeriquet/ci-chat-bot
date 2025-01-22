package slack

import (
	"maps"
	"testing"
)

func TestBuildJobParams(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		params      string
		expected    map[string]string
		errorString string
	}{
		{
			name:        "NoParameters",
			params:      "",
			expected:    map[string]string{},
			errorString: "",
		},
		{
			name:        "SimpleParameter",
			params:      "\"KEY1=VALUE1\"",
			expected:    map[string]string{"KEY1": "VALUE1"},
			errorString: "",
		},
		{
			name:        "IncorrectlyQuotedParameter",
			params:      "“KEY1=VALUE1”",
			expected:    map[string]string{"KEY1": "VALUE1"},
			errorString: "",
		},
		{
			name:        "IncorrectlyDeliminatedParameter",
			params:      "\"KEY1:VALUE1\"",
			expected:    nil,
			errorString: "unable to interpret `KEY1:VALUE1` as a parameter. Please ensure that all parameters are in the form of KEY=VALUE; nested parameters should be delimited with \\n",
		},
		{
			name:        "MarkDownLinkParameter",
			params:      "\"KEY1=<http://abc123.com|VALUE1>\"",
			expected:    map[string]string{"KEY1": "VALUE1"},
			errorString: "",
		},
		{
			name:        "One NestedParameter",
			params:      "\"DEVSCRIPTS_CONFIG=KEY1=VALUE1\"",
			expected:    map[string]string{"DEVSCRIPTS_CONFIG": "KEY1=VALUE1"},
			errorString: "",
		},
		{
			name:        "Two NestedParameters",
			params:      "\"DEVSCRIPTS_CONFIG=KEY1=VALUE1\\nKEY2=VALUE2\"",
			expected:    map[string]string{"DEVSCRIPTS_CONFIG": "KEY1=VALUE1\\nKEY2=VALUE2"},
			errorString: "",
		},
		{
			name:        "NestedParametersBad",
			params:      "\"DEVSCRIPTS_CONFIG=KEY1=VALUE1,KEY2=VALUE2\"",
			expected:    nil,
			errorString: "unable to interpret `DEVSCRIPTS_CONFIG=KEY1=VALUE1,KEY2=VALUE2` as a DEVSCRIPTS_CONFIG parameter. Please ensure that nested parameters are delimited by newlines",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BuildJobParams(tc.params)
			if !maps.Equal(got, tc.expected) {
				t.Errorf("got %q, expected %q", got, tc.expected)
			}
			if err != nil && err.Error() != tc.errorString {
				t.Errorf("got error %q, expected error %q", err, tc.errorString)
			}
		})
	}
}

func TestParseParameterValue(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "NoValue",
			value:    "",
			expected: "",
		},
		{
			name:     "SimpleValue",
			value:    "VALUE1",
			expected: "VALUE1",
		},
		{
			name:     "MarkDownLinkValue",
			value:    "<http://abc123.com|VALUE1>",
			expected: "VALUE1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseParameterValue(tc.value)
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}
