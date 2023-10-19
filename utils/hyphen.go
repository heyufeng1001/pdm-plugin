// Package utils
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/4/17
package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"golang.org/x/time/rate"
	"pdm-plugin.github.com/logs"
)

var client = &http.Client{}

func AccessBySystemCall[T any](ctx context.Context, url string, method string, header map[string]string,
	param map[string]string, body any) (T, error) {
	var ret T
	sb := bytes.Buffer{}
	sb.WriteString(url)
	sb.WriteByte('?')
	for k, v := range param {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('&')
	}
	args := []string{"--location", "--request", method, sb.String()}

	bodyJSON, err := sonic.Marshal(body)
	if err != nil {
		return ret, err
	}
	if len(bodyJSON) != 0 {
		args = append(args, []string{"--data-raw", `'` + string(bodyJSON) + `''`}...)
	}

	if header != nil {
		for k, v := range header {
			args = append(args, []string{`--header`, `'` + k + `: ` + v + `'`}...)
		}
	}

	cmd := exec.Command("curl", args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	logs.CtxInfo(ctx, "[AccessBySystemCall]ready to exec cmd: %s", cmd.String())
	err = cmd.Run()
	if err != nil {
		return ret, fmt.Errorf("[AccessBySystemCall]exec cmd failed: %w", err)
	}
	err = sonic.Unmarshal(stdout.Bytes(), &ret)
	if err != nil {
		return ret, fmt.Errorf("[AccessBySystemCall]unmarshal stdout failed: %w", err)
	}
	return ret, nil
}

func AccessResp(ctx context.Context, url string, method string, header map[string]string, param map[string]string,
	body any, setAuthorization func(*http.Request)) (*http.Response, error) {
	return access(ctx, url, method, header, param, body, setAuthorization)
}

func Access[T any](ctx context.Context, url string, method string, header map[string]string, param map[string]string,
	body any, setAuthorization func(*http.Request), isSuccess func(response *http.Response) bool) (T, error) {
	var ret T
	resp, err := access(ctx, url, method, header, param, body, setAuthorization)
	if err != nil {
		return ret, fmt.Errorf("access %s failed: %w", url, err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, fmt.Errorf("read resp body failed: %w", err)
	}

	logs.CtxInfo(ctx, "[Access]read resp body: %s", string(respBody))

	if isSuccess != nil && !isSuccess(resp) {
		return ret, fmt.Errorf("is success return false: %s", string(respBody))
	}
	err = sonic.Unmarshal(respBody, &ret)
	if err != nil {
		return ret, fmt.Errorf("unmarshal resp body failed: %w", err)
	}
	return ret, nil
}

var allLimiter = rate.NewLimiter(50, 50)

func access(ctx context.Context, url string, method string, header map[string]string, param map[string]string,
	body any, setAuthorization func(*http.Request)) (*http.Response, error) {
	allLimiter.Wait(ctx)
	bodyJSON, err := sonic.MarshalString(body)
	if err != nil {
		return nil, err
	}
	bodyReader := strings.NewReader(bodyJSON)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	if setAuthorization != nil {
		setAuthorization(req)
	}
	q := req.URL.Query()
	for k, v := range param {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range header {
		req.Header.Add(k, v)
	}

	logs.CtxInfo(ctx, "[access]ready to send http request, url: %s, method: %s, header: %v, body: %s", req.URL, req.Method, req.Header, bodyJSON)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("access %s failed: %w", url, err)
	}

	logs.CtxInfo(ctx, "[access]get http response, code: %d, header: %v", resp.StatusCode, resp.Header)

	return resp, nil
}

func TernaryForm[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

type Pair[F, S any] struct {
	First  F
	Second S
}

func MakePair[F, S any](first F, second S) *Pair[F, S] {
	return &Pair[F, S]{First: first, Second: second}
}

func (p *Pair[F, S]) Split() (F, S) {
	return p.First, p.Second
}

type ch[T any] chan T

type SafeChan[T any] struct {
	ch[T]
	once sync.Once
}

func NewSafeChan[T any](size ...int) *SafeChan[T] {
	return &SafeChan[T]{TernaryForm(len(size) == 0, make(chan T), make(chan T, size[0])), sync.Once{}}
}

func (s *SafeChan[T]) Close() {
	s.once.Do(func() {
		close(s.ch)
	})
}

func BatchGetObjectWithOrder[K, V any](ctx context.Context, ids []K, minParallelNum int, fc func(context.Context, K) (V, error)) ([]V, error) {
	logs.CtxInfo(ctx, "[BatchGetObjectWithOrder]batch get start, ids is: %v", ids)
	if len(ids) < minParallelNum {
		logs.CtxInfo(ctx, "[BatchGetObjectWithOrder]start serial get")
		ret := []V{}
		for i, id := range ids {
			v, err := fc(ctx, id)
			if err != nil {
				logs.CtxError(ctx, "[BatchGetObjectWithOrder][serial]get v by fc in %d failed: %s", i, err)
				return nil, err
			}
			ret = append(ret, v)
		}
		return ret, nil
	}
	logs.CtxInfo(ctx, "[BatchGetObjectWithOrder]start parallel get")
	vOrderChan, errChan, wg := make(chan *Pair[int, V]), make(chan error), sync.WaitGroup{}
	ret := make([]V, len(ids))
	var err error
	for i, id := range ids {
		ni, nid := i, id
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					e := fmt.Errorf("[BatchGetObjectWithOrder][parallel]panic happens in batch get: %v", r)
					logs.CtxError(ctx, e.Error())
					errChan <- e
				}
			}()
			v, e := fc(ctx, nid)
			if e != nil {
				logs.CtxError(ctx, "[BatchGetObjectWithOrder][parallel]get v by fc in %d failed: %s", ni, e)
				errChan <- e
				return
			}
			vOrderChan <- MakePair(ni, v)
			logs.CtxInfo(ctx, "[BatchGetObjectWithOrder]fc %d success", ni)
		}()
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("[BatchGetObjectWithOrder][parallel]panic happens in wg wait: %v", r)
				logs.CtxError(ctx, err.Error())
			}
		}()
		wg.Wait()
		close(vOrderChan)
		close(errChan)
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("[BatchGetObjectWithOrder][parallel]panic happens in err listen get: %v", r)
				logs.CtxError(ctx, err.Error())
			}
		}()
		for e := range errChan {
			if err == nil {
				// 只设置第一个捕获的错误
				logs.CtxError(ctx, "[BatchGetObjectWithOrder]err %s was caught", err)
				err = e
			}
		}
	}()
	for vo := range vOrderChan {
		ret[vo.First] = vo.Second
	}
	if err != nil {
		logs.CtxError(ctx, "[BatchGetObjectWithOrder][parallel]get v by fc failed: %s", err)
		return nil, err
	}
	logs.CtxInfo(ctx, "[BatchGetObjectWithOrder]success")
	return ret, nil
}

type SliceSet[K comparable, V any] struct {
	m     map[K]int
	slice []V
}

func NewSetSlice[K comparable, V any]() *SliceSet[K, V] {
	return &SliceSet[K, V]{map[K]int{}, []V{}}
}

func NewSetSliceFormSlice[K comparable](s []K) *SliceSet[K, K] {
	ret := &SliceSet[K, K]{map[K]int{}, []K{}}
	for _, k := range s {
		ret.Upsert(k, k)
	}
	return ret
}

func (s *SliceSet[K, V]) insert(key K, value V) {
	s.m[key] = len(s.slice)
	s.slice = append(s.slice, value)
}

func (s *SliceSet[K, V]) update(key K, value V) {
	s.slice[s.m[key]] = value
}

func (s *SliceSet[K, V]) Insert(key K, value V) bool {
	if _, ok := s.m[key]; ok {
		return false
	}
	s.insert(key, value)
	return true
}

func (s *SliceSet[K, V]) Update(key K, value V) bool {
	if _, ok := s.m[key]; !ok {
		return false
	}
	s.update(key, value)
	return true
}
func (s *SliceSet[K, V]) Upsert(key K, value V) {
	if _, ok := s.m[key]; ok {
		s.update(key, value)
		return
	}
	s.insert(key, value)
}

func (s *SliceSet[K, V]) Get(key K) (V, bool) {
	var v V
	i, ok := s.m[key]
	if !ok {
		return v, ok
	}
	return s.slice[i], ok
}

func (s *SliceSet[K, V]) GetSlice() []V {
	return s.slice
}

func (s *SliceSet[K, V]) GetMap() map[K]int {
	return s.m
}

func SimpleCopy[F, T any](from F) (T, error) {
	var t T
	f, err := sonic.Marshal(from)
	if err != nil {
		return t, err
	}
	err = sonic.Unmarshal(f, t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func DeepCopyByJson[T any](from T) (T, error) {
	var ret T
	b, err := json.Marshal(from)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(b, &ret)
	return ret, err
}

// MustSimpleCopy 这两个方法都是不安全且有性能代价的，慎用。
func MustSimpleCopy[F, T any](from F) T {
	t := SafeInit[T]()
	t, _ = SimpleCopy[F, T](from)
	return t
}

func MustGetStringFromPtr(ptr *string) string {
	return TernaryForm(ptr == nil, "", *ptr)
}

func ParseGitURL(url string) string {
	us := strings.Split(url, "/")
	if len(us) != 5 && len(us) != 6 {
		return ""
	}
	if us[2] == "code.byted.org" {
		if len(us[4]) <= 4 {
			return ""
		}
		return fmt.Sprintf("%s/%s", us[3], us[4][:len(us[4])-4])
	} else if us[2] == "review.byted.org" {
		if len(us) == 6 {
			return fmt.Sprintf("%s/%s/%s", us[3], us[4], us[5])
		}
		return fmt.Sprintf("%s/%s", us[3], us[4])
	} else {
		return ""
	}
}

func MustAtoi(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

func MustDo[K, V any](key K, fc func(K) (V, error)) V {
	return MustEasyDo(func() (V, error) {
		return fc(key)
	})
}

func MustEasyDo[V any](fc func() (V, error)) V {
	v, err := fc()
	if err != nil {
		return SafeInit[V]()
	}
	return v
}

func Paging[T any](arr []T, offset, limit int) []T {
	return arr[TernaryForm((offset)*limit <= len(arr), (offset)*limit, len(arr)):TernaryForm((offset+1)*limit <= len(arr), (offset+1)*limit, len(arr))]
}

func SafeInit[T any]() T {
	var t T
	tt := reflect.TypeOf(t)
	switch tt.Kind() {
	case reflect.Ptr:
		return reflect.New(tt.Elem()).Interface().(T)
	default:
		return t
	}
}

func ToMap(in interface{}, tagName string, omit func(reflect.Value, *string) bool) (map[string]any, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("[ToMap]only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		// 若tag为空，直接跳过即可
		if tagValue := fi.Tag.Get(tagName); tagValue != "" && (omit == nil || !omit(v.Field(i), &tagValue)) {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

// TODO(@hyphen) need complete
func IDL2GORM[I, G any](idl I) (G, error) {
	var g G
	it := reflect.TypeOf(idl)

	switch it.Kind() {
	case reflect.Struct:

	case reflect.Ptr:

	default:
		// 基本类型：直接复制
		//gt := reflect.TypeOf(g)
	}
	return g, nil
}

func MustIDL2GORM[I, G any](idl I) G {
	var g G
	g, err := IDL2GORM[I, G](idl)
	if err != nil {
		return SafeInit[G]()
	}
	return g
}

func SafeAssert[T any](v any) T {
	t, ok := v.(T)
	if !ok {
		return SafeInit[T]()
	}
	return t
}

func WithRetry(ctx context.Context, fc func() error, times ...int) error {
	var err error
	rt := 3
	if len(times) > 0 {
		rt = times[0]
	}
	for i := 0; i < rt; i++ {
		err = fc()
		if err == nil {
			return nil
		}
		logs.CtxInfo(ctx, "[WithRetry]fc exec, times: %d failed: %s", i, err)
		time.Sleep(time.Millisecond * 100)
	}
	return err
}
