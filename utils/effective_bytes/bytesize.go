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

package effectivebytes

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type ByteSize uint64

const (
	B  ByteSize = 1
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
)

var longUnitMap = map[ByteSize]string{
	B:  "byte",
	KB: "kilobyte",
	MB: "megabyte",
	GB: "gigabyte",
	TB: "terabyte",
}

var shortUnitMap = map[ByteSize]string{
	B:  "B",
	KB: "KB",
	MB: "MB",
	GB: "GB",
	TB: "TB",
}

var unitMap = map[string]ByteSize{
	"B":     B,
	"BYTE":  B,
	"BYTES": B,

	"K":         KB,
	"KB":        KB,
	"KILOBYTE":  KB,
	"KILOBYTES": KB,

	"M":         MB,
	"MB":        MB,
	"MEGABYTE":  MB,
	"MEGABYTES": MB,

	"G":         GB,
	"GB":        GB,
	"GIGABYTE":  GB,
	"GIGABYTES": GB,

	"T":         TB,
	"TB":        TB,
	"TERABYTE":  TB,
	"TERABYTES": TB,
}

var (
	LongUnits bool   = false
	Format    string = "%.8f"
)

func Parse(s string) (ByteSize, error) {
	s = strings.TrimSpace(s)

	split := make([]string, 0)
	for i, r := range s {
		if !unicode.IsDigit(r) && r != '.' {
			split = append(split, strings.TrimSpace(string(s[:i])))
			split = append(split, strings.TrimSpace(string(s[i:])))
			break
		}
	}

	if len(split) != 2 {
		return 0, errors.New("unrecognized size suffix")
	}

	unit, ok := unitMap[strings.ToUpper(split[1])]
	if !ok {
		return 0, errors.New("unrecognized size suffix " + split[1])

	}

	value, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		return 0, err
	}

	bytesize := ByteSize(value * float64(unit))

	return bytesize, nil
}

func (b ByteSize) Format(format string, unit string, longUnits bool) string {
	return b.format(format, unit, longUnits)
}

func (b ByteSize) String() string {
	return b.format(Format, "", LongUnits)
}

func (b ByteSize) format(format string, unit string, longUnits bool) string {
	var unitSize ByteSize
	if unit != "" {
		var ok bool
		unitSize, ok = unitMap[strings.ToUpper(unit)]
		if !ok {
			return "Unrecognized unit: " + unit
		}
	} else {
		switch {
		case b >= TB:
			unitSize = TB
		case b >= GB:
			unitSize = GB
		case b >= MB:
			unitSize = MB
		case b >= KB:
			unitSize = KB
		default:
			unitSize = B
		}
	}

	if longUnits {
		var s string
		value := fmt.Sprintf(format, float64(b)/float64(unitSize))
		if printS, _ := strconv.ParseFloat(strings.TrimSpace(value), 64); printS > 0 && printS != 1 {
			s = "s"
		}

		return fmt.Sprintf(format+longUnitMap[unitSize]+s, float64(b)/float64(unitSize))
	}

	return fmt.Sprintf(format+shortUnitMap[unitSize], float64(b)/float64(unitSize))
}

func Float64(s string) (float64, error) {
	s = strings.TrimSpace(s)

	split := make([]string, 0)
	for i, r := range s {
		if !unicode.IsDigit(r) && r != '.' {
			split = append(split, strings.TrimSpace(string(s[:i])))
			split = append(split, strings.TrimSpace(string(s[i:])))
			break
		}
	}

	if len(split) != 2 {
		return 0, errors.New("unrecognized size suffix")
	}

	value, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
