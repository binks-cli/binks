package agent

import "testing"

func TestIsAIQuery(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", false},
		{"no prefix", "hello world", false},
		{"only prefix", ">>", false},
		{"prefix with space only", ">>   ", false},
		{"prefix with text", ">> ask something", true},
		{"prefix with leading space", "   >> ask", true},
		{"prefix with no space after", ">>ask", true},
		{"prefix with tabs", "\t>>\tquestion", true},
		{"accidental prefix in middle", "echo >> not ai", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsAIQuery(tc.input)
			if got != tc.want {
				t.Errorf("IsAIQuery(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
