// Copyright (c) 2023 Sorint.lab S.p.A.
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

package dto

import (
	"errors"
	"fmt"
	"regexp"

	effectivebytes "github.com/ercole-io/ercole/v2/utils/effective_bytes"
)

type OracleExadataMeasurement struct {
	UnparsedValue string  `json:"unparsedValue"`
	Symbol        string  `json:"symbol"`
	Quantity      float64 `json:"quantity"`
}

func ToOracleExadataMeasurement(s string) (res *OracleExadataMeasurement, err error) {
	res = &OracleExadataMeasurement{}

	if s == "0" {
		s = "0B"
	}

	res.UnparsedValue = s

	res.Quantity, err = effectivebytes.Float64(s)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`[a-zA-Z]+`)
	match := re.FindStringSubmatch(s)

	if len(match) > 1 {
		return nil, errors.New("invalid OracleExadataMeasurement")
	}

	res.Symbol = match[0]

	return res, nil
}

func (m *OracleExadataMeasurement) ToTb() (*OracleExadataMeasurement, error) {
	b, err := effectivebytes.Parse(m.UnparsedValue)
	if err != nil {
		return nil, err
	}

	convertedValue := b.Format(effectivebytes.Format, "TB", false)

	q, err := effectivebytes.Float64(convertedValue)
	if err != nil {
		return nil, err
	}

	return &OracleExadataMeasurement{
		UnparsedValue: b.String(),
		Symbol:        "TB",
		Quantity:      q,
	}, nil
}

func (m *OracleExadataMeasurement) SetUnparsedValue() {
	m.UnparsedValue = fmt.Sprintf("%v%s", m.Quantity, m.Symbol)
}
