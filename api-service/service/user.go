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

package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func (as *APIService) ListUsers() ([]model.User, error) {
	return as.Database.ListUsers()
}

func (as *APIService) GetUser(username string) (*model.User, error) {
	return as.Database.GetUser(username)
}

func (as *APIService) AddUser(user model.User) error {
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return err
	}

	user.Password, user.Salt = generateHashAndSalt(user.Password, salt, p)

	return as.Database.AddUser(user)
}

func (as *APIService) UpdateUserGroups(updatedUser model.User) error {
	return as.Database.UpdateUserGroups(updatedUser)
}

func (as *APIService) RemoveUser(username string) error {
	return as.Database.RemoveUser(username)
}

func generateHashAndSalt(password string, salt []byte, p *params) (string, string) {
	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, b64Salt
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
