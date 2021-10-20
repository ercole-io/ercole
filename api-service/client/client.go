// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

type ApiSvcClientInterface interface {
	GetAlertsByFilter(filter dto.AlertsFilter) ([]model.Alert, error)
	AckAlerts(filter dto.AlertsFilter) error
}

type Client struct {
	remoteEndpoint string
	client         *http.Client
	config         config.APIService
}

func NewClient(config config.APIService) *Client {
	return &Client{
		remoteEndpoint: strings.TrimSuffix(config.RemoteEndpoint, "/"),
		client:         &http.Client{},
		config:         config,
	}
}

func (c *Client) doRequest(ctx context.Context, path, method string, body []byte) (*http.Response, error) {
	url := utils.NewAPIUrlNoParams(
		c.remoteEndpoint,
		c.config.AuthenticationProvider.Username,
		c.config.AuthenticationProvider.Password,
		path)

	req, err := http.NewRequest(method, url.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) getResponse(ctx context.Context, path, method string, body []byte) (*http.Response, error) {
	resp, err := c.doRequest(ctx, path, method, body)
	if err != nil {
		return nil, err
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		return resp, fmt.Errorf("Api error (code: %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) getParsedResponse(ctx context.Context, path, method string, body []byte, response interface{}) (*http.Response, error) {
	resp, err := c.getResponse(ctx, path, method, body)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	return resp, d.Decode(response)
}
