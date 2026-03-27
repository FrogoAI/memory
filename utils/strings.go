//nolint:revive // package name is intentional
package utils

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	mrand "math/rand/v2"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"github.com/mfonda/simhash"
	"github.com/twmb/murmur3"
	"golang.org/x/text/unicode/norm"
)

// Email parsing constants.
const (
	EmailTagStart = "+"
	EmailAt       = "@"

	singleKey        = 1
	pairedKeys       = 2
	shortIDRandBytes = 2
	tinyIDRandBytes  = 4
	uint32Bytes      = 4
	crc16Bytes       = 2
	base62Radix      = 62
	tinyIDTrimLen    = 5
)

// ABTest assigns data to a group bucket using murmur3 hashing with the given salt and group sizes.
func ABTest(data, salt []byte, groups ...uint64) uint64 {
	var total uint64
	for _, group := range groups {
		total += group
	}

	if total == 0 {
		return 0
	}

	return murmur3.Sum64(append(data, salt...)) % total
}

// SimHash returns the simhash fingerprint of the given data.
func SimHash(data []byte) uint64 {
	return simhash.Simhash(simhash.NewWordFeatureSet(data))
}

// SimHashCompare returns the Hamming distance between two simhash fingerprints.
func SimHashCompare(val1, val2 uint64) uint8 {
	return simhash.Compare(val1, val2)
}

// EmailUserName returns the local part of an email address (before the @).
func EmailUserName(email string) string {
	res := strings.Index(email, EmailAt)
	if res <= -1 {
		return email
	}

	return email[:res]
}

// EmailDomain returns the domain part of an email address (after the @).
func EmailDomain(email string) string {
	res := strings.Index(email, EmailAt)
	if res <= -1 {
		return email
	}

	res += len(EmailAt)

	return email[res:]
}

// SanitizeEmail lowercases the email and strips any +tag portion from the local part.
func SanitizeEmail(email string) string {
	return strings.Join(SplitBetweenTokens(strings.ToLower(email), EmailTagStart, EmailAt), EmailAt)
}

// NFDLowerString trims, NFD-normalizes, and lowercases the given string.
func NFDLowerString(str string) string {
	return strings.ToLower(norm.NFD.String(strings.TrimSpace(str)))
}

// CommonString strips all characters except letters, digits, and spaces from the string.
func CommonString(str string) string {
	var result strings.Builder

	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// SplitBetweenTokens takes a string and one or two tokens, and cuts
// everything between the two tokens (or two copies of the first token).
func SplitBetweenTokens(data string, keys ...string) []string {
	if data == "" {
		return []string{}
	}

	var key1, key2 string

	switch {
	case len(keys) == singleKey:
		key1 = keys[0]
		key2 = keys[0]
	case len(keys) >= pairedKeys:
		key1 = keys[0]
		key2 = keys[1]
	default:
		return []string{data}
	}

	if key1 == "" || key2 == "" {
		return []string{data}
	}

	s := strings.Index(data, key1)
	if s <= -1 {
		return []string{data}
	}

	part1 := data[0:s]

	s += len(key1)

	e := strings.Index(data[s:], key2)
	if e <= -1 {
		return []string{part1}
	}

	e += len(key2)

	part2 := data[s+e:]

	return []string{part1, part2}
}

// ByteSliceToString cast given bytes to string, without allocation memory
func ByteSliceToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b)) //nolint:gosec // zero-alloc []byte→string conversion via unsafe pointer cast
}

// GetShortID return short id
func GetShortID() ([]byte, error) {
	b := make([]byte, shortIDRandBytes)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	r := make([]byte, uint32Bytes)
	//nolint:gosec // Nanosecond() returns [0, 999999999], fits in uint32
	binary.BigEndian.PutUint32(r, uint32(time.Now().Nanosecond()))
	src := append(b, r...)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return dst, nil
}

// GetTinyID return tiny id
func GetTinyID() ([]byte, error) {
	b := make([]byte, tinyIDRandBytes)

	_, err := rand.Read(b) //nolint:gosec // crypto/rand.Read is correct here; suppresses deprecated-API warning
	if err != nil {
		return nil, err
	}

	r := make([]byte, uint32Bytes)
	//nolint:gosec // Nanosecond() returns [0, 999999999], fits in uint32
	binary.BigEndian.PutUint32(r, uint32(time.Now().Nanosecond()))
	b = append(b, r...)
	val := binary.BigEndian.Uint64(b)

	//nolint:gosec // uint64->int64 reinterprets high bit; big.Int handles negative values correctly
	return []byte(big.NewInt(int64(val)).Text(base62Radix))[tinyIDTrimLen:], nil
}

// Between function to get content between two keys
func Between(data string, keys ...string) string {
	var key1, key2 string

	switch {
	case len(keys) == singleKey:
		key1 = keys[0]
		key2 = keys[0]
	case len(keys) >= pairedKeys:
		key1 = keys[0]
		key2 = keys[1]
	default:
		return ""
	}

	if key1 == "" || key2 == "" {
		return ""
	}

	s := strings.Index(data, key1)
	if s <= -1 {
		return ""
	}

	s += len(key1)

	e := strings.Index(data[s:], key2)
	if e <= -1 {
		return ""
	}

	return strings.TrimSpace(data[s : s+e])
}

// SafeGet return value of pointer, and return default value if it nil
func SafeGet[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}

	return *ptr
}

// MaskField replaces the middle portion of str with asterisks, keeping the specified number of characters at each end.
func MaskField(str string, keepUnmaskedFront int, keepUnmaskedEnd int) string {
	var result strings.Builder

	size := len(str)
	defaultResult := strings.Repeat("*", size)

	if size <= (keepUnmaskedFront+keepUnmaskedEnd)*2 {
		return defaultResult
	}

	_, err := result.WriteString(str[:keepUnmaskedFront])
	if err != nil {
		return defaultResult
	}

	_, err = result.WriteString(strings.Repeat("*", size-keepUnmaskedFront-keepUnmaskedEnd))
	if err != nil {
		return defaultResult
	}

	_, err = result.WriteString(str[size-keepUnmaskedEnd:])
	if err != nil {
		return defaultResult
	}

	return result.String()
}

// SplitByChunks splits the string into chunks of the given size.
func SplitByChunks(s string, chunkSize int) []string {
	if chunkSize <= 0 {
		return nil
	}

	var chunks []string

	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}

		chunks = append(chunks, s[i:end])
	}

	return chunks
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

// RandStringBytes returns a random lowercase alphabetic string of length n.
func RandStringBytes(n int) string {
	if n <= 0 {
		return ""
	}

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mrand.IntN(len(letterBytes))] //nolint:gosec // math/rand is intentional — not security-sensitive
	}

	return string(b)
}

// HashName returns a short hash of the name: its first letter followed by a hex-encoded CRC-16.
func HashName(name string) string {
	name = strings.ToLower(name)
	val := crc16XModem([]byte(name))

	r := make([]byte, crc16Bytes)
	binary.BigEndian.PutUint16(r, val)

	return name[:1] + hex.EncodeToString(r)
}
