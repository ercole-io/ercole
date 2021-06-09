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
	"github.com/ercole-io/ercole/v2/utils"
)

type ApiSvcClientInterface interface {
	AckAlertsByFilter(filter dto.AlertsFilter) error
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

func (c *Client) AckAlertsByFilter(filter dto.AlertsFilter) error {
	b := struct {
		Filter dto.AlertsFilter `json:"filter"`
	}{
		Filter: filter,
	}
	body, err := json.Marshal(b)
	if err != nil {
		return utils.NewError(err, "Can't marshal")
	}

	_, err = c.getResponse(context.TODO(), "/alerts/ack", "POST", body)
	if err != nil {
		return err
	}

	return nil
}
