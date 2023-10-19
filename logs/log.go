// Package logs
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package logs

import (
	"context"
)

type ctxKeyLogIDType string

const (
	CtxKeyLogID ctxKeyLogIDType = "ctx_key_log_id"
)

func CtxInfo(ctx context.Context, format string, args ...any) {
	v := ctx.Value(CtxKeyLogID)
	logger.WithContext(ctx).WithField(string(CtxKeyLogID), v).Infof(format, args...)
}

func CtxError(ctx context.Context, format string, args ...any) {
	v := ctx.Value(CtxKeyLogID)
	logger.WithContext(ctx).WithField(string(CtxKeyLogID), v).Errorf(format, args...)
}

func CtxWarn(ctx context.Context, format string, args ...any) {
	v := ctx.Value(CtxKeyLogID)
	logger.WithContext(ctx).WithField(string(CtxKeyLogID), v).Warnf(format, args...)
}
