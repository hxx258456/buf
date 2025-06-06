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

package bufimage_test

import (
	"testing"

	"buf.build/go/standard/xtesting"
	"github.com/bufbuild/buf/private/buf/buftesting"
	"github.com/bufbuild/buf/private/bufpkg/bufimage"
	"github.com/bufbuild/buf/private/pkg/prototesting"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestCompareGoogleapis(t *testing.T) {
	xtesting.SkipIfShort(t)
	// Don't run in parallel as it allocates a lot of memory
	// cannot directly compare with source code info as buf alpha protoc creates additional source
	// code infos that protoc does not
	image := testBuildGoogleapis(t, false)
	fileDescriptorSet := bufimage.ImageToFileDescriptorSet(image)
	actualProtocFileDescriptorSet := testBuildActualProtocGoogleapis(t, false)
	prototesting.AssertFileDescriptorSetsEqual(
		t,
		fileDescriptorSet,
		actualProtocFileDescriptorSet,
	)
}

func testBuildActualProtocGoogleapis(t *testing.T, includeSourceInfo bool) *descriptorpb.FileDescriptorSet {
	googleapisDirPath := buftesting.GetGoogleapisDirPath(t, buftestingDirPath)
	filePaths := buftesting.GetProtocFilePaths(t, googleapisDirPath, 0)
	fileDescriptorSet := buftesting.GetActualProtocFileDescriptorSet(t, true, includeSourceInfo, googleapisDirPath, filePaths)
	assert.Equal(t, buftesting.NumGoogleapisFilesWithImports, len(fileDescriptorSet.GetFile()))

	return fileDescriptorSet
}
