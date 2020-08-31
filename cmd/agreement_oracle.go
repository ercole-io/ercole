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
	"github.com/spf13/cobra"
)

// agreementOracleCmd represents the oracle command
var agreementOracleCmd = &cobra.Command{
	Use:   "oracle",
	Short: "Perform an api request related to oracle agreements",
	Long:  `oracle perform an api request related to oracle agreements`,
}

func init() {
	agreementCmd.AddCommand(agreementOracleCmd)
}
