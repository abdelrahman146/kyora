package id

import (
	"crypto/rand"
	"io"
	"regexp"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/segmentio/ksuid"
)

func KsuidWithPrefix(prefix string) string {
	return prefix + ksuid.New().String()
}

func UlidWithPrefix(prefix string) string {
	return prefix + ulid.Make().String()
}

func Ksuid() string {
	return ksuid.New().String()
}

func Ulid() string {
	return ulid.Make().String()
}

func Base62(size int) string {
	return ulid.MustNew(ulid.Now(), ulid.DefaultEntropy()).String()[0:size]
}

func Base62WithPrefix(prefix string, size int) string {
	return prefix + Base62(size)
}

func RandomNumber(length int) (string, error) {
	const table = "1234567890"
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

func RandomString(length int) (string, error) {
	const table = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

var nonAlphaNum = regexp.MustCompile(`[^A-Za-z0-9]+`)

// NewCodeFromString generates a code of specified length from the input string.
// It removes non-alphanumeric characters, uppercases the result, and pads with 'X' if needed.
func NewCodeFromString(s string, length int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return strings.Repeat("X", length)
	}
	// Remove non-alphanumeric and split by spaces after cleanup
	cleaned := nonAlphaNum.ReplaceAllString(s, " ")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return strings.Repeat("X", length)
	}
	parts := strings.Fields(cleaned)
	// Prefer first token; fall back to entire cleaned string
	token := parts[0]
	token = nonAlphaNum.ReplaceAllString(token, "")
	if len(token) >= length {
		return strings.ToUpper(token[:length])
	}
	// Pad to length using X if shorter
	token = strings.ToUpper(token)
	if len(token) < length {
		token = strings.Repeat("X", length-len(token)) + token
	}
	return token
}

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
