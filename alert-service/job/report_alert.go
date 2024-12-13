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
package job

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/alert-service/database"
	"github.com/ercole-io/ercole/v2/alert-service/emailer"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

type ReportAlertJob struct {
	Database database.MongoDatabaseInterface
	Config   config.Configuration
	Log      logger.Logger
	Emailer  emailer.Emailer
}

func (r *ReportAlertJob) Run() {
	cronIntervals := map[string]int{
		"@daily":   1,
		"@weekly":  7,
		"@monthly": 30,
	}

	crontab := strings.ToLower(r.Config.AlertService.ReportAlertJob.Crontab)

	days, validCron := cronIntervals[crontab]
	if !validCron {
		r.Log.Error("report alert job - invalid crontab configuration")
		return
	}

	now := time.Now()
	end := now.AddDate(0, 0, -days)

	alerts, err := r.Database.FindAlertsByDate(end, now)
	if err != nil {
		r.Log.Error(err)
		return
	}

	alertMails := r.getAlertMails(alerts)

	if len(alertMails.Alerts) == 0 {
		r.Log.Infof("report alert job - no new alerts found")
		return
	}

	subject := fmt.Sprintf("Ercole alert messages - %s", now.Format("02/01/2006"))

	message := ""

	switch {
	case len(alerts) <= 20:
		message = `<style>
table {
  font-family: arial, sans-serif;
  border-collapse: collapse;
  width: 100%;
}
td, th {
  border: 1px solid #dddddd;
  text-align: left;
  padding: 8px;
}
tr:nth-child(even) {
  background-color: #dddddd;
}
</style><table><tr><th>Type</th><th>Date</th><th>Severity</th><th>Hostnames</th><th>Code</th><th>Description</th></tr>`

		for _, alert := range alertMails.Alerts {
			hostname := ""
			if val, ok := alert.OtherInfo["hostname"]; ok {
				hostname = val.(string)
			}

			description := alert.Description
			if len(alert.Description) > 100 {
				description = fmt.Sprintf("%s ...", alert.Description[:96])
			}

			row := fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,
				alert.AlertCategory,
				alert.Date.Format("02/01/2006 15:04"),
				alert.AlertSeverity,
				hostname,
				alert.AlertCode,
				description,
			)

			message += row
		}

		message += "</table>"

		if err := r.Emailer.SendHtmlEmail(subject, message, alertMails.To); err != nil {
			r.Log.Error(err)
			return
		}

	default:
		file, err := r.createAlertReportXlsx(alertMails.Alerts)
		if err != nil {
			r.Log.Error(err)
			return
		}

		var buf bytes.Buffer
		if err := file.Write(&buf); err != nil {
			r.Log.Error(err)
			return
		}

		if err := r.Emailer.SendReportEmail(subject, alertMails.To, buf); err != nil {
			r.Log.Error(err)
			return
		}
	}
}

func (r *ReportAlertJob) createAlertReportXlsx(alerts []model.Alert) (*excelize.File, error) {
	sheet := "Alerts"
	headers := []string{
		"Type",
		"Date",
		"Severity",
		"Hostnames",
		"Code",
		"Description",
	}

	file, err := exutils.NewXLSX(r.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)
	for _, val := range alerts {
		nextAxis := axisHelp.NewRow()
		file.SetCellValue(sheet, nextAxis(), val.AlertCategory)
		file.SetCellValue(sheet, nextAxis(), val.Date.Format("02/01/2006 15:04"))
		file.SetCellValue(sheet, nextAxis(), val.AlertSeverity)

		hostname := ""
		if v, ok := val.OtherInfo["hostname"]; ok {
			hostname = v.(string)
		}

		file.SetCellValue(sheet, nextAxis(), hostname)
		file.SetCellValue(sheet, nextAxis(), val.AlertCode)
		file.SetCellValue(sheet, nextAxis(), val.Description)
	}

	return file, nil
}

type AlertMails struct {
	Alerts []model.Alert
	To     []string
}

func (r *ReportAlertJob) getAlertMails(alerts []model.Alert) AlertMails {
	to := r.Config.AlertService.Emailer.To
	alertMails := make([]model.Alert, 0, len(alerts))

	alertConfigMap := map[string]struct {
		Enable bool
		To     []string
	}{
		model.AlertCodeNewServer:               {r.Config.AlertService.Emailer.AlertType.NewHost.Enable, r.Config.AlertService.Emailer.AlertType.NewHost.To},
		model.AlertCodeNewDatabase:             {r.Config.AlertService.Emailer.AlertType.NewDatabase.Enable, r.Config.AlertService.Emailer.AlertType.NewDatabase.To},
		model.AlertCodeNewLicense:              {r.Config.AlertService.Emailer.AlertType.NewLicense.Enable, r.Config.AlertService.Emailer.AlertType.NewLicense.To},
		model.AlertCodeNewOption:               {r.Config.AlertService.Emailer.AlertType.NewOption.Enable, r.Config.AlertService.Emailer.AlertType.NewOption.To},
		model.AlertCodeUnlistedRunningDatabase: {r.Config.AlertService.Emailer.AlertType.NewUnlistedRunningDatabase.Enable, r.Config.AlertService.Emailer.AlertType.NewUnlistedRunningDatabase.To},
		model.AlertCodeIncreasedCPUCores:       {r.Config.AlertService.Emailer.AlertType.NewHostCpu.Enable, r.Config.AlertService.Emailer.AlertType.NewHostCpu.To},
		model.AlertCodeMissingPrimaryDatabase:  {r.Config.AlertService.Emailer.AlertType.MissingPrimaryDatabase.Enable, r.Config.AlertService.Emailer.AlertType.MissingPrimaryDatabase.To},
		model.AlertCodeMissingDatabase:         {r.Config.AlertService.Emailer.AlertType.MissingDatabase.Enable, r.Config.AlertService.Emailer.AlertType.MissingDatabase.To},
		model.AlertCodeAgentError:              {r.Config.AlertService.Emailer.AlertType.AgentError.Enable, r.Config.AlertService.Emailer.AlertType.AgentError.To},
		model.AlertCodeNoData:                  {r.Config.AlertService.Emailer.AlertType.NoData.Enable, r.Config.AlertService.Emailer.AlertType.NoData.To},
	}

	for _, alert := range alerts {
		if alert.AlertSeverity == model.AlertSeverityWarning && !r.Config.AlertService.Emailer.AlertSeverity.Warning {
			continue
		}

		if config, exists := alertConfigMap[alert.AlertCode]; exists && config.Enable {
			to = append(to, config.To...)
			alertMails = append(alertMails, alert)
		}
	}

	to = utils.RemoveDuplicate(to)

	return AlertMails{
		Alerts: alertMails,
		To:     to,
	}
}
