package misc

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

//MD5Hash ...
func MD5Hash(s string) string {
	checksum := md5.Sum([]byte(s))
	return hex.EncodeToString(checksum[:])
}

func cryptMd5(plainText string, salt []byte) string {
	srcBytes := []byte(plainText)

	saltedBytes := make([]byte, 0, len(salt)+len(srcBytes))
	saltedBytes = append(saltedBytes, srcBytes...)
	saltedBytes = append(saltedBytes, salt...)

	hashBytes := md5.Sum(saltedBytes)
	hashWithSaltBytes := make([]byte, 0, len(hashBytes)+len(salt))
	hashWithSaltBytes = append(hashWithSaltBytes, hashBytes[:]...)
	hashWithSaltBytes = append(hashWithSaltBytes, salt...)
	return base64.StdEncoding.EncodeToString(hashWithSaltBytes)
}

//CryptWithMd5 ...
func CryptWithMd5(s, salt string) string {
	return cryptMd5(s, []byte(salt))
}

//VerifyHashWithMd5 ...
func VerifyHashWithMd5(plainText, hash, salt string) bool {
	hashWithSaltBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}

	hashSizeInBytes := 128 / 8
	hashBytesLen := len(hashWithSaltBytes)
	if hashBytesLen < hashSizeInBytes {
		return false
	}

	return hash == cryptMd5(plainText, []byte(salt))
}
