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

func init() {
	textHandler := slog.NewTextHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.SourceKey {
					source, _ := a.Value.Any().(*slog.Source)
					if source != nil {
						source.File = filepath.Base(source.File)
					}
				}
				return a
			},
		})
	Logger = slog.New(textHandler)

}
