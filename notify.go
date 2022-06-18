package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

type BridgeUpdateNotification struct {
	Environment BeeperEnv `json:"-"`

	Channel  BeeperChannel `json:"channel"`
	Bridge   BridgeType    `json:"bridge"`
	Image    string        `json:"image"`
	Password string        `json:"password"`

	DeployNext bool `json:"deployNext,omitempty"`
}

func (bun BridgeUpdateNotification) prepareRequest() (*http.Request, error) {
	url := env(fmt.Sprintf("BEEPER_%s_ADMIN_API_URL", bun.Environment))
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&bun); err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	} else if req, err := http.NewRequest(http.MethodPost, url, &body); err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	} else {
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}
}

func (bun BridgeUpdateNotification) Send() error {
	if req, err := bun.prepareRequest(); err != nil {
		return err
	} else if resp, err := client.Do(req); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	} else if respBody, err := io.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("failed to read response body (status: %d): %w", resp.StatusCode, err)
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("got non-200 status code %d: %s", resp.StatusCode, respBody)
	} else {
		log.Printf("Successfully notified Beeper %s/%s about update to %s (status: %d, body: %s)", bun.Environment, bun.Channel, bun.Bridge, resp.StatusCode, respBody)
		return nil
	}
}

func (bun *BridgeUpdateNotification) Fill(bridge BridgeType, image string) *BridgeUpdateNotification {
	if bun.Bridge == "" {
		bun.Bridge = bridge
	}
	bun.Image = image
	bun.Password = env(fmt.Sprintf("BEEPER_%s_ADMIN_NIGHTLY_PASS", bun.Environment))
	return bun
}
