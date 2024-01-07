// Package handler
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/12/29
package handler

import (
	"context"
	"testing"

	"pdm-plugin.github.com/config"
)

func TestSync(t *testing.T) {
	config.Init("../../conf/conf.yaml")
	syncFullData(context.Background())
}
