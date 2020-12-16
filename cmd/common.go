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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/jinzhu/now"
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
var alertStatus string
var alertStatusAll bool
var from string
var to string
var olderThan string
var outputFormat string
var mode string
var mode2 string
var metric string
var searchOracleDatabaseAgreementAgreementID string
var searchOracleDatabaseAgreementMetrics string
var searchOracleDatabaseAgreementPartID string
var searchOracleDatabaseAgreementItemDescription string
var searchOracleDatabaseAgreementCSI string
var searchOracleDatabaseAgreementReferenceNumber string
var searchOracleDatabaseAgreementUnlimited string
var searchOracleDatabaseAgreementCatchAll string
var searchOracleDatabaseAgreementLicensesCountLTE int
var searchOracleDatabaseAgreementLicensesCountGTE int
var searchOracleDatabaseAgreementUsersCountLTE int
var searchOracleDatabaseAgreementUsersCountGTE int
var searchOracleDatabaseAgreementAvailableCountLTE int
var searchOracleDatabaseAgreementAvailableCountGTE int

type apiOption struct {
	addOption func(cmd *cobra.Command)
	addParam  func(params url.Values)
}

var fullOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().BoolVarP(&summary, "summary", "s", false, "Summary mode")
	},
	addParam: func(params url.Values) {
		params.Set("full", strconv.FormatBool(!summary))
	},
}

var modeOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&mode, "mode", "m", "full", "Output mode (full, summary, lms)")
	},
	addParam: func(params url.Values) {
		params.Set("mode", mode)
	},
}

var modeAlertOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&mode2, "mode", "m", "all", "Output mode (all, aggregated-code-severity, aggregated-category-severity)")
	},
	addParam: func(params url.Values) {
		params.Set("mode", mode2)
	},
}

var windowTimeOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().IntVarP(&windowTime, "window-time", "w", 12, "Window time")
	},
	addParam: func(params url.Values) {
		params.Set("window-time", strconv.Itoa(windowTime))
	},
}

var locationOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&location, "location", "l", "", "Filter by location")
	},
	addParam: func(params url.Values) {
		params.Set("location", location)
	},
}

var environmentOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&environment, "environment", "e", "", "Filter by environment")
	},
	addParam: func(params url.Values) {
		params.Set("environment", environment)
	},
}

var sortingOptions apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort by field")
		cmd.Flags().BoolVar(&sortDesc, "desc-order", false, "Sort descending")
	},
	addParam: func(params url.Values) {
		if sortBy != "" {
			params.Set("sort-by", sortBy)
			params.Set("sort-desc", strconv.FormatBool(sortDesc))
		}
	},
}

var limitOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().IntVarP(&limit, "limit", "n", 15, "Limit the number of databases")
	},
	addParam: func(params url.Values) {
		params.Set("limit", strconv.Itoa(limit))
	},
}

var severityOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&severity, "severity", "", "Filter by severity")
	},
	addParam: func(params url.Values) {
		params.Set("severity", severity)
	},
}

var alertStatusOptions apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&alertStatus, "status", model.AlertStatusNew, "Filter by alert status")
		cmd.Flags().BoolVar(&alertStatusAll, "all", false, "Also show read alerts")
	},
	addParam: func(params url.Values) {
		if alertStatusAll {
			params.Set("status", "")
		} else {
			params.Set("status", alertStatus)
		}
	},
}

var fromOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&from, "from", "f", "", "From")
	},
	addParam: func(params url.Values) {
		if from != "" {
			if val, err := time.Parse(time.RFC3339, from); err == nil {
				from = val.Format(time.RFC3339)
			} else if val, err := time.Parse("2006-01-02", from); err == nil {
				from = val.Format(time.RFC3339)
			} else {
				fmt.Fprintf(os.Stderr, "Unable to parse the value of the from option: %v\n", err)
				os.Exit(1)
			}
		}
		params.Set("from", from)
	},
}

var fromToWindowOptions apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&from, "from", "F", "", "Filter alerts with a date >= from")
		cmd.Flags().StringVarP(&to, "to", "t", "", "Filter alerts with a date <= to")
	},
	addParam: func(params url.Values) {
		if from != "" {
			if val, err := time.Parse(time.RFC3339, from); err == nil {
				from = val.Format(time.RFC3339)
			} else if val, err := time.Parse("2006-01-02", from); err == nil {
				from = val.Format(time.RFC3339)
			} else {
				fmt.Fprintf(os.Stderr, "Unable to parse the value of the from option: %v\n", err)
				os.Exit(1)
			}
		}
		if to != "" {
			if val, err := time.Parse(time.RFC3339, to); err == nil {
				to = val.Format(time.RFC3339)
			} else if val, err := time.Parse("2006-01-02", to); err == nil {
				to = now.New(val).EndOfDay().Format(time.RFC3339)
			} else {
				fmt.Fprintf(os.Stderr, "Unable to parse the value of the to option: %v\n", err)
				os.Exit(1)
			}
		}
		params.Set("from", from)
		params.Set("to", to)
	},
}

var olderThanOptions apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&olderThan, "older-than", "t", "", "Filter data older than day")
	},
	addParam: func(params url.Values) {
		if olderThan != "" {
			if val, err := time.Parse(time.RFC3339, olderThan); err == nil {
				olderThan = val.Format(time.RFC3339)
			} else if val, err := time.Parse("2006-01-02", olderThan); err == nil {
				olderThan = now.New(val).EndOfDay().Format(time.RFC3339)
			} else {
				fmt.Fprintf(os.Stderr, "Unable to parse the value of the older than option: %v\n", err)
				os.Exit(1)
			}
		}
		params.Set("older-than", olderThan)
	},
}

var metricOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&metric, "metric", "m", "all", "metric (work, version)")
	},
	addParam: func(params url.Values) {
		params.Set("metric", metric)
	},
}

var searchOracleDatabaseAgreementsOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&searchOracleDatabaseAgreementAgreementID, "agreement-id", "", "")
		cmd.Flags().StringVar(&searchOracleDatabaseAgreementPartID, "part-id", "", "")
		cmd.Flags().StringVar(&searchOracleDatabaseAgreementItemDescription, "item-description", "", "")
		cmd.Flags().StringVar(&searchOracleDatabaseAgreementCSI, "csi", "", "")
		cmd.Flags().StringVar(&searchOracleDatabaseAgreementMetrics, "metric", "", "")
		cmd.Flags().StringVar(&searchOracleDatabaseAgreementReferenceNumber, "reference-number", "", "")
		cmd.Flags().StringVar(&searchOracleDatabaseAgreementUnlimited, "unlimited", "", "NULL,true,false")
		cmd.Flags().StringVar(&searchOracleDatabaseAgreementCatchAll, "catch-all", "", "NULL,true,false")
		cmd.Flags().IntVar(&searchOracleDatabaseAgreementLicensesCountLTE, "licenses-count-lte", -1, "")
		cmd.Flags().IntVar(&searchOracleDatabaseAgreementLicensesCountGTE, "licenses-count-gte", -1, "")
		cmd.Flags().IntVar(&searchOracleDatabaseAgreementUsersCountLTE, "users-count-lte", -1, "")
		cmd.Flags().IntVar(&searchOracleDatabaseAgreementUsersCountGTE, "users-count-gte", -1, "")
		cmd.Flags().IntVar(&searchOracleDatabaseAgreementAvailableCountLTE, "available-count-lte", -1, "")
		cmd.Flags().IntVar(&searchOracleDatabaseAgreementAvailableCountGTE, "available-count-gte", -1, "")
	},
	addParam: func(params url.Values) {
		params.Set("agreement-id", searchOracleDatabaseAgreementAgreementID)
		params.Set("part-id", searchOracleDatabaseAgreementPartID)
		params.Set("item-description", searchOracleDatabaseAgreementItemDescription)
		params.Set("csi", searchOracleDatabaseAgreementCSI)
		params.Set("metrics", searchOracleDatabaseAgreementMetrics)
		params.Set("unlimited", searchOracleDatabaseAgreementUnlimited)
		params.Set("catch-all", searchOracleDatabaseAgreementCatchAll)
		params.Set("licenses-count-lte", fmt.Sprintf("%d", searchOracleDatabaseAgreementLicensesCountLTE))
		params.Set("licenses-count-gte", fmt.Sprintf("%d", searchOracleDatabaseAgreementLicensesCountGTE))
		params.Set("users-count-lte", fmt.Sprintf("%d", searchOracleDatabaseAgreementUsersCountLTE))
		params.Set("users-count-gte", fmt.Sprintf("%d", searchOracleDatabaseAgreementUsersCountGTE))
		params.Set("available-count-lte", fmt.Sprintf("%d", searchOracleDatabaseAgreementAvailableCountLTE))
		params.Set("available-count-gte", fmt.Sprintf("%d", searchOracleDatabaseAgreementAvailableCountGTE))
	},
}

func simpleAPIRequestCommand(
	use string,
	short string,
	long string,
	searchArguments bool,
	anotherOptions []apiOption,
	customResponseTypes bool,
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
			for _, opt := range anotherOptions {
				opt.addParam(params)
			}

			//Make the http request
			req, _ := http.NewRequest("GET", utils.NewAPIUrl(
				ercoleConfig.APIService.RemoteEndpoint,
				ercoleConfig.APIService.AuthenticationProvider.Username,
				ercoleConfig.APIService.AuthenticationProvider.Password,
				endpointPath,
				params,
			).String(), bytes.NewReader([]byte{}))

			if customResponseTypes {
				switch outputFormat {
				case "json":
					outputFormat = "application/json"
				case "xlsx":
					outputFormat = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
				case "lms":
					outputFormat = "application/vnd.oracle.lms+vnd.ms-excel.sheet.macroEnabled.12"
				case "mongohostdata", "mhd":
					outputFormat = "application/vnd.ercole.mongohostdata+json"
				}
				req.Header.Set("Accept", outputFormat)
			}

			resp, err := http.DefaultClient.Do(req)
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

				if resp.Header.Get("Content-Type") == "application/json" {
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
				} else {
					os.Stdout.Write(out)
				}
			}

		},
	}

	if !searchArguments {
		cmd.Args = cobra.ExactArgs(0)
	}
	for _, opt := range anotherOptions {
		opt.addOption(cmd)
	}
	if customResponseTypes {
		cmd.Flags().StringVarP(&outputFormat, "format", "f", "application/json", "Output format")
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
	olderThanOption bool,
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
			if olderThanOption {
				olderThanOptions.addParam(params)
			}

			//Make the http request
			resp, err := http.Get(
				utils.NewAPIUrl(
					ercoleConfig.APIService.RemoteEndpoint,
					ercoleConfig.APIService.AuthenticationProvider.Username,
					ercoleConfig.APIService.AuthenticationProvider.Password,
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
	if olderThanOption {
		olderThanOptions.addOption(cmd)
	}
	return cmd
}

func simpleChartRequestCommandSingleValue(
	use string,
	short string,
	long string,
	searchArguments bool,
	anotherOptions []apiOption,
	customResponseTypes bool,
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
			for _, opt := range anotherOptions {
				opt.addParam(params)
			}

			//Make the http request
			req, _ := http.NewRequest("GET", utils.NewAPIUrl(
				ercoleConfig.ChartService.RemoteEndpoint,
				ercoleConfig.APIService.AuthenticationProvider.Username,
				ercoleConfig.APIService.AuthenticationProvider.Password,
				endpointPath,
				params,
			).String(), bytes.NewReader([]byte{}))

			if customResponseTypes {
				switch outputFormat {
				case "json":
					outputFormat = "application/json"
				}
				req.Header.Set("Accept", outputFormat)
			}

			resp, err := http.DefaultClient.Do(req)
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

				if resp.Header.Get("Content-Type") == "application/json" {
					var res interface{}
					err = json.Unmarshal(out, &res)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Failed to unmarshal response body: %v (%s)\n", err, string(out))
						os.Exit(1)
					}

					enc := json.NewEncoder(os.Stdout)
					enc.SetIndent("", "    ")
					enc.Encode(res)
				} else {
					os.Stdout.Write(out)
				}
			}

		},
	}

	if !searchArguments {
		cmd.Args = cobra.ExactArgs(0)
	}
	for _, opt := range anotherOptions {
		opt.addOption(cmd)
	}
	if customResponseTypes {
		cmd.Flags().StringVarP(&outputFormat, "format", "f", "application/json", "Output format")
	}
	return cmd
}
