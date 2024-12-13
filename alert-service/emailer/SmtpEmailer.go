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

package emailer

import (
	"bytes"
	"crypto/tls"
	"io"

	gomail "gopkg.in/gomail.v2"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
)

type SMTPEmailer struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
}

func (emailer *SMTPEmailer) SendEmail(subject string, text string, to []string) error {
	if !emailer.Config.AlertService.Emailer.Enabled {
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailer.Config.AlertService.Emailer.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(emailer.Config.AlertService.Emailer.SMTPServer,
		emailer.Config.AlertService.Emailer.SMTPPort,
		emailer.Config.AlertService.Emailer.SMTPUsername,
		emailer.Config.AlertService.Emailer.SMTPPassword)

	if emailer.Config.AlertService.Emailer.DisableSSLCertificateValidation {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	err := d.DialAndSend(m)
	if err != nil {
		return utils.NewError(err, "EMAILER")
	} else {
		return nil
	}
}

func (emailer *SMTPEmailer) send(m gomail.Message) error {
	d := gomail.NewDialer(emailer.Config.AlertService.Emailer.SMTPServer,
		emailer.Config.AlertService.Emailer.SMTPPort,
		emailer.Config.AlertService.Emailer.SMTPUsername,
		emailer.Config.AlertService.Emailer.SMTPPassword)

	if emailer.Config.AlertService.Emailer.DisableSSLCertificateValidation {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	err := d.DialAndSend(&m)
	if err != nil {
		return utils.NewError(err, "EMAILER")
	} else {
		return nil
	}
}

func (emailer *SMTPEmailer) SendHtmlEmail(subject, text string, to []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", emailer.Config.AlertService.Emailer.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	return emailer.send(*m)
}

func (emailer *SMTPEmailer) SendReportEmail(subject string, to []string, attachmentBuff bytes.Buffer) error {
	m := gomail.NewMessage()
	m.SetHeader("From", emailer.Config.AlertService.Emailer.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", "Please see attached file.")
	m.Attach("alert_report.xlsx", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(attachmentBuff.Bytes())
		return err
	}))

	return emailer.send(*m)
}
