package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RandUUID 产生随机UUID
func RandUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

// RandStr 产生随机字符串
func RandStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

// RandDigitStr 产生随机数字字符串
func RandDigitStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []byte("0123456789")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

// RandStr 产生随机16进制数字字符串
func RandHEXStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []byte("0123456789abcdef")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}
