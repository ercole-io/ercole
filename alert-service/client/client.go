package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

type AlertSvcClientInterface interface {
	ThrowNewAlert(alert model.Alert) error
}

type Client struct {
	remoteEndpoint string
	client         *http.Client
	config         config.AlertService
}

func NewClient(config config.AlertService) *Client {
	return &Client{
		remoteEndpoint: strings.TrimSuffix(config.RemoteEndpoint, "/"),
		client:         &http.Client{},
		config:         config,
	}
}

func (c *Client) doRequest(ctx context.Context, path, method string, body []byte) (*http.Response, error) {
	url := utils.NewAPIUrlNoParams(
		c.remoteEndpoint,
		c.config.PublisherUsername,
		c.config.PublisherPassword,
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
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, err
		}

		return resp, fmt.Errorf("Api error (code: %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) ThrowNewAlert(alert model.Alert) error {
	body, err := json.Marshal(alert)
	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "Can't marshal alert")
	}

	_, err = c.getResponse(context.TODO(), "/alerts", "POST", body)
	if err != nil {
		return err
	}

	return nil
}
