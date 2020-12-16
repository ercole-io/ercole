// Copyright (c) 2020 Sorint.lab S.p.A.
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

package controller

import (
	"errors"

	"github.com/ercole-io/ercole/v2/utils"
)

//go:generate mockgen -source ../service/service.go -destination=fake_service_test.go -package=controller

//Common data
var errMock error = errors.New("MockError")
var aerrMock utils.AdvancedErrorInterface = utils.NewAdvancedErrorPtr(errMock, "mock")

type FailingReader struct{}

func (fr *FailingReader) Read(p []byte) (n int, err error) {
	return 0, aerrMock
}

// type patchingFunctionSimilarTo struct{ pf model.PatchingFunction }

// func (sp *patchingFunctionSimilarTo) Matches(x interface{}) bool {
// 	fmt.Printf("%#v\n%#v\n", sp.pf, x)
// 	if val, ok := x.(model.PatchingFunction); !ok {
// 		return false
// 	} else if val.Code != sp.pf.Code {
// 		return false
// 	} else {
// 		return reflect.DeepEqual(sp.pf.Vars, val.Vars)
// 	}
// }

// func (sp *patchingFunctionSimilarTo) String() string {
// 	return fmt.Sprintf("is similar to %v", sp.pf)
// }
