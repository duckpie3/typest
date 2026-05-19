package typing

import "testing"

func TestCountCharacters(t *testing.T) {
	cases := []struct {
		name  string
		words []string
		want  int
	}{
		{name: "empty", words: []string{}, want: 0},
		{name: "single", words: []string{"hello"}, want: 5},
		{name: "two-words", words: []string{"hello", "world"}, want: 11},
		{name: "double-space", words: []string{"a", "", "b"}, want: 4},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := countCharacters(tc.words); got != tc.want {
				t.Fatalf("countCharacters(%v) = %d, want %d", tc.words, got, tc.want)
			}
		})
	}
}
