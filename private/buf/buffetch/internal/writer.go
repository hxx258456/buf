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

package internal

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"buf.build/go/app"
	"buf.build/go/standard/xio"
	"github.com/klauspost/compress/zstd"
)

type writer struct {
	logger *slog.Logger

	// never set for now (no corresponding option)
	httpEnabled  bool
	localEnabled bool
	stdioEnabled bool
}

func newWriter(
	logger *slog.Logger,
	options ...WriterOption,
) *writer {
	writer := &writer{
		logger: logger,
	}
	for _, option := range options {
		option(writer)
	}
	return writer
}

func (w *writer) PutFile(
	ctx context.Context,
	container app.EnvStdoutContainer,
	fileRef FileRef,
	options ...PutFileOption,
) (io.WriteCloser, error) {
	putFileOptions := newPutFileOptions()
	for _, option := range options {
		option(putFileOptions)
	}
	switch t := fileRef.(type) {
	case SingleRef:
		return w.putSingle(
			ctx,
			container,
			t,
			putFileOptions.noFileCompression,
		)
	case ArchiveRef:
		return w.putArchiveFile(
			ctx,
			container,
			t,
			putFileOptions.noFileCompression,
		)
	default:
		return nil, fmt.Errorf("unknown FileRef type: %T", fileRef)
	}
}

func (w *writer) putSingle(
	ctx context.Context,
	container app.EnvStdoutContainer,
	singleRef SingleRef,
	noFileCompression bool,
) (io.WriteCloser, error) {
	return w.putFileWriteCloser(ctx, container, singleRef, noFileCompression)
}

func (w *writer) putArchiveFile(
	ctx context.Context,
	container app.EnvStdoutContainer,
	archiveRef ArchiveRef,
	noFileCompression bool,
) (io.WriteCloser, error) {
	return w.putFileWriteCloser(ctx, container, archiveRef, noFileCompression)
}

func (w *writer) putFileWriteCloser(
	ctx context.Context,
	container app.EnvStdoutContainer,
	fileRef FileRef,
	noFileCompression bool,
) (_ io.WriteCloser, retErr error) {
	writeCloser, err := w.putFileWriteCloserPotentiallyUncompressed(ctx, container, fileRef)
	if err != nil {
		return nil, err
	}
	defer func() {
		if retErr != nil {
			retErr = errors.Join(retErr, writeCloser.Close())
		}
	}()
	if noFileCompression {
		return writeCloser, nil
	}
	switch compressionType := fileRef.CompressionType(); compressionType {
	case CompressionTypeNone:
		return writeCloser, nil
	case CompressionTypeGzip:
		gzipWriteCloser := gzip.NewWriter(writeCloser)
		return xio.CompositeWriteCloser(
			gzipWriteCloser,
			xio.ChainCloser(
				gzipWriteCloser,
				writeCloser,
			),
		), nil
	case CompressionTypeZstd:
		zstdWriteCloser, err := zstd.NewWriter(writeCloser)
		if err != nil {
			return nil, err
		}
		return xio.CompositeWriteCloser(
			zstdWriteCloser,
			xio.ChainCloser(
				zstdWriteCloser,
				writeCloser,
			),
		), nil
	default:
		return nil, fmt.Errorf("unknown CompressionType: %v", compressionType)
	}
}

func (w *writer) putFileWriteCloserPotentiallyUncompressed(
	ctx context.Context,
	container app.EnvStdoutContainer,
	fileRef FileRef,
) (io.WriteCloser, error) {
	switch fileScheme := fileRef.FileScheme(); fileScheme {
	case FileSchemeHTTP:
		if !w.httpEnabled {
			return nil, NewWriteHTTPDisabledError()
		}
		return nil, fmt.Errorf("http not supported for writes: %v", fileRef.Path())
	case FileSchemeHTTPS:
		if !w.httpEnabled {
			return nil, NewWriteHTTPDisabledError()
		}
		return nil, fmt.Errorf("https not supported for writes: %v", fileRef.Path())
	case FileSchemeLocal:
		if !w.localEnabled {
			return nil, NewWriteLocalDisabledError()
		}
		return os.Create(fileRef.Path())
	case FileSchemeStdio, FileSchemeStdout:
		if !w.stdioEnabled {
			return nil, NewWriteStdioDisabledError()
		}
		return xio.NopWriteCloser(container.Stdout()), nil
	case FileSchemeStdin:
		return nil, errors.New("cannot write to stdin")
	case FileSchemeNull:
		return xio.DiscardWriteCloser, nil
	default:
		return nil, fmt.Errorf("unknown FileScheme: %v", fileScheme)
	}
}

type putFileOptions struct {
	noFileCompression bool
}

func newPutFileOptions() *putFileOptions {
	return &putFileOptions{}
}
