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
	"context"
	"log/slog"

	"buf.build/go/app"
	"github.com/bufbuild/buf/private/bufpkg/bufprotoplugin"
	"github.com/bufbuild/buf/private/pkg/storage/storageos"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct {
	logger            *slog.Logger
	storageosProvider storageos.Provider
}

func newGenerator(
	logger *slog.Logger,
	storageosProvider storageos.Provider,
) *generator {
	return &generator{
		logger:            logger,
		storageosProvider: storageosProvider,
	}
}

func (g *generator) Generate(
	ctx context.Context,
	container app.EnvStderrContainer,
	pluginName string,
	requests []*pluginpb.CodeGeneratorRequest,
	options ...GenerateOption,
) (_ *pluginpb.CodeGeneratorResponse, retErr error) {
	generateOptions := newGenerateOptions()
	for _, option := range options {
		option(generateOptions)
	}
	handlerOptions := []HandlerOption{
		HandlerWithPluginPath(generateOptions.pluginPath...),
		HandlerWithProtocPath(generateOptions.protocPath...),
	}
	handler, err := NewHandler(
		g.logger,
		g.storageosProvider,
		pluginName,
		handlerOptions...,
	)
	if err != nil {
		return nil, err
	}
	return bufprotoplugin.NewGenerator(
		g.logger,
		handler,
	).Generate(
		ctx,
		container,
		requests,
	)
}

type generateOptions struct {
	pluginPath []string
	protocPath []string
}

func newGenerateOptions() *generateOptions {
	return &generateOptions{}
}
