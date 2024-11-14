// Copyright (c) 2022 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	mathRand "math/rand"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength     = 16
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
	minSpecialChar = 2
	minNum         = 2
	minUpperCase   = 2
	passwordLength = 16
	usernameLength = 16
)

var sample = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func GenerateRandomBytes() ([]byte, error) {
	b := make([]byte, saltLength)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateHashAndSalt(password string, salt []byte) (string, string) {
	p := Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  saltLength,
		KeyLength:   32,
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, b64Salt
}

func SuggestPassword() string {
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

	var password strings.Builder

	for i := 0; i < minSpecialChar; i++ {
		random := r.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	for i := 0; i < minNum; i++ {
		random := r.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	for i := 0; i < minUpperCase; i++ {
		random := r.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := r.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}

	inRune := []rune(password.String())
	r.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}

func SuggestUsername() string {
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

	b := make([]rune, usernameLength)
	for i := range b {
		b[i] = sample[r.Intn(len(sample))]
	}

	return string(b)
}
