// Package client
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package client

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
	"pdm-plugin.github.com/config"
	"pdm-plugin.github.com/utils"
)

// IDataClient 对于数据存在三个操作：查询，更新，新建
// 具体来说：在设计/任务表中查询，在任务表中新增/更新
//
//go:generate mockgen -destination api_mock.go -package client . IDataClient
type IDataClient interface {
	Query(ctx context.Context, body *QueryParam) (QueryResp, error)
	Create(ctx context.Context, body *CreateParam) error
	BatchCreate(ctx context.Context, body *BatchCreateParam) error
	Update(ctx context.Context, body *UpdateParam) error
	BatchUpdate(ctx context.Context, body *BatchUpdateParam) error
	BatchDelete(ctx context.Context, body *BatchDeleteParam) error
	Get(ctx context.Context, body *GetParam) (GetResp, error)
}

func NewIDataClient() IDataClient {
	return &dataClient{}
}

type dataClient struct {
}

var buLimiter = rate.NewLimiter(10, 10)

func (d *dataClient) BatchUpdate(ctx context.Context, body *BatchUpdateParam) error {
	buLimiter.Wait(ctx)
	_, err := utils.Access[any](
		ctx,
		config.Config().JianDaoHost+config.Config().URLBatchUpdateData,
		http.MethodPost,
		map[string]string{
			"Content-Type": "application/json",
		},
		nil,
		body,
		func(request *http.Request) {
			request.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		},
		func(response *http.Response) bool {
			return response.StatusCode == http.StatusOK
		},
	)
	return err
}

var bcLimiter = rate.NewLimiter(10, 10)

func (d *dataClient) BatchCreate(ctx context.Context, body *BatchCreateParam) error {
	bcLimiter.Wait(ctx)
	_, err := utils.Access[any](
		ctx,
		config.Config().JianDaoHost+config.Config().URLBatchCreateData,
		http.MethodPost,
		map[string]string{
			"Content-Type": "application/json",
		},
		nil,
		body,
		func(request *http.Request) {
			request.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		},
		func(response *http.Response) bool {
			return response.StatusCode == http.StatusOK
		},
	)
	return err
}

var bdLimiter = rate.NewLimiter(10, 10)

func (d *dataClient) BatchDelete(ctx context.Context, body *BatchDeleteParam) error {
	bdLimiter.Wait(ctx)
	_, err := utils.Access[any](
		ctx,
		config.Config().JianDaoHost+config.Config().URLBatchDeleteData,
		http.MethodPost,
		map[string]string{
			"Content-Type": "application/json",
		},
		nil,
		body,
		func(request *http.Request) {
			request.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		},
		func(response *http.Response) bool {
			return response.StatusCode == http.StatusOK
		},
	)
	return err
}

var cLimiter = rate.NewLimiter(20, 20)

func (d *dataClient) Create(ctx context.Context, body *CreateParam) error {
	cLimiter.Wait(ctx)
	_, err := utils.Access[any](
		ctx,
		config.Config().JianDaoHost+config.Config().URLCreateData,
		http.MethodPost,
		map[string]string{
			"Content-Type": "application/json",
		},
		nil,
		body,
		func(request *http.Request) {
			request.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		},
		func(response *http.Response) bool {
			return response.StatusCode == http.StatusOK
		},
	)
	return err
}

var uLimiter = rate.NewLimiter(10, 10)

func (d *dataClient) Update(ctx context.Context, body *UpdateParam) error {
	uLimiter.Wait(ctx)
	_, err := utils.Access[any](
		ctx,
		config.Config().JianDaoHost+config.Config().URLUpdateData,
		http.MethodPost,
		map[string]string{
			"Content-Type": "application/json",
		},
		nil,
		body,
		func(request *http.Request) {
			request.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		},
		func(response *http.Response) bool {
			return response.StatusCode == http.StatusOK
		},
	)
	return err
}

var qLimiter = rate.NewLimiter(20, 20)

func (d *dataClient) Query(ctx context.Context, body *QueryParam) (QueryResp, error) {
	qLimiter.Wait(ctx)
	return utils.Access[QueryResp](
		ctx,
		config.Config().JianDaoHost+config.Config().URLQueryData,
		http.MethodPost,
		map[string]string{
			"Content-Type": "application/json",
		},
		nil,
		body,
		func(request *http.Request) {
			request.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		},
		func(response *http.Response) bool {
			return response.StatusCode == http.StatusOK
		},
	)
}

var gLimiter = rate.NewLimiter(30, 20)

func (d *dataClient) Get(ctx context.Context, body *GetParam) (GetResp, error) {
	gLimiter.Wait(ctx)
	return utils.Access[GetResp](
		ctx,
		config.Config().JianDaoHost+config.Config().URLGetData,
		http.MethodPost,
		map[string]string{
			"Content-Type": "application/json",
		},
		nil,
		body,
		func(request *http.Request) {
			request.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		},
		func(response *http.Response) bool {
			return response.StatusCode == http.StatusOK
		},
	)
}
