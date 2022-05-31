package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

// passsalt 密码加密盐
var authAlgorithm = "pbkdf2_sha256"
var authIterations = 15000
var authDigest = sha256.New
var authSaltLen = 12
var authHashLen = 32

func GetPasswd(passwd string) string {
	return encryptPasswd(passwd, salt(authSaltLen), authIterations)
}

func encryptPasswd(passwd string, passsalt []byte, iter int) string {
	dk := pbkdf2.Key([]byte(passwd), passsalt, iter, authHashLen, authDigest)
	hash := base64.StdEncoding.EncodeToString(dk)
	return fmt.Sprintf("%s$%d$%s$%s", authAlgorithm, authIterations, passsalt, hash)
}

func salt(n int) []byte {
	rand.Seed(time.Now().UnixNano())
	letters := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func VerifyPasswd(encryptpasswd, passwd string) bool {
	var result bool
	passlist := strings.Split(encryptpasswd, "$")
	if len(passlist) != 4 {
		return result
	}

	salt := passlist[2]
	iter, err := strconv.Atoi(passlist[1])
	if err != nil {
		return result
	}
	// TODO: select algorithm
	encryptingpasswd := encryptPasswd(passwd, []byte(salt), iter)
	if encryptpasswd == encryptingpasswd {
		result = true
	}
	return result
}

func CheckGoodPassword(password string) error {
	if len(password) < 8 {
		return errors.New("密码不得小于8位")
	}
	return nil
}
