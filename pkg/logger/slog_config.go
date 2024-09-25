/*
 * SPDX-License-Identifier: Apache-2.0
 *  SPDX-FileCopyrightText: Huawei Inc.
 */

package logger

import (
	"log/slog"
	"os"
	"path/filepath"
)

var Logger *slog.Logger

// use a simple inbuilt slog logger to log to stdout.
func init() {
	var textHandler slog.Handler = slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.SourceKey {
					source, _ := a.Value.Any().(*slog.Source)
					if source != nil {
						// by default the full file path is set. This must be replaced with just file name.
						source.File = filepath.Base(source.File)
					}
				}
				return a
			},
		})
	Logger = slog.New(AgentLogHandler{textHandler})

}

// AgentLogHandler necessary to extend the logger to add more context information. Not used for now.
type AgentLogHandler struct {
	// embedded struct field to take any existing methods of the field, and we override only what we need.
	slog.Handler
}
