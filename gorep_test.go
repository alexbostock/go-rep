package main

import "testing"

func TestTest(t *testing.T) {
	cases := []struct {
		regex string
		str   string
		want  bool
	}{
		{"a", "a", true},
		{"abc", "abc", true},
		{"abcde", "abcde", true},
		{"abc", "abcd", false},
		{"abcde", "abcd", false},
		{"a?", "a", true},
		{"a?b", "b", true},
		{"a?b", "ab", true},
		{"a?a", "a", true},

		// TODO: pass these
		//{"a?a", "aa", true},
		//{"a?ab", "ab", true},
		//{"a?ab", "aab", true},
		//{".+a", "sdga", true},

		{"a+", "a", true},
		{"a+", "aa", true},
		{"a+", "aaa", true},
		{"a+", "aaab", false},
		{"a+", "b", false},
		{"a+", "baaa", false},

		{".", "a", true},
		{".+", "adgsfhfjyf", true},
		{".+a", "sdgb", false},
		{"...", "abc", true},
	}

	for _, c := range cases {
		sm, err := Compile(c.regex)
		if err != nil {
			t.Errorf(c.regex, c.str, "compiler error")
		}

		test := Test(sm, c.str)
		if test != c.want {
			t.Errorf(c.regex, c.str)
		}
	}
}
