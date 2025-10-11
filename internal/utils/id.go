package utils

import (
	"crypto/rand"
	"io"

	"github.com/oklog/ulid/v2"
	"github.com/segmentio/ksuid"
)

type idHelper struct{}

func (idHelper) NewKsuidWithPrefix(prefix string) string {
	return prefix + ksuid.New().String()
}

func (idHelper) NewUlidWithPrefix(prefix string) string {
	return prefix + ulid.Make().String()
}

func (idHelper) NewKsuid() string {
	return ksuid.New().String()
}

func (idHelper) NewUlid() string {
	return ulid.Make().String()
}

func (idHelper) NewBase62(size int) string {
	return ulid.MustNew(ulid.Now(), ulid.DefaultEntropy()).String()[0:size]
}

func (h idHelper) NewBase62WithPrefix(prefix string, size int) string {
	return prefix + h.NewBase62(size)
}

func (idHelper) RandomNumber(length int) (string, error) {
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

func (idHelper) RandomString(length int) (string, error) {
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

var ID = idHelper{}
