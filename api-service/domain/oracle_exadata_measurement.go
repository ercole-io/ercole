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

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/ercole-io/ercole/v2/utils"
	effectivebytes "github.com/ercole-io/ercole/v2/utils/effective_bytes"
)

type OracleExadataMeasurement struct {
	unparsedValue string
	Symbol        string
	Quantity      float64
}

func (m OracleExadataMeasurement) String() string {
	return fmt.Sprintf("%.2f %s", m.Quantity, m.Symbol)
}

func (m OracleExadataMeasurement) Human(symbol string) (string, error) {
	c, err := Convert(m, symbol)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2f %s", c.Quantity, c.Symbol), nil
}

func (m *OracleExadataMeasurement) Add(qty float64, symbol string) {
	t := *m

	safeOem, err := Convert(t, symbol)
	if err != nil {
		return
	}

	safeOem.Quantity += qty

	og, err := Convert(*safeOem, t.Symbol)
	if err != nil {
		return
	}

	m.Quantity = og.Quantity
	m.Symbol = og.Symbol
}

func (m *OracleExadataMeasurement) Sub(oem OracleExadataMeasurement) {
	if m == nil {
		return
	}

	safeOem, err := Convert(oem, m.Symbol)
	if err != nil {
		return
	}

	m.Quantity -= safeOem.Quantity
}

func NewOracleExadataMeasurement() *OracleExadataMeasurement {
	return &OracleExadataMeasurement{
		Symbol: "MIB",
	}
}

func Convert(m OracleExadataMeasurement, symbol string) (*OracleExadataMeasurement, error) {
	if m.String() == "" {
		return nil, errors.New("invalid OracleExadataMeasurement, cannot convert to TiB")
	}

	b, err := effectivebytes.Parse(m.String())
	if err != nil {
		return nil, err
	}

	convertedValue := b.Format(effectivebytes.Format, symbol, false)

	q, err := effectivebytes.Float64(convertedValue)
	if err != nil {
		return nil, err
	}

	trimmedvalue := utils.TruncateFloat64(q)

	return &OracleExadataMeasurement{
		unparsedValue: m.String(),
		Symbol:        symbol,
		Quantity:      trimmedvalue,
	}, nil
}

func IntToOracleExadataMeasurement(d int, symbol string) (res *OracleExadataMeasurement, err error) {
	s := strconv.Itoa(d)

	if s == "0" {
		s = "0B"
	} else {
		s = fmt.Sprintf("%d%s", d, symbol)
	}

	res = &OracleExadataMeasurement{unparsedValue: s}

	f, err := effectivebytes.Float64(s)
	if err != nil {
		return nil, err
	}

	res.Quantity = utils.TruncateFloat64(f)

	re := regexp.MustCompile(`[a-zA-Z]+`)
	match := re.FindStringSubmatch(s)

	if len(match) > 1 {
		return nil, errors.New("invalid OracleExadataMeasurement")
	}

	res.Symbol = match[0]

	return res, nil
}

func StringToOracleExadataMeasurement(s string) (res *OracleExadataMeasurement, err error) {
	if s == "0" || s == "" {
		s = "0B"
	}

	res = &OracleExadataMeasurement{unparsedValue: s}

	f, err := effectivebytes.Float64(s)
	if err != nil {
		return nil, err
	}

	res.Quantity = utils.TruncateFloat64(f)

	re := regexp.MustCompile(`[a-zA-Z]+`)
	match := re.FindStringSubmatch(s)

	if len(match) > 1 {
		return nil, errors.New("invalid OracleExadataMeasurement")
	}

	res.Symbol = match[0]

	return res, nil
}

func GetPercentage(measure, total OracleExadataMeasurement) string {
	sameSymbol := "MIB"

	safemeasure, err := Convert(measure, sameSymbol)
	if err != nil {
		return ""
	}

	safetotal, err := Convert(total, sameSymbol)
	if err != nil {
		return ""
	}

	if safetotal.Quantity != 0 {
		perc := (safemeasure.Quantity * 100) / safetotal.Quantity

		return fmt.Sprintf("%.2f%%", perc)
	}

	return fmt.Sprintf("0%%")
}
