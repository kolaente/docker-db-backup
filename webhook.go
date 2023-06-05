package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func callWebhook() error {
	if config.CompletionWebhookURL == "" {
		return nil
	}

	res, err := http.Get(config.CompletionWebhookURL)
	if err != nil {
		return err
	}

	if res.StatusCode > 399 {
		buf := bytes.Buffer{}
		_, _ = buf.ReadFrom(res.Body)
		return fmt.Errorf("recived an error status code while calling the webhook: %d, message: %s", res.StatusCode, buf.String())
	}

	return nil
}
