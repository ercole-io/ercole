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

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var cmd = &cobra.Command{
		Use:   "print-add-template",
		Short: "Print a Oracle/Database agreement",
		Long:  `Print a Oracle/Database agreement to stdout`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`{
	"agreementID": "5051863",
	"partsID": [
		"L10006",
		"A90620"
	],
	"csi": "6871235",
	"referenceNumber": "10032246681",
	"unlimited": false,
	"count": 30,
	"catchAll": false,
	"hosts": [
		"pluto",
		"pippo"
	]
}`)
		},
	}

	agreementOracleDatabaseCmd.AddCommand(cmd)
}
