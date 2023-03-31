package strcase

import (
	"testing"
)

func toCamel(tb testing.TB) {
	cases := [][]string{
		{"test_case", "TestCase"},
		{"test.case", "TestCase"},
		{"test", "Test"},
		{"TestCase", "TestCase"},
		{" test  case ", "TestCase"},
		{"", ""},
		{"many_many_words", "ManyManyWords"},
		{"AnyKind of_string", "AnyKindOfString"},
		{"odd-fix", "OddFix"},
		{"numbers2And55with000", "Numbers2And55With000"},
		{"ID", "Id"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := ToCamel(in)
		if result != out {
			tb.Errorf("%q (%q != %q)", in, result, out)
		}
	}
}

func TestToCamel(t *testing.T) {
	toCamel(t)
}

func BenchmarkToCamel(b *testing.B) {
	benchmarkCamelTest(b, toCamel)
}

func toLowerCamel(tb testing.TB) {
	cases := [][]string{
		{"foo-bar", "fooBar"},
		{"TestCase", "testCase"},
		{"", ""},
		{"AnyKind of_string", "anyKindOfString"},
		{"AnyKind.of-string", "anyKindOfString"},
		{"ID", "id"},
		{"some string", "someString"},
		{" some string", "someString"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := ToLowerCamel(in)
		if result != out {
			tb.Errorf("%q (%q != %q)", in, result, out)
		}
	}
}

func TestToLowerCamel(t *testing.T) {
	toLowerCamel(t)
}

func BenchmarkToLowerCamel(b *testing.B) {
	benchmarkCamelTest(b, toLowerCamel)
}

func benchmarkCamelTest(b *testing.B, fn func(testing.TB)) {
	for n := 0; n < b.N; n++ {
		fn(b)
	}
}
