// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/amreo/ercole-services/utils"
	"github.com/spf13/cobra"
)

//common options
var verbose bool
var summary bool
var sortBy string
var sortDesc bool
var location string
var environment string
var windowTime int
var limit int
var severity string

func simpleAPIRequestCommand(
	use string,
	short string,
	long string,
	searchArguments bool,
	fullOption bool,
	windowTimeOption bool,
	locationOption bool,
	environmentOption bool,
	sortableResult bool,
	limitOption bool,
	severityOption bool,
	endpointPath string,
	errorMessageFormat string,
	httpErrorMessageFormat string) *cobra.Command {

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			//Set query params
			params := url.Values{}
			if searchArguments {
				params.Set("search", strings.Join(args, " "))
			}
			if fullOption {
				params.Set("full", strconv.FormatBool(!summary))
			}
			if windowTimeOption {
				params.Set("window-time", strconv.Itoa(windowTime))
			}
			if locationOption {
				params.Set("location", location)
			}
			if environmentOption {
				params.Set("environment", environment)
			}
			if limitOption {
				params.Set("limit", strconv.Itoa(limit))
			}
			if sortableResult && sortBy != "" {
				params.Set("sort-by", sortBy)
				params.Set("sort-desc", strconv.FormatBool(sortDesc))
			}
			if severityOption {
				params.Set("severity", severity)
			}
			//Make the http request
			resp, err := http.Get(
				utils.NewAPIUrl(
					ercoleConfig.APIService.RemoteEndpoint,
					ercoleConfig.APIService.UserUsername,
					ercoleConfig.APIService.UserPassword,
					endpointPath,
					params,
				).String())
			if err != nil {
				fmt.Fprintf(os.Stderr, errorMessageFormat, err)
				os.Exit(1)
			} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
				out, _ := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				fmt.Fprintf(os.Stderr, httpErrorMessageFormat, resp.StatusCode, string(out))
				os.Exit(1)
			} else {
				out, _ := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				var res []interface{}
				err = json.Unmarshal(out, &res)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to unmarshal response body: %v (%s)\n", err, string(out))
					os.Exit(1)
				}

				for _, item := range res {
					enc := json.NewEncoder(os.Stdout)
					enc.SetIndent("", "    ")
					enc.Encode(item)
				}
			}

		},
	}

	if !searchArguments {
		cmd.Args = cobra.ExactArgs(0)
	}
	if fullOption {
		cmd.Flags().BoolVarP(&summary, "summary", "s", false, "Summary mode")
	}
	if windowTimeOption {
		cmd.Flags().IntVarP(&windowTime, "window-time", "w", 12, "Window time")
	}
	if locationOption {
		cmd.Flags().StringVarP(&location, "location", "l", "", "Filter by location")
	}
	if environmentOption {
		cmd.Flags().StringVarP(&environment, "environment", "e", "", "Filter by environment")
	}
	if sortableResult {
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort by field")
		cmd.Flags().BoolVar(&sortDesc, "desc-order", false, "Sort descending")
	}
	if limitOption {
		cmd.Flags().IntVarP(&limit, "limit", "n", 15, "Limit the number of databases")
	}
	if severityOption {
		cmd.Flags().StringVar(&severity, "severity", "", "Limit the number of databases")
	}
	return cmd
}

func simpleSingleValueAPIRequestCommand(
	use string,
	short string,
	long string,
	searchArguments bool,
	locationOption bool,
	environmentOption bool,
	endpointPath string,
	errorMessageFormat string,
	httpErrorMessageFormat string) *cobra.Command {

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			//Set query params
			params := url.Values{}
			if searchArguments {
				params.Set("search", strings.Join(args, " "))
			}
			if locationOption {
				params.Set("location", location)
			}
			if environmentOption {
				params.Set("environment", environment)
			}

			//Make the http request
			resp, err := http.Get(
				utils.NewAPIUrl(
					ercoleConfig.APIService.RemoteEndpoint,
					ercoleConfig.APIService.UserUsername,
					ercoleConfig.APIService.UserPassword,
					endpointPath,
					params,
				).String())
			if err != nil {
				fmt.Fprintf(os.Stderr, errorMessageFormat, err)
				os.Exit(1)
			} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
				out, _ := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				fmt.Fprintf(os.Stderr, httpErrorMessageFormat, resp.StatusCode, string(out))
				os.Exit(1)
			} else {
				out, _ := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()
				var res interface{}
				err = json.Unmarshal(out, &res)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to unmarshal response body: %v (%s)\n", err, string(out))
					os.Exit(1)
				}

				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "    ")
				enc.Encode(res)
			}

		},
	}

	if !searchArguments {
		cmd.Args = cobra.ExactArgs(0)
	}
	if locationOption {
		cmd.Flags().StringVarP(&location, "location", "l", "", "Filter by location")
	}
	if environmentOption {
		cmd.Flags().StringVarP(&environment, "environment", "e", "", "Filter by environment")
	}
	return cmd
}
