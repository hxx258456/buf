// Copyright 2020-2025 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris

package buffetch

import (
	"testing"

	"buf.build/go/app"
	"github.com/bufbuild/buf/private/buf/buffetch/internal"
)

func TestGetParsedRefSuccess_UnixOnly(t *testing.T) {
	t.Parallel()
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBinpb,
			"",
			internal.FileSchemeStdin,
			internal.CompressionTypeNone,
			nil,
		),
		app.DevStdinFilePath,
	)
	testGetParsedRefSuccess(
		t,
		internal.NewDirectParsedSingleRef(
			formatBinpb,
			"",
			internal.FileSchemeStdout,
			internal.CompressionTypeNone,
			nil,
		),
		app.DevStdoutFilePath,
	)
}
