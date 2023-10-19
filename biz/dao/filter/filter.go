// Package filter
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package filter

import (
	"context"

	"github.com/bytedance/sonic"
	"pdm-plugin.github.com/logs"
)

type Filter struct {
	Rel  Relation `json:"rel"`
	Cond []*Cond  `json:"cond"`
}

type JSONFilter string

type Relation string

const (
	RelationAnd = "and"
	RelationOr  = "or"
)

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) WithCond(cond *Cond) *Filter {
	f.Cond = append(f.Cond, cond)
	return f
}

func (f *Filter) WithRel(rel Relation) *Filter {
	f.Rel = rel
	return f
}

func (f *Filter) JSON(ctx context.Context) JSONFilter {
	j, _ := sonic.MarshalString(f)
	logs.CtxInfo(ctx, "[Filter.JSON]generate filter json: %s", j)
	return JSONFilter(j)
}

type Cond struct {
	Field  string     `json:"field"`
	Type   string     `json:"type"`
	Method CondMethod `json:"method"`
	Value  []any      `json:"value"`
}

type CondMethod string

const (
	CondMethodEq = "eq"
)
