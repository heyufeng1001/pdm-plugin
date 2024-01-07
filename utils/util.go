// Package utils
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package utils

import (
	"sync"

	"github.com/bytedance/sonic"
)

func MustMarshal(v any) string {
	s, _ := sonic.MarshalString(v)
	return s
}

var GlobalLock sync.Mutex
