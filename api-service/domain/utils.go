// Copyright (c) 2024 Sorint.lab S.p.A.
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

package domain

func ToUpperLevelLayers[model interface{}, dto interface{}](m []model, f func(model) (*dto, error)) ([]dto, error) {
	res := make([]dto, 0, len(m))

	for _, v := range m {
		k, err := f(v)
		if err != nil {
			return nil, err
		}

		res = append(res, *k)
	}

	return res, nil
}
