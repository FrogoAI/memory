//nolint:revive // package name is intentional
package utils

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/FrogoAI/testutils"
	"github.com/twmb/murmur3"
)

func murmur3Sum64(data, salt []byte) uint64 {
	return murmur3.Sum64(append(data, salt...))
}

const hexConst = 16

func TestSimHash(t *testing.T) {
	testcases := []struct {
		name    string
		input   []byte
		input2  []byte
		output  string
		compare uint8
	}{
		{
			input:   []byte{},
			output:  "ffffffffffffffff",
			input2:  []byte{},
			compare: 0,
		},
		{
			input:   []byte("maksym"),
			output:  "cdb0cb6c28c1bafb",
			input2:  []byte("maksim"),
			compare: 10,
		},
		{
			input:   []byte("maxim"),
			output:  "8ff2a4602ca09ad7",
			input2:  []byte("maksym"),
			compare: 20,
		},
		{
			input:   []byte("test string"),
			output:  "9d2ffffe9fdfffef",
			input2:  []byte("another string"),
			compare: 19,
		},
		{
			input:   []byte("sim string"),
			output:  "d9bffdde6b7fd79e",
			input2:  []byte("sime string"),
			compare: 15,
		},
		{
			input:   []byte("different string"),
			output:  "f9efddfe6b77d68e",
			input2:  []byte("similar number"),
			compare: 23,
		},
		{
			input:   []byte("token 1"),
			output:  "bf63ff6dbe05f7fe",
			input2:  []byte("token 2"),
			compare: 2,
		},
	}

	for _, test := range testcases {
		t.Run(string(test.input)+" vs "+string(test.input2), func(t *testing.T) {
			val := SimHash(test.input)
			val2 := SimHash(test.input2)

			fmt.Println(CRC16(strconv.FormatUint(val, hexConst)))
			fmt.Println(CRC16(strconv.FormatUint(val2, hexConst)))

			testutils.Equal(t,
				strconv.FormatUint(val, hexConst),
				test.output,
			)
			testutils.Equal(t, SimHashCompare(val, val2), test.compare)
		})
	}
}

func TestSplitBetweenTokens(t *testing.T) {
	testcases := []struct {
		name      string
		data      string
		arguments []string
		result    []string
	}{
		{
			name:      "split_between_different_tokens",
			data:      `some_string_which_we_should_split@should_not_be_visible;must_be_present`,
			arguments: []string{"@", ";"},
			result:    []string{"some_string_which_we_should_split", "must_be_present"},
		},
		{
			name:      "split_between_same_token_tokens",
			data:      `some_string_which_we_should_split;should_not_be_visible;must_be_present`,
			arguments: []string{";", ";"},
			result:    []string{"some_string_which_we_should_split", "must_be_present"},
		},
		{
			name:      "split_between_single_token",
			data:      `some_string_which_we_should_split;should_not_be_visible;must_be_present`,
			arguments: []string{";"},
			result:    []string{"some_string_which_we_should_split", "must_be_present"},
		},
		{
			name:      "return_fist_part_for_single_token",
			data:      `some_string_which_we_should_split;should_not_be_visible`,
			arguments: []string{";"},
			result:    []string{"some_string_which_we_should_split"},
		},
		{
			name:      "return_income_string_if_no_arguments",
			data:      `some_string_which_we_should_split;should_be_also_visible`,
			arguments: []string{},
			result:    []string{"some_string_which_we_should_split;should_be_also_visible"},
		},
		{
			name:      "return_income_string_if_no_match",
			data:      `some_string_which_we_should_split;should_be_also_visible`,
			arguments: []string{"@"},
			result:    []string{"some_string_which_we_should_split;should_be_also_visible"},
		},
		{
			name:      "return_empty_for_empty_input",
			data:      ``,
			arguments: []string{"@"},
			result:    []string{},
		},
		{
			name:      "if_both_token_are_empty",
			data:      `some_string_which_we_should_split`,
			arguments: []string{"", ""},
			result:    []string{"some_string_which_we_should_split"},
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			result := SplitBetweenTokens(test.data, test.arguments...)
			testutils.Equal(t, result, test.result)
		})
	}
}

func TestSanitizeEmail(t *testing.T) {
	testcases := []struct {
		name   string
		email  string
		result string
	}{
		{
			name:   "email_with_tag",
			email:  "testemail+example@gmail.com",
			result: `testemail@gmail.com`,
		},
		{
			name:   "email_with_two_tags",
			email:  "testemail+exa+mple@gmail.com",
			result: `testemail@gmail.com`,
		},
		{
			name:   "user_name_with_tag_without_domain",
			email:  "testemail+exa",
			result: `testemail`,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			result := SanitizeEmail(test.email)
			testutils.Equal(t, result, test.result)
		})
	}
}

func TestEmailDomain(t *testing.T) {
	testcases := []struct {
		name   string
		email  string
		domain string
	}{
		{
			name:   "valid_email",
			email:  "testemail@gmail.com",
			domain: `gmail.com`,
		},
		{
			name:   "empty_email",
			email:  "",
			domain: ``,
		},
		{
			name:   "only_domain",
			email:  "@gmail.com",
			domain: `gmail.com`,
		},
		{
			name:   "only_username",
			email:  "testmail@",
			domain: ``,
		},
		{
			name:   "at_not_present",
			email:  "test;mail.com",
			domain: `test;mail.com`,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			domain := EmailDomain(test.email)
			testutils.Equal(t, domain, test.domain)
		})
	}
}

func TestSplitByChunks(t *testing.T) {
	chunks := SplitByChunks("teststring", 3)
	testutils.Equal(t, chunks, []string{"tes", "tst", "rin", "g"})
}

func TestGetTinyID(t *testing.T) {
	shortID, err := GetTinyID()
	testutils.Equal(t, err, nil)
	shortID2, err := GetTinyID()
	testutils.Equal(t, err, nil)
	testutils.NotEqual(t, shortID, shortID2)
}

func TestRandStringBytes(t *testing.T) {
	val := RandStringBytes(1)
	testutils.Equal(t, len(val), 1)
}

func TestHashName(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		result string
	}{
		{
			name:   "test",
			input:  "test",
			result: "t9b06",
		},
		{
			name:   "weavers",
			input:  "weavers",
			result: "w0709",
		},
		{
			name:   "InsideGallery",
			input:  "InsideGallery",
			result: "i53a5",
		},
	}

	for _, tst := range testcases {
		t.Run(tst.name, func(t *testing.T) {
			val := HashName(tst.input)
			testutils.Equal(t, val, tst.result)
		})
	}
}

func TestByteSliceToString(t *testing.T) {
	testcases := map[string]struct {
		bytes []byte
		out   string
	}{
		"inStrs:empty":     {bytes: []byte{}, out: ""},
		"inStrs:nil":       {bytes: nil, out: ""},
		"inStrs:non_empty": {bytes: []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33}, out: "Hello world!"},
	}

	for k, c := range testcases {
		c := c

		t.Run(k, func(t *testing.T) {
			str := ByteSliceToString(c.bytes)
			testutils.Equal(t, str, c.out)
		})
	}
}

func TestByteSliceToStringNative(t *testing.T) {
	testcases := map[string]struct {
		bytes []byte
		out   string
	}{
		"inStrs:empty":     {bytes: []byte{}, out: ""},
		"inStrs:nil":       {bytes: nil, out: ""},
		"inStrs:non_empty": {bytes: []byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33}, out: "Hello world!"},
	}

	for k, c := range testcases {
		c := c

		t.Run(k, func(t *testing.T) {
			str := ByteSliceToString(c.bytes)
			testutils.Equal(t, str, c.out)
		})
	}
}

func TestBetween(t *testing.T) {
	testcases := map[string]struct {
		data   string
		keys   []string
		result string
	}{
		"empty_data": {
			data:   "",
			keys:   []string{"[RESULT]"},
			result: "",
		},
		"not_key_in_data": {
			data:   "Some text",
			keys:   []string{"[RESULT]"},
			result: "",
		},
		"single_key_in_data": {
			data:   "Some [RESULT]text",
			keys:   []string{"[RESULT]"},
			result: "",
		},
		"key_is_present": {
			data:   "Some [RESULT]text[RESULT]",
			keys:   []string{"[RESULT]"},
			result: "text",
		},
		"trip_space_in_result": {
			data:   "Some [RESULT] text \n[RESULT]",
			keys:   []string{"[RESULT]"},
			result: "text",
		},
		"no_key": {
			data:   "Some [RESULT] text \n[RESULT]",
			keys:   []string{},
			result: "",
		},
		"between_two_keys": {
			data:   "Some < text\n > \n",
			keys:   []string{"<", ">"},
			result: "text",
		},
	}

	for name, test := range testcases {
		test := test

		t.Run(name, func(t *testing.T) {
			result := Between(test.data, test.keys...)
			testutils.Equal(t, test.result, result)
		})
	}
}

func TestABTest(t *testing.T) {
	cases := []struct {
		name   string
		data   []byte
		salt   []byte
		groups []uint64
		want   uint64
	}{
		{
			name:   "zero_total_returns_zero",
			data:   []byte("user1"),
			salt:   []byte("salt"),
			groups: []uint64{0, 0},
			want:   0,
		},
		{
			name:   "empty_groups_returns_zero",
			data:   []byte("user1"),
			salt:   []byte("salt"),
			groups: []uint64{},
			want:   0,
		},
		{
			name:   "single_group",
			data:   []byte("user1"),
			salt:   []byte("salt"),
			groups: []uint64{100},
			want:   murmur3Sum64([]byte("user1"), []byte("salt")) % 100,
		},
		{
			name:   "two_groups",
			data:   []byte("test"),
			salt:   []byte("s"),
			groups: []uint64{50, 50},
			want:   murmur3Sum64([]byte("test"), []byte("s")) % 100,
		},
		{
			name:   "deterministic",
			data:   []byte("user1"),
			salt:   []byte("salt"),
			groups: []uint64{50, 50},
			want:   ABTest([]byte("user1"), []byte("salt"), 50, 50),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ABTest(tc.data, tc.salt, tc.groups...)
			testutils.Equal(t, result, tc.want)
		})
	}
}

func TestEmailUserName(t *testing.T) {
	cases := []struct {
		name  string
		email string
		want  string
	}{
		{name: "valid_email", email: "user@example.com", want: "user"},
		{name: "no_at_sign", email: "useronly", want: "useronly"},
		{name: "empty_string", email: "", want: ""},
		{name: "at_sign_only", email: "@domain.com", want: ""},
		{name: "multiple_at", email: "user@sub@domain.com", want: "user"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, EmailUserName(tc.email), tc.want)
		})
	}
}

func TestNFDLowerString(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty", input: "", want: ""},
		{name: "already_lower", input: "hello", want: "hello"},
		{name: "upper_case", input: "HELLO", want: "hello"},
		{name: "mixed_case", input: "HeLLo WoRLd", want: "hello world"},
		{name: "leading_trailing_spaces", input: "  hello  ", want: "hello"},
		{name: "unicode_accented", input: "\u00C9", want: "e\u0301"},
		{name: "unicode_nfd_decomposition", input: "\u00F1", want: "n\u0303"},
		{name: "tabs_and_newlines", input: "\t hello \n", want: "hello"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, NFDLowerString(tc.input), tc.want)
		})
	}
}

func TestCommonString(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty", input: "", want: ""},
		{name: "only_letters", input: "hello", want: "hello"},
		{name: "letters_and_digits", input: "abc123", want: "abc123"},
		{name: "with_spaces", input: "hello world", want: "hello world"},
		{name: "strips_punctuation", input: "hello, world!", want: "hello world"},
		{name: "strips_special_chars", input: "a@b#c$d%e", want: "abcde"},
		{name: "unicode_letters", input: "café!", want: "café"},
		{name: "only_special_chars", input: "!@#$%^&*()", want: ""},
		{name: "mixed_content", input: "user-name_123 (test)", want: "username123 test"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, CommonString(tc.input), tc.want)
		})
	}
}

func TestGetShortID(t *testing.T) {
	id1, err1 := GetShortID()
	testutils.Equal(t, err1, nil)

	id2, err2 := GetShortID()
	testutils.Equal(t, err2, nil)

	// IDs should be non-empty
	testutils.NotEqual(t, len(id1), 0)
	testutils.NotEqual(t, len(id2), 0)

	// hex-encoded 6 bytes = 12 characters
	testutils.Equal(t, len(id1), 12)
	testutils.Equal(t, len(id2), 12)
}

func TestSafeGet(t *testing.T) {
	cases := []struct {
		name         string
		ptr          *int
		defaultValue int
		want         int
	}{
		{name: "nil_returns_default", ptr: nil, defaultValue: 42, want: 42},
		{name: "non_nil_returns_value", ptr: intPtr(99), defaultValue: 42, want: 99},
		{name: "zero_value_ptr", ptr: intPtr(0), defaultValue: 42, want: 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, SafeGet(tc.ptr, tc.defaultValue), tc.want)
		})
	}
}

func intPtr(v int) *int { return &v }

func TestSafeGetString(t *testing.T) {
	var nilStr *string

	testutils.Equal(t, SafeGet(nilStr, "default"), "default")

	val := "hello"
	testutils.Equal(t, SafeGet(&val, "default"), "hello")
}

func TestRandStringBytesExtended(t *testing.T) {
	cases := []struct {
		name string
		n    int
		want int
	}{
		{name: "zero_length", n: 0, want: 0},
		{name: "negative_length", n: -1, want: 0},
		{name: "one", n: 1, want: 1},
		{name: "ten", n: 10, want: 10},
		{name: "hundred", n: 100, want: 100},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := RandStringBytes(tc.n)
			testutils.Equal(t, len(result), tc.want)

			// verify all characters are lowercase letters
			for _, c := range result {
				if c < 'a' || c > 'z' {
					t.Fatalf("unexpected character: %c", c)
				}
			}
		})
	}
}

func TestSplitByChunksExtended(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		chunkSize int
		want      []string
	}{
		{name: "empty_string", input: "", chunkSize: 3, want: nil},
		{name: "chunk_size_zero", input: "abc", chunkSize: 0, want: nil},
		{name: "chunk_size_negative", input: "abc", chunkSize: -1, want: nil},
		{name: "exact_fit", input: "abcdef", chunkSize: 3, want: []string{"abc", "def"}},
		{name: "remainder", input: "abcdefg", chunkSize: 3, want: []string{"abc", "def", "g"}},
		{name: "chunk_larger_than_input", input: "ab", chunkSize: 10, want: []string{"ab"}},
		{name: "single_char_chunks", input: "abc", chunkSize: 1, want: []string{"a", "b", "c"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.Equal(t, SplitByChunks(tc.input, tc.chunkSize), tc.want)
		})
	}
}

func TestMaskField(t *testing.T) {
	testcases := []struct {
		name              string
		str               string
		keepUnmaskedFront int
		keepUnmaskedEnd   int
		expected          string
	}{
		{
			name:              "mask long string",
			str:               "secret-string-here",
			keepUnmaskedFront: 2,
			keepUnmaskedEnd:   3,
			expected:          "se*************ere",
		},
		{
			name:              "mask short string",
			str:               "sec",
			keepUnmaskedFront: 2,
			keepUnmaskedEnd:   3,
			expected:          "***",
		},
		{
			name:              "mask medium string",
			str:               "secret",
			keepUnmaskedFront: 2,
			keepUnmaskedEnd:   3,
			expected:          "******",
		},
		{
			name:              "mask minimum to have show string",
			str:               "secret12345",
			keepUnmaskedFront: 2,
			keepUnmaskedEnd:   3,
			expected:          "se******345",
		},
		{
			name:              "mask without keep",
			str:               "secret12345",
			keepUnmaskedFront: 0,
			keepUnmaskedEnd:   0,
			expected:          "***********",
		},
		{
			name:              "mask front only",
			str:               "sec",
			keepUnmaskedFront: 1,
			keepUnmaskedEnd:   0,
			expected:          "s**",
		},
		{
			name:              "mask back only",
			str:               "sec",
			keepUnmaskedFront: 0,
			keepUnmaskedEnd:   1,
			expected:          "**c",
		},
		{
			name:              "empty str",
			str:               "",
			keepUnmaskedFront: 2,
			keepUnmaskedEnd:   3,
			expected:          "",
		},
	}

	for _, test := range testcases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result := MaskField(test.str, test.keepUnmaskedFront, test.keepUnmaskedEnd)
			testutils.Equal(t, result, test.expected)
		})
	}
}
