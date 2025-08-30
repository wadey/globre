package globre

import (
	"regexp"
	"testing"
)

func TestConvert(t *testing.T) {
	var convertTests = []struct {
		in  string
		out string
		yes []string
		no  []string
	}{
		{"*.foo.com", `^.*\.foo\.com$`,
			[]string{"bar.foo.com", ".foo.com"},
			[]string{"foo.com"}},
		{"foo{,-dev*,-prod}.bar.com", `^foo(|-dev.*|-prod)\.bar\.com$`,
			[]string{"foo.bar.com", "foo-dev.bar.com", "foo-prod.bar.com", "foo-dev1.bar.com"},
			[]string{"foo-prod1.bar.com"}},
		{"foo{,-bar}{,-dev}.baz.com", `^foo(|-bar)(|-dev)\.baz\.com$`,
			[]string{"foo.baz.com", "foo-bar.baz.com", "foo-dev.baz.com", "foo-bar-dev.baz.com"},
			[]string{"foo-dev-bar.baz.com"}},
		{"foo{,{-a,-b}}.baz.com", `^foo(|(-a|-b))\.baz\.com$`,
			[]string{"foo.baz.com", "foo-a.baz.com", "foo-b.baz.com"},
			[]string{"foo-a-b.baz.com"}},
	}
	for _, tt := range convertTests {
		t.Run(tt.in, func(t *testing.T) {
			s, err := Convert(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			// t.Logf("%#q", s)
			if tt.out != s {
				t.Fatalf("got %#q, want %#q", s, tt.out)
			}

			re, err := regexp.Compile(s)
			if err != nil {
				t.Fatal(err)
			}

			for _, y := range tt.yes {
				if !re.MatchString(y) {
					t.Fatalf("want match, but did not: %v", y)
				}
			}
			for _, n := range tt.no {
				if re.MatchString(n) {
					t.Fatalf("do not want match, but did: %v", n)
				}
			}
		})
	}
}

func TestConvertSeparators(t *testing.T) {
	var convertTests = []struct {
		in  string
		sep string
		out string
		yes []string
		no  []string
	}{
		{"*.foo.com", ".", `^[^\.]*\.foo\.com$`,
			[]string{"bar.foo.com", ".foo.com"},
			[]string{"foo.foo.foo.com"}},
		{"**.foo.com", ".", `^.*\.foo\.com$`,
			[]string{"foo.foo.foo.com"},
			[]string{}},
		{"/foo/*/baz", "/", `^/foo/[^/]*/baz$`,
			[]string{"/foo/bar/baz"},
			[]string{"/foo/bar/bar/baz"}},
		{"/foo/**/baz", "/", `^/foo/.*/baz$`,
			[]string{"/foo/bar/bar/baz"},
			[]string{}},
		{"/**", "/", `^/.*$`,
			[]string{"/foo/bar/bar/baz"},
			[]string{""}},
	}
	for _, tt := range convertTests {
		t.Run(tt.in, func(t *testing.T) {
			s, err := ConvertSeparators(tt.in, tt.sep)
			if err != nil {
				t.Fatal(err)
			}
			// t.Logf("%#q", s)
			if tt.out != s {
				t.Fatalf("got %#q, want %#q", s, tt.out)
			}

			re, err := regexp.Compile(s)
			if err != nil {
				t.Fatal(err)
			}

			for _, y := range tt.yes {
				if !re.MatchString(y) {
					t.Fatalf("want match, but did not: %v", y)
				}
			}
			for _, n := range tt.no {
				if re.MatchString(n) {
					t.Fatalf("do not want match, but did: %v", n)
				}
			}
		})
	}
}
