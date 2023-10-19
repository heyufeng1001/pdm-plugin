// Package client
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package client

import (
	"context"
	"strings"

	"github.com/bytedance/sonic"
	"pdm-plugin.github.com/biz/dao/filter"
	"pdm-plugin.github.com/logs"
)

type BaseEntry struct {
	AppID         string `json:"app_id"`
	EntryID       string `json:"entry_id"`
	TransactionID string `json:"transaction_id"`
}

type EntryValue struct {
	Value any `json:"value"`
}

func NewEntryValue(v any) EntryValue {
	return EntryValue{Value: v}
}

type EntryItem map[string]any

type JsonEntryItem string

func (e EntryItem) JSON(ctx context.Context) JsonEntryItem {
	j, _ := sonic.MarshalString(e)
	logs.CtxInfo(ctx, "[Filter.JSON]generate filter json: %s", j)
	return JsonEntryItem(j)
}

func (e EntryItem) Get(key string) (string, bool) {
	v, ok := e[key]
	if ok {
		r, _ := v.(string)
		return r, true
	}
	return "", false
}

func (e EntryItem) GetAny(key string) (any, bool) {
	v, ok := e[key]
	if ok {
		return v, true
	}
	return nil, false
}

func (e EntryItem) GetChildren(key string) ([]EntryItem, bool) {
	v, ok := e[key]
	if !ok {
		return nil, false
	}
	vms, ok := v.([]any)
	if !ok {
		return nil, false
	}
	ret := []EntryItem{}
	for _, vm := range vms {
		m, ok := vm.(map[string]any)
		if !ok {
			continue
		}
		ret = append(ret, m)
	}
	return ret, true
}

func (e EntryItem) Get0(key string) (string, bool) {
	widgets := strings.Split(key, ".")
	switch len(widgets) {
	case 1:
		v, ok := e[widgets[0]]
		if ok {
			return v.(string), true
		}
		return "", false
	case 2:
		v, ok := e[widgets[0]]
		if !ok {
			return "", false
		}
		vms, ok := v.([]any)
		if !ok {
			return "", false
		}
		for _, vm := range vms {
			m, ok := vm.(map[string]any)
			if !ok {
				continue
			}
			r, ok := m[widgets[1]]
			if ok {
				return r.(string), true
			}
		}
		return "", false
	default:
	}
	return "", false
}

func (e EntryItem) MGet(key string) string {
	s, ok := e.Get(key)
	if !ok {
		return ""
	}
	return s
}

func (e EntryItem) MGetAny(key string) any {
	v, ok := e.GetAny(key)
	if !ok {
		return nil
	}
	return v
}

type QueryParam struct {
	BaseEntry
	DataID string         `json:"data_id,omitempty"`
	Fields []string       `json:"fields,omitempty"`
	Filter *filter.Filter `json:"filter,omitempty"`
	Limit  int            `json:"limit"`
}

type CreateParam struct {
	BaseEntry
	Data EntryItem `json:"data"`
}

type BatchCreateParam struct {
	BaseEntry
	DataList []EntryItem `json:"data_list"`
}

type UpdateParam struct {
	BaseEntry
	DataID string    `json:"data_id"`
	Data   EntryItem `json:"data"`
}

type BatchUpdateParam struct {
	BaseEntry
	DataIDs []string  `json:"data_ids"`
	Data    EntryItem `json:"data"`
}

type GetParam struct {
	BaseEntry
	DataID string `json:"data_id"`
}

type BatchDeleteParam struct {
	BaseEntry
	DataIDs []string `json:"data_ids"`
}

type QueryResp struct {
	Data []EntryItem `json:"data"`
}

type GetResp struct {
	Data EntryItem `json:"data"`
}
