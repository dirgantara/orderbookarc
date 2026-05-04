package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// hmacSign computes HMAC-SHA256 hex of payload using secret.
// Matches orderbook middleware: payload = method + requestURI + timestamp + nonce
func hmacSign(payload, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// PlaceOrderRequest mirrors orderbook's PlaceOrderRequest DTO.
type PlaceOrderRequest struct {
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`
	Type     string `json:"type"`
	Price    string `json:"price,omitempty"`
	Quantity string `json:"quantity"`
}

// callOrderbook sends a signed request to the orderbook server.
func callOrderbook(orderbookURL, apiKeyID, apiKey, apiSecret, method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, orderbookURL+path, reqBody)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := uuid.New().String()

	// payload = method + requestURI + timestamp + nonce  (matches middleware)
	payload := method + path + ts + nonce
	sig := hmacSign(payload, apiSecret)
	fmt.Println("fullpath::", orderbookURL+path)
	fmt.Println("payload::", payload)

	req.Header.Set("X-API-KEY-ID", apiKeyID)
	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("X-TIMESTAMP", ts)
	req.Header.Set("X-NONCE", nonce)
	req.Header.Set("X-SIGNATURE", sig)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("orderbook request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	return respBody, resp.StatusCode, err
}
