/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"xpanse-agent/pkg/logger"
)

type HttpRequestLogger struct {
	http.RoundTripper
}

func (httpRequestLogger *HttpRequestLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	// Log the request URL
	logger.Logger.Info(fmt.Sprintf("Request URL: %s", req.URL.String()))

	// Read and log the request body if it exists
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		// Log the request body
		logger.Logger.Info(fmt.Sprintf("Request Body: %s", string(bodyBytes)))

		// Restore the io.ReadCloser to its original state
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Make the actual HTTP request using the embedded RoundTripper
	resp, err := httpRequestLogger.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Log the response status code
	logger.Logger.Info(fmt.Sprintf("Response Status: %d", resp.StatusCode))

	// Read and log the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Log the response body
	logger.Logger.Info("Response Body: ", slog.String("msg", string(responseBody)))

	// Restore the response body so it can be read later by application code
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(responseBody))

	return resp, nil
}
