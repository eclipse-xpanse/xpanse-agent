/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package logger

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type HttpRequestLogger struct {
	http.RoundTripper
}

func (httpRequestLogger *HttpRequestLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request URL
	Logger.Info(fmt.Sprintf("Request URL: %s", req.URL.String()))

	// Read and log the request body if it exists
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		// Log the request body
		Logger.Info(fmt.Sprintf("Request Body: %s", string(bodyBytes)))

		// Restore the io.ReadCloser to its original state
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Make the actual HTTP request using the embedded RoundTripper
	resp, err := httpRequestLogger.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Log the response status code
	Logger.Info(fmt.Sprintf("Response Status: %d", resp.StatusCode))

	// Read and log the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Log the response body
	Logger.Info("Response Body: ", slog.String("msg", string(responseBody)))

	// Restore the response body so it can be read later by application code
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	return resp, nil
}
