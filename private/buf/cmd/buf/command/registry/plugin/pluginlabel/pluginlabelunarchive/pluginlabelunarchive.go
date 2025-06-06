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

package pluginlabelunarchive

import (
	"context"
	"fmt"

	pluginv1beta1 "buf.build/gen/go/bufbuild/registry/protocolbuffers/go/buf/registry/plugin/v1beta1"
	"buf.build/go/app/appcmd"
	"buf.build/go/app/appext"
	"connectrpc.com/connect"
	"github.com/bufbuild/buf/private/buf/bufcli"
	"github.com/bufbuild/buf/private/bufpkg/bufparse"
	"github.com/bufbuild/buf/private/bufpkg/bufregistryapi/bufregistryapiplugin"
	"github.com/spf13/pflag"
)

// NewCommand returns a new Command.
func NewCommand(
	name string,
	builder appext.SubCommandBuilder,
	deprecated string,
) *appcmd.Command {
	flags := newFlags()
	return &appcmd.Command{
		Use:        name + " <remote/owner/plugin:label>",
		Short:      "Unarchive a plugin label",
		Args:       appcmd.ExactArgs(1),
		Deprecated: deprecated,
		Run: builder.NewRunFunc(
			func(ctx context.Context, container appext.Container) error {
				return run(ctx, container, flags)
			},
		),
		BindFlags: flags.Bind,
	}
}

type flags struct {
}

func newFlags() *flags {
	return &flags{}
}

func (f *flags) Bind(flagSet *pflag.FlagSet) {
}

func run(
	ctx context.Context,
	container appext.Container,
	_ *flags,
) error {
	pluginRef, err := bufparse.ParseRef(container.Arg(0))
	if err != nil {
		return appcmd.WrapInvalidArgumentError(err)
	}
	labelName := pluginRef.Ref()
	if labelName == "" {
		return appcmd.NewInvalidArgumentError("label is required")
	}
	clientConfig, err := bufcli.NewConnectClientConfig(container)
	if err != nil {
		return err
	}
	pluginFullName := pluginRef.FullName()
	labelServiceClient := bufregistryapiplugin.NewClientProvider(clientConfig).V1Beta1LabelServiceClient(pluginFullName.Registry())
	// UnarchiveLabelsResponse is empty.
	if _, err := labelServiceClient.UnarchiveLabels(
		ctx,
		connect.NewRequest(
			&pluginv1beta1.UnarchiveLabelsRequest{
				LabelRefs: []*pluginv1beta1.LabelRef{
					{
						Value: &pluginv1beta1.LabelRef_Name_{
							Name: &pluginv1beta1.LabelRef_Name{
								Owner:  pluginFullName.Owner(),
								Plugin: pluginFullName.Name(),
								Label:  labelName,
							},
						},
					},
				},
			},
		),
	); err != nil {
		if connect.CodeOf(err) == connect.CodeNotFound {
			return bufcli.NewLabelNotFoundError(pluginRef)
		}
		return err
	}
	_, err = fmt.Fprintf(container.Stdout(), "Unarchived %s.\n", pluginRef)
	return err
}
