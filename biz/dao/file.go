// Package dao
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/8
package dao

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/bytedance/sonic"
	"pdm-plugin.github.com/config"
	"pdm-plugin.github.com/logs"
	"pdm-plugin.github.com/utils"
)

type GetUploadTokenResp struct {
	TokenAndURLList []*GetUploadTokenItem `json:"token_and_url_list"`
}

type GetUploadTokenItem struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type UploadResp struct {
	Key string `json:"key"`
}

func UploadImage(ctx context.Context, files []*utils.Pair[string, string], txID string, entryID string) ([]string, error) {
	logs.CtxInfo(ctx, "[UploadImage]start: %s, %s", utils.MustMarshal(files), txID)
	tokens, err := utils.Access[GetUploadTokenResp](
		ctx,
		config.Config().JianDaoHost+config.Config().URLGetUploadToken,
		http.MethodPost,
		map[string]string{
			"Content-Type": "application/json",
		},
		nil,
		map[string]string{
			"app_id":         config.Config().AppID,
			"entry_id":       entryID,
			"transaction_id": txID,
		},
		func(request *http.Request) {
			request.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		},
		func(response *http.Response) bool {
			return response.StatusCode == http.StatusOK
		},
	)
	if err != nil {
		logs.CtxError(ctx, "[UploadImage]get tokens failed: %s", err)
		return nil, err
	}
	if len(tokens.TokenAndURLList) == 0 {
		logs.CtxError(ctx, "[UploadImage]tokens is empty")
		return nil, errors.New("token is empty")
	}
	token := tokens.TokenAndURLList[0]
	keys := []string{}

	for idx, file := range files {
		logs.CtxInfo(ctx, "[UploadImage]ready to handle: %d, %s", idx, utils.MustMarshal(file))
		fromURL, fileName := file.Split()
		image, err := http.Get(fromURL)
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]get img failed: %d, %s", idx, err)
			return nil, err
		}

		// 上传文件
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		err = writer.WriteField("token", token.Token)
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]write field token failed: %d, %s", idx, err)
			return nil, err
		}
		fw, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]create field file failed: %d, %s", idx, err)
			return nil, err
		}
		_, err = io.Copy(fw, image.Body)
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]copy failed: %d, %s", idx, err)
			return nil, err
		}

		err = writer.Close()
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]close writer failed: %d, %s", idx, err)
			return nil, err
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, token.URL, payload)
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]new http req failed: %d, %s", idx, err)
			return nil, err
		}
		req.Header.Add("Authorization", "Bearer "+config.Config().AppSecret)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		logs.CtxInfo(ctx, "[UploadImage]ready to send http request, url: %s, method: %s, header: %v", req.URL, req.Method, req.Header)

		res, err := client.Do(req)
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]do http failed: %d, %s", idx, err)
			return nil, err
		}

		logs.CtxInfo(ctx, "[UploadImage]get http response, code: %d, header: %v", res.StatusCode, res.Header)

		body, err := io.ReadAll(res.Body)
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]read all failed: %d, %s", idx, err)
			return nil, err
		}
		logs.CtxInfo(ctx, "[UploadImage]read resp body: %s", string(body))

		ret := UploadResp{}
		err = sonic.Unmarshal(body, &ret)
		if err != nil {
			logs.CtxError(ctx, "[UploadImage]unmarshal failed: %d, %s", idx, err)
			return nil, err
		}
		logs.CtxInfo(ctx, "[UploadImage]ready to append key: %d, %s", idx, ret.Key)
		keys = append(keys, ret.Key)
	}
	logs.CtxInfo(ctx, "[UploadImage]success: %v", keys)
	return keys, nil
}
