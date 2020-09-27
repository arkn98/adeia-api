/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"adeia/pkg/util/constants"

	"github.com/alexedwards/argon2id"
	"github.com/dchest/uniuri"
	"github.com/trustelem/zxcvbn"
)

// GenerateRandomBytes generates cryptographically secure random bytes for a
// specified size (in bytes).
func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}

// EncodeHex encodes the given byte slice into hex string.
func EncodeHex(b []byte) string {
	return hex.EncodeToString(b)
}

// EncodeBase64 encodes the given byte slice into a url-safe base64 string (without padding).
func EncodeBase64(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// DecodeBase64 decodes the given url-safe base64 string (without padding) into a byte slice.
func DecodeBase64(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// Hash hashes the give byte slice using SHA256.
func Hash(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}

// HashPassword uses argon2id to generate a hash from the password.
func HashPassword(p string) (hash string, err error) {
	return argon2id.CreateHash(p, argon2id.DefaultParams)
}

// ComparePwdHash compares the password and hash.
func ComparePwdHash(p, h string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(p, h)
}

// NewEmpID generates a short, URL-friendly alpha-numeric employee ID.
func NewEmpID() string {
	return uniuri.NewLenChars(constants.EmployeeIDLength, []byte(constants.EmployeeIDChars))
}

// PasswordStrength returns the strength of a password (on a scale of 0 - 4).
func PasswordStrength(password string) int {
	// TODO: add site-specific, user-specific inputs to penalize weak passwords
	return zxcvbn.PasswordStrength(password, []string{}).Score
}
