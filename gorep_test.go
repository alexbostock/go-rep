package main

import "testing"

func TestTest(t *testing.T) {
	cases := []struct {
		regex string
		str   string
		want  bool
	}{}

	for _, c := range cases {
		sm, err := Compile(c.regex)
		if err != nil {
			t.Errorf(c.regex, c.str, "compiler error")
		}

		test := Test(sm, c.str)
		if test != c.want {
			t.Errorf(c.regex, c.str, "match wrong")
		}
	}
}
