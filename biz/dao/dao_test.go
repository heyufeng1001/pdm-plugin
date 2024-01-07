// Package dao
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/5
package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"pdm-plugin.github.com/config"
	"pdm-plugin.github.com/logs"
)

func TestQueryDesign(t *testing.T) {
	config.Init("../../conf/conf.yaml")
	ok, err := QueryDesignExisted(context.Background(), &EntryData{
		Type:  "面料",
		Code:  "M241S017",
		Color: "茹伊粉紫",
	})
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(ok)
}

func TestQueryFull(t *testing.T) {
	config.Init("../../conf/conf.yaml")
	items, err := QueryDesignFullData(context.Background())
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(len(items))
}

func TestQueryTask(t *testing.T) {
	config.Init("../../conf/conf.yaml")
	item, err := QueryTask(context.Background(), &EntryData{
		Type:  "面料",
		Code:  "M250Z021",
		Color: "朗格伦绿",
		Year:  "moodytiger-2025年春夏企划案",
	})
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(item.Get("_id"))
}

func TestCreateTask(t *testing.T) {
	config.Init("../../conf/conf.yaml")
	err := CreateTaskItem(context.Background(), &EntryData{
		Type:   "面料",
		Code:   "12345",
		Color:  "waibiwaibi",
		Status: "正常",
	}, nil, nil, "")
	if err != nil {
		t.Log(err)
		return
	}
}

func TestQueryDetail(t *testing.T) {
	logs.Init()
	config.Init("../../conf/conf.yaml")
	ok, err := QueryDetail(context.Background(), &EntryData{
		Code:  "MWV3241",
		Color: "",
		Year:  "moodytiger-2024年秋冬企划案",
		WLMC:  "20D消光尼丝纺（底压光sasa",
	})
	t.Log(ok)
	t.Log(err)
	time.Sleep(time.Second)
}

func TestUpdateTask(t *testing.T) {
	config.Init("../../conf/conf.yaml")
	err := UpdateTaskStatus(context.Background(), "6523e0e49609ca273f0a142b", TaskStatusCanceled)
	if err != nil {
		t.Log(err)
		return
	}
}

func TestGetItem(t *testing.T) {
	config.Init("../../conf/conf.yaml")
	item, err := GetSingleDesignItem(context.Background(), "651e5642311fb9c0011fac24")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(item)
}

func TestQb(t *testing.T) {
	config.Init("../../conf/conf.yaml")

	_, _, _, _, _, _, gy, _, _ := queryBase(context.Background(), "面料", "M247S945")
	fmt.Println(gy)
	fmt.Println(gy == nil)
}
