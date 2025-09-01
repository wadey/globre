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
		{"foo[{}]", `^foo[\{\}]$`,
			[]string{"foo{"}, []string{"foo("}},
		{"foo[*]", `^foo[\*]$`,
			[]string{"foo*"}, []string{"foo("}},
		// https://github.com/gobwas/glob/issues/62
		{"{,a}{,a}a", `^(|a)(|a)a$`,
			[]string{"a", "aa", "aaa"},
			[]string{"aaaa", ""}},
		// https://github.com/gobwas/glob/issues/66
		{"{**/daxing,daxing}/**/*dev*.yaml", `^(.*.*/daxing|daxing)/.*.*/.*dev.*\.yaml$`,
			[]string{"playground/daxing/generated/dev.yaml"},
			[]string{}},
		// https://github.com/gobwas/glob/issues/61
		{"start*art", `^start.*art$`,
			[]string{"start-art"},
			[]string{"start"}},
		// https://github.com/gobwas/glob/issues/50
		{"{google.*,*yandex:*.exe:page.*}", `^(google\..*|.*yandex:.*\.exe:page\..*)$`,
			[]string{"yandex:service.exe:page.12345"},
			[]string{}},
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

func TestConvertErr(t *testing.T) {
	var convertTests = []struct {
		in  string
		err string
	}{
		{"[.foo.com", "error parsing regexp: missing closing ]: `[\\.foo\\.com$`"},
		{"{).foo.com", "error parsing regexp: missing closing ): `^(\\)\\.foo\\.com$`"},
	}
	for _, tt := range convertTests {
		t.Run(tt.in, func(t *testing.T) {
			_, err := Compile(tt.in)
			if err == nil {
				t.Fatalf("did not error: %#q", tt.in)
			}
			if tt.err != err.Error() {
				t.Fatalf("got %s, want %s", err, tt.err)
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
		// https://github.com/gobwas/glob/issues/66
		{"{**/daxing,daxing}/**/*dev*.yaml", "/", `^(.*/daxing|daxing)/.*/[^/]*dev[^/]*\.yaml$`,
			[]string{"playground/daxing/generated/dev.yaml"},
			[]string{}},
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
