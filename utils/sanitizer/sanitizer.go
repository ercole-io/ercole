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

package sanitizer

import (
	"fmt"
	"reflect"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/ercole-io/ercole/v2/logger"
)

type Sanitizer struct {
	log logger.Logger

	policy *bluemonday.Policy
}

func NewSanitizer(log logger.Logger) *Sanitizer {
	return &Sanitizer{
		log:    log,
		policy: bluemonday.StrictPolicy(),
	}
}

func (s *Sanitizer) Sanitize(obj interface{}) (interface{}, error) {
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	if err := s.sanitizeRecursive(copy, original); err != nil {
		return nil, err
	}

	return copy.Interface(), nil
}

func (s *Sanitizer) sanitizeRecursive(copy, original reflect.Value) error {
	switch original.Kind() {

	case reflect.Ptr:
		// Unwrap and call once again

		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return nil // pointer is nil
		}

		copy.Set(reflect.New(originalValue.Type()))
		if err := s.sanitizeRecursive(copy.Elem(), originalValue); err != nil {
			return err
		}

	case reflect.Interface:
		originalValue := original.Elem()

		if !originalValue.IsValid() {
			return nil
		}

		copyValue := reflect.New(originalValue.Type()).Elem()
		if err := s.sanitizeRecursive(copyValue, originalValue); err != nil {
			return err
		}
		copy.Set(copyValue)

	case reflect.Struct:

		if original.CanInterface() {
			i := original.Interface()

			// If it is a time.Time, copy directly entire struct
			if _, ok := i.(time.Time); ok {
				copy.Set(original)
				return nil
			}
		}

		for i := 0; i < original.NumField(); i += 1 {
			if err := s.sanitizeRecursive(copy.Field(i), original.Field(i)); err != nil {
				return err
			}
		}

	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			if err := s.sanitizeRecursive(copy.Index(i), original.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()

			if err := s.sanitizeRecursive(copyValue, originalValue); err != nil {
				return err
			}
			copy.SetMapIndex(key, copyValue)
		}

	case reflect.String:
		sanitizedString := s.sanitizeString(original.Interface().(string))
		copy.SetString(sanitizedString)

	default:
		if !copy.CanSet() {
			return fmt.Errorf("Can't set value")
		}
		copy.Set(original)
	}

	return nil
}

func (s *Sanitizer) sanitizeString(a string) string {
	b := s.policy.Sanitize(a)

	if a != b {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(a),
			FromFile: "Original",

			B:      difflib.SplitLines(b),
			ToFile: "Sanitized",

			Context: 3,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)

		s.log.Warnf("A string has been sanitized:\n%s", text)
	}

	return b
}
