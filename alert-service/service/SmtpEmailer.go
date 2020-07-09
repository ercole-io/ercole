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

package service

import (
	"crypto/tls"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	gomail "gopkg.in/gomail.v2"
)

type SMTPEmailer struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
}

func (this *SMTPEmailer) SendEmail(subject string, text string, to []string) utils.AdvancedErrorInterface {
	if !this.Config.AlertService.Emailer.Enabled {
		return nil
	}
	m := gomail.NewMessage()
	m.SetHeader("From", this.Config.AlertService.Emailer.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(this.Config.AlertService.Emailer.SMTPServer,
		this.Config.AlertService.Emailer.SMTPPort,
		this.Config.AlertService.Emailer.SMTPUsername,
		this.Config.AlertService.Emailer.SMTPPassword)
	if this.Config.AlertService.Emailer.DisableSSLCertificateValidation {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	err := d.DialAndSend(m)
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "EMAILER")
	} else {
		return nil
	}
}
