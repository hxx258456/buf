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

package bufprotopluginexec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"

	"buf.build/go/app"
	"buf.build/go/standard/xio"
	"buf.build/go/standard/xlog/xslog"
	"buf.build/go/standard/xos/xexec"
	"github.com/bufbuild/buf/private/pkg/protoencoding"
	"github.com/bufbuild/buf/private/pkg/storage"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"github.com/bufbuild/buf/private/pkg/tmp"
	"github.com/bufbuild/protoplugin"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type protocProxyHandler struct {
	logger            *slog.Logger
	storageosProvider storageos.Provider
	protocPath        string
	protocExtraArgs   []string
	pluginName        string
}

func newProtocProxyHandler(
	logger *slog.Logger,
	storageosProvider storageos.Provider,
	protocPath string,
	protocExtraArgs []string,
	pluginName string,
) *protocProxyHandler {
	return &protocProxyHandler{
		logger:            logger,
		storageosProvider: storageosProvider,
		protocPath:        protocPath,
		protocExtraArgs:   protocExtraArgs,
		pluginName:        pluginName,
	}
}

func (h *protocProxyHandler) Handle(
	ctx context.Context,
	pluginEnv protoplugin.PluginEnv,
	responseWriter protoplugin.ResponseWriter,
	request protoplugin.Request,
) (retErr error) {
	defer xslog.DebugProfile(h.logger, slog.String("plugin", filepath.Base(h.pluginName)))()

	// We should send the complete FileDescriptorSet with source-retention options to --descriptor_set_in.
	//
	// This is used via the FileDescriptorSet below.
	request, err := request.WithSourceRetentionOptions()
	if err != nil {
		return err
	}

	protocVersion, err := h.getProtocVersion(ctx, pluginEnv)
	if err != nil {
		return err
	}
	if h.pluginName == "kotlin" && !getKotlinSupportedAsBuiltin(protocVersion) {
		return fmt.Errorf("kotlin is not supported for protoc version %s", versionString(protocVersion))
	}
	if h.pluginName == "rust" && !getRustSupportedAsBuiltin(protocVersion) {
		return fmt.Errorf("rust is not supported for protoc version %s", versionString(protocVersion))
	}
	// When we create protocProxyHandlers in NewHandler, we always prefer protoc-gen-.* plugins
	// over builtin plugins, so we only get here if we did not find protoc-gen-js, so this
	// is an error
	if h.pluginName == "js" && !getJSSupportedAsBuiltin(protocVersion) {
		return errors.New("js moved to a separate plugin hosted at https://github.com/protocolbuffers/protobuf-javascript in v21, you must install this plugin")
	}
	fileDescriptorSet := &descriptorpb.FileDescriptorSet{
		File: request.AllFileDescriptorProtos(),
	}
	fileDescriptorSetData, err := protoencoding.NewWireMarshaler().Marshal(fileDescriptorSet)
	if err != nil {
		return err
	}
	descriptorFilePath := app.DevStdinFilePath
	var tmpFile tmp.File
	if descriptorFilePath == "" {
		// since we have no stdin file (i.e. Windows), we're going to have to use a temporary file
		tmpFile, err = tmp.NewFile(ctx, bytes.NewReader(fileDescriptorSetData))
		if err != nil {
			return err
		}
		defer func() {
			retErr = errors.Join(retErr, tmpFile.Close())
		}()
		descriptorFilePath = tmpFile.Path()
	}
	tmpDir, err := tmp.NewDir(ctx)
	if err != nil {
		return err
	}
	defer func() {
		retErr = errors.Join(retErr, tmpDir.Close())
	}()
	args := slices.Concat(h.protocExtraArgs, []string{
		fmt.Sprintf("--descriptor_set_in=%s", descriptorFilePath),
		fmt.Sprintf("--%s_out=%s", h.pluginName, tmpDir.Path()),
	})
	if getSetExperimentalAllowProto3OptionalFlag(protocVersion) {
		args = append(
			args,
			"--experimental_allow_proto3_optional",
		)
	}
	if parameter := request.Parameter(); parameter != "" {
		args = append(
			args,
			fmt.Sprintf("--%s_opt=%s", h.pluginName, parameter),
		)
	}
	args = append(
		args,
		request.CodeGeneratorRequest().GetFileToGenerate()...,
	)
	stdin := xio.DiscardReader
	if descriptorFilePath != "" && descriptorFilePath == app.DevStdinFilePath {
		stdin = bytes.NewReader(fileDescriptorSetData)
	}
	if err := xexec.Run(
		ctx,
		h.protocPath,
		xexec.WithArgs(args...),
		xexec.WithEnv(pluginEnv.Environ),
		xexec.WithStdin(stdin),
		xexec.WithStderr(pluginEnv.Stderr),
	); err != nil {
		// TODO: strip binary path as well?
		// We don't know if this is a system error or plugin error, so we assume system error
		return handlePotentialTooManyFilesError(err)
	}
	if getFeatureProto3OptionalSupported(protocVersion) {
		responseWriter.SetFeatureProto3Optional()
	}
	// We always claim support for all Editions in the response because the invocation to
	// "protoc" will fail if it can't handle the input editions. That way, we don't have to
	// track which protoc versions support which editions and synthesize this information.
	// And that also lets us support users passing "--experimental_editions" to protoc.
	responseWriter.SetFeatureSupportsEditions(descriptorpb.Edition_EDITION_PROTO2, descriptorpb.Edition_EDITION_MAX)

	// no need for symlinks here, and don't want to support
	readWriteBucket, err := h.storageosProvider.NewReadWriteBucket(tmpDir.Path())
	if err != nil {
		return err
	}
	return storage.WalkReadObjects(
		ctx,
		readWriteBucket,
		"",
		func(readObject storage.ReadObject) error {
			data, err := io.ReadAll(readObject)
			if err != nil {
				return err
			}
			responseWriter.AddFile(readObject.Path(), string(data))
			return nil
		},
	)
}

func (h *protocProxyHandler) getProtocVersion(
	ctx context.Context,
	pluginEnv protoplugin.PluginEnv,
) (*pluginpb.Version, error) {
	stdoutBuffer := bytes.NewBuffer(nil)
	if err := xexec.Run(
		ctx,
		h.protocPath,
		xexec.WithArgs(slices.Concat(h.protocExtraArgs, []string{"--version"})...),
		xexec.WithEnv(pluginEnv.Environ),
		xexec.WithStdout(stdoutBuffer),
	); err != nil {
		// TODO: strip binary path as well?
		return nil, handlePotentialTooManyFilesError(err)
	}
	return parseVersionForCLIVersion(strings.TrimSpace(stdoutBuffer.String()))
}
