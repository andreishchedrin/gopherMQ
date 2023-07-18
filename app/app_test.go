package app

import (
	"andreishchedrin/gopherMQ/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
)

func TestAppPushPull(t *testing.T) {
	cfg, err := config.NewConfig("../config/.env")
	if err != nil {
		panic(err)
	}

	app := NewApp(cfg)
	go func() {
		app.Start()
	}()

	t.Run("push-pull", func(t *testing.T) {
		app.HttpServer.App.Post("/push", app.HttpServer.PushHandler)

		jsonBody1 := map[string]interface{}{
			"channel": "channel1",
			"message": "payload",
		}
		body1, _ := json.Marshal(jsonBody1)

		req1 := httptest.NewRequest("POST", "/push", bytes.NewReader(body1))
		req1.Header.Set("Content-Type", "application/json")

		resp1, err := app.HttpServer.App.Test(req1)

		if err != nil {
			t.Errorf("error: %v", err)
		}

		if resp1.StatusCode != 200 {
			t.Errorf("got %v, want %v", resp1.StatusCode, 200)
		}

		app.HttpServer.App.Get("/pull", app.HttpServer.PullHandler)

		jsonBody2 := map[string]interface{}{
			"channel": "channel1",
		}
		body2, _ := json.Marshal(jsonBody2)

		req2 := httptest.NewRequest("GET", "/pull", bytes.NewReader(body2))
		req2.Header.Set("Content-Type", "application/json")

		resp2, err := app.HttpServer.App.Test(req2)

		if err != nil {
			t.Errorf("error: %v", err)
		}

		if resp2.StatusCode != 200 {
			t.Errorf("got %v, want %v", resp2.StatusCode, 200)
		}

		b2, err := io.ReadAll(resp2.Body)
		if err != nil {
			t.Error(err)
		}

		//TODO temporary assert package go ver conflict
		if string(b2) != "payload" {
			t.Error("test failed")
		}
	})
}
