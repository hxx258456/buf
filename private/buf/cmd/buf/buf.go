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

package buf

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"buf.build/go/app"
	"buf.build/go/app/appcmd"
	"buf.build/go/app/appext"
	"connectrpc.com/connect"
	"github.com/bufbuild/buf/private/buf/bufcli"
	"github.com/bufbuild/buf/private/buf/bufctl"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/alpha/protoc"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/alpha/registry/token/tokendelete"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/alpha/registry/token/tokenget"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/alpha/registry/token/tokenlist"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/bufpluginv1"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/bufpluginv1beta1"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/bufpluginv2"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/lsp"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/price"
	betaplugindelete "github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/registry/plugin/plugindelete"
	betapluginpush "github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/registry/plugin/pluginpush"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/registry/webhook/webhookcreate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/registry/webhook/webhookdelete"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/registry/webhook/webhooklist"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/beta/studioagent"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/breaking"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/build"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/config/configinit"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/config/configlsbreakingrules"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/config/configlslintrules"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/config/configlsmodules"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/config/configmigrate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/convert"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/curl"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/dep/depgraph"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/dep/depprune"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/dep/depupdate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/export"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/format"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/generate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/lint"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/lsfiles"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/mod/modlsbreakingrules"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/mod/modlslintrules"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/mod/modopen"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/plugin/pluginprune"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/plugin/pluginpush"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/plugin/pluginupdate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/push"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulecommit/modulecommitaddlabel"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulecommit/modulecommitinfo"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulecommit/modulecommitlist"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulecommit/modulecommitresolve"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulecreate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/moduledelete"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/moduledeprecate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/moduleinfo"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulelabel/modulelabelarchive"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulelabel/modulelabelinfo"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulelabel/modulelabellist"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulelabel/modulelabelunarchive"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/modulesettings/modulesettingsupdate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/module/moduleundeprecate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/organization/organizationcreate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/organization/organizationdelete"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/organization/organizationinfo"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/organization/organizationupdate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/plugincommit/plugincommitaddlabel"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/plugincommit/plugincommitinfo"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/plugincommit/plugincommitlist"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/plugincommit/plugincommitresolve"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/plugincreate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/plugindelete"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/plugininfo"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/pluginlabel/pluginlabelarchive"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/pluginlabel/pluginlabelinfo"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/pluginlabel/pluginlabellist"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/pluginlabel/pluginlabelunarchive"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/plugin/pluginsettings/pluginsettingsupdate"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/registrycc"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/registrylogin"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/registrylogout"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/sdk/version"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/registry/whoami"
	"github.com/bufbuild/buf/private/buf/cmd/buf/command/stats"
	"github.com/bufbuild/buf/private/bufpkg/bufcobra"
	"github.com/bufbuild/buf/private/bufpkg/bufconnect"
	"github.com/bufbuild/buf/private/bufpkg/bufmodule"
	"github.com/bufbuild/buf/private/pkg/slogapp"
	"github.com/bufbuild/buf/private/pkg/syserror"
	"github.com/spf13/cobra"
)

// Main is the entrypoint to the buf CLI.
func Main(name string) {
	appcmd.Main(context.Background(), NewRootCommand(name))
}

// NewRootCommand returns a new root command.
//
// This is public for use in testing.
func NewRootCommand(name string) *appcmd.Command {
	builder := appext.NewBuilder(
		name,
		appext.BuilderWithTimeout(120*time.Second),
		appext.BuilderWithInterceptor(newErrorInterceptor()),
		appext.BuilderWithLoggerProvider(slogapp.LoggerProvider),
	)
	return &appcmd.Command{
		Use:                 name,
		Short:               "The Buf CLI",
		Long:                "A tool for working with Protocol Buffers and managing resources on the Buf Schema Registry (BSR)",
		Version:             bufcli.Version,
		BindPersistentFlags: builder.BindRoot,
		SubCommands: []*appcmd.Command{
			build.NewCommand("build", builder),
			export.NewCommand("export", builder),
			format.NewCommand("format", builder),
			lint.NewCommand("lint", builder),
			breaking.NewCommand("breaking", builder),
			generate.NewCommand("generate", builder),
			lsfiles.NewCommand("ls-files", builder),
			stats.NewCommand("stats", builder),
			push.NewCommand("push", builder),
			convert.NewCommand("convert", builder),
			curl.NewCommand("curl", builder),
			{
				Use:   "dep",
				Short: "Work with dependencies",
				SubCommands: []*appcmd.Command{
					depgraph.NewCommand("graph", builder),
					depprune.NewCommand("prune", builder, ``, false),
					depupdate.NewCommand("update", builder, ``, false),
				},
			},
			{
				Use:   "config",
				Short: "Work with configuration files",
				SubCommands: []*appcmd.Command{
					configinit.NewCommand("init", builder, ``, false, false),
					configmigrate.NewCommand("migrate", builder),
					configlslintrules.NewCommand("ls-lint-rules", builder),
					configlsbreakingrules.NewCommand("ls-breaking-rules", builder),
					configlsmodules.NewCommand("ls-modules", builder),
				},
			},
			{
				Use:        "mod",
				Short:      `Manage Buf modules, all commands are deprecated and have moved to the "buf config", "buf dep", or "buf registry" subcommands.`,
				Deprecated: `all commands are deprecated and have moved to the "buf config", "buf dep", or "buf registry" subcommands.`,
				Hidden:     true,
				SubCommands: []*appcmd.Command{
					// Deprecated and hidden.
					configinit.NewCommand("init", builder, deprecatedMessage("buf config init", "buf mod init"), true, true),
					// Deprecated and hidden.
					depprune.NewCommand("prune", builder, deprecatedMessage("buf dep prune", "buf mod update"), true),
					// Deprecated and hidden.
					depupdate.NewCommand("update", builder, deprecatedMessage("buf dep update", "buf mod update"), true),
					// Deprecated and hidden.
					modopen.NewCommand("open", builder),
					// Deprecated and hidden.
					registrycc.NewCommand("clear-cache", builder, deprecatedMessage("buf registry cc", "buf mod clear-cache"), true, "cc"),
					// Deprecated and hidden.
					modlslintrules.NewCommand("ls-lint-rules", builder),
					// Deprecated and hidden.
					modlsbreakingrules.NewCommand("ls-breaking-rules", builder),
				},
			},
			{
				Use:   "plugin",
				Short: "Work with check plugins",
				SubCommands: []*appcmd.Command{
					pluginpush.NewCommand("push", builder),
					pluginupdate.NewCommand("update", builder),
					pluginprune.NewCommand("prune", builder),
				},
			},
			{
				Use:   "registry",
				Short: "Manage assets on the Buf Schema Registry",
				SubCommands: []*appcmd.Command{
					registrylogin.NewCommand("login", builder),
					registrylogout.NewCommand("logout", builder),
					whoami.NewCommand("whoami", builder),
					registrycc.NewCommand("cc", builder, ``, false),
					{
						Use:        "commit",
						Short:      `Manage a module's commits, all commands are deprecated and have moved to the "buf registry module commit" subcommands`,
						Deprecated: `all commands are deprecated and have moved to the "buf registry module commit" subcommands.`,
						Hidden:     true,
						SubCommands: []*appcmd.Command{
							modulecommitaddlabel.NewCommand("add-label", builder, deprecatedMessage("buf registry module commit add-label", "buf registry commit add-label")),
							modulecommitinfo.NewCommand("info", builder, deprecatedMessage("buf registry module commit info", "buf registry commit info")),
							modulecommitlist.NewCommand("list", builder, deprecatedMessage("buf registry module commit list", "buf registry commit list")),
							modulecommitresolve.NewCommand("resolve", builder, deprecatedMessage("buf registry module commit resolve", "buf registry commit resolve")),
						},
					},
					{
						Use:   "sdk",
						Short: "Manage Generated SDKs",
						SubCommands: []*appcmd.Command{
							version.NewCommand("version", builder),
						},
					},
					{
						Use:        "label",
						Short:      `Manage a module's labels, all commands are deprecated and have moved to the "buf registry module label" subcommands`,
						Deprecated: `all commands are deprecated and have moved to the "buf registry module label" subcommands.`,
						Hidden:     true,
						SubCommands: []*appcmd.Command{
							modulelabelarchive.NewCommand("archive", builder, deprecatedMessage("buf registry module label archive", "buf registry label archive")),
							modulelabelinfo.NewCommand("info", builder, deprecatedMessage("buf registry module label info", "buf registry label info")),
							modulelabellist.NewCommand("list", builder, deprecatedMessage("buf registry module label list", "buf registry label list")),
							modulelabelunarchive.NewCommand("unarchive", builder, deprecatedMessage("buf registry module label unarchive", "buf registry label unarchive")),
						},
					},
					{
						Use:   "organization",
						Short: "Manage organizations",
						SubCommands: []*appcmd.Command{
							organizationcreate.NewCommand("create", builder),
							organizationdelete.NewCommand("delete", builder),
							organizationinfo.NewCommand("info", builder),
							organizationupdate.NewCommand("update", builder),
						},
					},
					{
						Use:   "module",
						Short: "Manage BSR modules",
						SubCommands: []*appcmd.Command{
							{
								Use:   "commit",
								Short: "Manage a module's commits",
								SubCommands: []*appcmd.Command{
									modulecommitaddlabel.NewCommand("add-label", builder, ""),
									modulecommitinfo.NewCommand("info", builder, ""),
									modulecommitlist.NewCommand("list", builder, ""),
									modulecommitresolve.NewCommand("resolve", builder, ""),
								},
							},
							{
								Use:   "label",
								Short: "Manage a module's labels",
								SubCommands: []*appcmd.Command{
									modulelabelarchive.NewCommand("archive", builder, ""),
									modulelabelinfo.NewCommand("info", builder, ""),
									modulelabellist.NewCommand("list", builder, ""),
									modulelabelunarchive.NewCommand("unarchive", builder, ""),
								},
							},
							{
								Use:   "settings",
								Short: "Manage a module's settings",
								SubCommands: []*appcmd.Command{
									modulesettingsupdate.NewCommand("update", builder, ""),
								},
							},
							modulecreate.NewCommand("create", builder),
							moduleinfo.NewCommand("info", builder),
							moduledelete.NewCommand("delete", builder),
							moduledeprecate.NewCommand("deprecate", builder),
							modulesettingsupdate.NewCommand("update", builder, deprecatedMessage("buf registry module settings update", "buf registry update")),
							moduleundeprecate.NewCommand("undeprecate", builder),
						},
					},
					{
						Use:   "plugin",
						Short: "Manage BSR plugins",
						SubCommands: []*appcmd.Command{
							{
								Use:   "commit",
								Short: "Manage a plugin's commits",
								SubCommands: []*appcmd.Command{
									plugincommitaddlabel.NewCommand("add-label", builder, ""),
									plugincommitinfo.NewCommand("info", builder, ""),
									plugincommitlist.NewCommand("list", builder, ""),
									plugincommitresolve.NewCommand("resolve", builder, ""),
								},
							},
							{
								Use:   "label",
								Short: "Manage a plugin's labels",
								SubCommands: []*appcmd.Command{
									pluginlabelarchive.NewCommand("archive", builder, ""),
									pluginlabelinfo.NewCommand("info", builder, ""),
									pluginlabellist.NewCommand("list", builder, ""),
									pluginlabelunarchive.NewCommand("unarchive", builder, ""),
								},
							},
							{
								Use:   "settings",
								Short: "Manage a plugin's settings",
								SubCommands: []*appcmd.Command{
									pluginsettingsupdate.NewCommand("update", builder),
								},
							},
							plugincreate.NewCommand("create", builder),
							plugininfo.NewCommand("info", builder),
							plugindelete.NewCommand("delete", builder),
						},
					},
				},
			},
			{
				Use:   "beta",
				Short: "Beta commands. Unstable and likely to change",
				SubCommands: []*appcmd.Command{
					lsp.NewCommand("lsp", builder),
					price.NewCommand("price", builder),
					bufpluginv1beta1.NewCommand("buf-plugin-v1beta1", builder),
					bufpluginv1.NewCommand("buf-plugin-v1", builder),
					bufpluginv2.NewCommand("buf-plugin-v2", builder),
					studioagent.NewCommand("studio-agent", builder),
					{
						Use:   "registry",
						Short: "Manage assets on the Buf Schema Registry",
						SubCommands: []*appcmd.Command{
							{
								Use:   "webhook",
								Short: "Manage webhooks for a repository on the Buf Schema Registry",
								SubCommands: []*appcmd.Command{
									webhookcreate.NewCommand("create", builder),
									webhookdelete.NewCommand("delete", builder),
									webhooklist.NewCommand("list", builder),
								},
							},
							{
								Use:   "plugin",
								Short: "Manage plugins on the Buf Schema Registry",
								SubCommands: []*appcmd.Command{
									betapluginpush.NewCommand("push", builder),
									betaplugindelete.NewCommand("delete", builder),
								},
							},
						},
					},
				},
			},
			{
				Use:    "alpha",
				Short:  "Alpha commands. Unstable and recommended only for experimentation. These may be deleted",
				Hidden: true,
				SubCommands: []*appcmd.Command{
					protoc.NewCommand("protoc", builder),
					{
						Use:   "registry",
						Short: "Manage assets on the Buf Schema Registry",
						SubCommands: []*appcmd.Command{
							{
								Use:   "token",
								Short: "Manage user tokens",
								SubCommands: []*appcmd.Command{
									tokenget.NewCommand("get", builder),
									tokenlist.NewCommand("list", builder),
									tokendelete.NewCommand("delete", builder),
								},
							},
						},
					},
				},
			},
		},
		ModifyCobra: func(cobraCommand *cobra.Command) error {
			cobraCommand.AddCommand(bufcobra.NewWebpagesCommand("webpages", cobraCommand))
			return nil
		},
	}
}

// newErrorInterceptor returns a CLI interceptor that wraps Buf CLI errors.
func newErrorInterceptor() appext.Interceptor {
	return func(next func(context.Context, appext.Container) error) func(context.Context, appext.Container) error {
		return func(ctx context.Context, container appext.Container) error {
			return wrapError(next(ctx, container))
		}
	}
}

// wrapError is used when a CLI command fails, regardless of its error code.
// Note that this function will wrap the error so that the underlying error
// can be recovered via 'errors.Is'.
func wrapError(err error) error {
	if err == nil {
		return nil
	}

	var connectErr *connect.Error
	isConnectError := errors.As(err, &connectErr)
	// If error is empty and not a system error or Connect error, we return it as-is.
	if !isConnectError && err.Error() == "" {
		return err
	}
	if isConnectError {
		var augmentedConnectError *bufconnect.AugmentedConnectError
		isAugmentedConnectErr := errors.As(err, &augmentedConnectError)
		if isPossibleNewCLIOldBSRError(connectErr) && isAugmentedConnectErr {
			return fmt.Errorf("Failure: %[1]s for https://%[2]s%[3]s\n"+
				"This version of the buf CLI may require APIs that have not yet been deployed to https://%[2]s\n"+
				"To resolve this failure, you can either:\n"+
				"- Try using an older version of the buf CLI\n"+
				"- Contact the site admin for https://%[2]s to upgrade the instance",
				connectErr,
				augmentedConnectError.Addr(),
				augmentedConnectError.Procedure(),
			)
		}
		connectCode := connectErr.Code()
		switch {
		case connectCode == connect.CodeUnauthenticated, isEmptyUnknownError(err):
			loginCommand := "buf registry login"
			authErr, ok := bufconnect.AsAuthError(err)
			if !ok {
				// This code should be unreachable.
				return fmt.Errorf("Failure: you are not authenticated. "+
					"Set the %[1]s environment variable or run %q, using a Buf API token as the password. "+
					"If you have set the %[1]s or run the login command, "+
					"your token may have expired. "+
					"For details, visit https://buf.build/docs/bsr/authentication",
					bufconnect.TokenEnvKey,
					loginCommand,
				)
			}
			// Invalid token found in env var.
			if authErr.HasToken() && authErr.TokenEnvKey() != "" {
				return fmt.Errorf("Failure: the %[1]s environment variable is not valid for %[2]s. "+
					"Set %[1]s to a valid Buf API token, or unset it. "+
					"For details, visit https://buf.build/docs/bsr/authentication",
					authErr.TokenEnvKey(), authErr.Remote(),
				)
			}
			if authErr.Remote() != bufconnect.DefaultRemote {
				loginCommand = fmt.Sprintf("%s %s", loginCommand, authErr.Remote())
			}
			// Invalid token found in netrc.
			if authErr.HasToken() {
				return fmt.Errorf("Failure: your Buf API token for %s is invalid. "+
					"Run %q using a valid Buf API token. "+
					"For details, visit https://buf.build/docs/bsr/authentication",
					authErr.Remote(),
					loginCommand,
				)
			}
			// No token found.
			return fmt.Errorf("Failure: you are not authenticated for %s. "+
				"Set the %s environment variable or run %q, "+
				"using a Buf API token as the password. "+
				"For details, visit https://buf.build/docs/bsr/authentication",
				authErr.Remote(),
				bufconnect.TokenEnvKey,
				loginCommand,
			)
		case connectCode == connect.CodeUnavailable:
			msg := `Failure: the server hosted at that remote is unavailable.`
			// If the returned error is Unavailable, then determine if this is a DNS error.  If so,
			// get the address used so that we can display a more helpful error message.
			if dnsError := (&net.DNSError{}); errors.As(err, &dnsError) && dnsError.IsNotFound {
				return fmt.Errorf(`%s Are you sure "%s" is a valid remote address?`, msg, dnsError.Name)
			}
			// If the unavailable error wraps a tls.CertificateVerificationError, show a more specific
			// error message to the user to aid in troubleshooting.
			if tlsErr := wrappedTLSError(err); tlsErr != nil {
				return fmt.Errorf("tls certificate verification: %w", tlsErr)
			}
			return errors.New(msg)
		}
	}

	sysError, isSysError := syserror.As(err)
	if isSysError {
		err = fmt.Errorf(
			"it looks like you have found a bug in buf. "+
				"Please file an issue at https://github.com/bufbuild/buf/issues "+
				"and provide the command you ran, as well as the following message: %w",
			sysError.Unwrap(),
		)
	}

	var importNotExistError *bufmodule.ImportNotExistError
	if errors.As(err, &importNotExistError) {
		// There must be a better place to do this, perhaps in the Controller, but this works for now.
		err = app.WrapError(bufctl.ExitCodeFileAnnotation, importNotExistError)
	}

	return appFailureError(err)
}

// isEmptyUnknownError returns true if the given
// error is non-nil, but has an empty message
// and an unknown error code.
//
// This is relevant for errors returned by
// envoyauthd when the client does not provide
// an authentication header.
func isEmptyUnknownError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "" && connect.CodeOf(err) == connect.CodeUnknown
}

// wrappedTLSError returns an unwrapped TLS error or nil if the error is another type of error.
func wrappedTLSError(err error) error {
	if tlsErr := (&tls.CertificateVerificationError{}); errors.As(err, &tlsErr) {
		return tlsErr
	}
	return nil
}

func appFailureError(err error) error {
	return fmt.Errorf("Failure: %w", err)
}

// isPossibleNewCLIOldBSRError determines if an error might be from a newer
// version of the CLI interacting with an older version of the BSR.
func isPossibleNewCLIOldBSRError(connectErr *connect.Error) bool {
	switch connectErr.Code() {
	case connect.CodeUnknown:
		// Older versions of the BSR return errors of this shape
		// for unrecognized services.
		// NOTE: This handling can be removed once all BSR instances
		// are upgraded past v1.7.0.
		return connectErr.Message() == fmt.Sprintf("%d %s", http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	case connect.CodeUnimplemented:
		// RPC was known, but unimplemented in the BSR version.
		return true
	default:
		return false
	}
}

// deprecatedMessage returns a message indicating that a command is deprecated.
func deprecatedMessage(newCommand, oldCommand string) string {
	return fmt.Sprintf(
		`use "%s" instead. However, "%s" will continue to work.`,
		newCommand, oldCommand,
	)
}
