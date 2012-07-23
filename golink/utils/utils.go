package utils

import (
    "crypto/sha1"
    "fmt"
)

// hash a string
func PasswordHash(pwd string) string {
    hasher := sha1.New()
    hasher.Write([]byte(pwd))
    return fmt.Sprintf("%x", hasher.Sum(nil))
}
