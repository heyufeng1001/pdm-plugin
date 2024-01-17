// Package handler
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/12/29
package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"pdm-plugin.github.com/biz/dao"
	"pdm-plugin.github.com/biz/dao/client"
	"pdm-plugin.github.com/config"
	"pdm-plugin.github.com/logs"
	utils2 "pdm-plugin.github.com/utils"
)

func HandleSync(ctx context.Context, c *app.RequestContext) {
	//go syncFullData(ctx)
	c.JSON(http.StatusOK, utils.H{"message": "success"})
}

func syncFullData(ctx context.Context) {
	// 获取版单设计的全量数据
	logs.CtxInfo(ctx, "[QueryDesignFullData]ready to query full data")
	body := &client.QueryParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryDesignID,
		},
		Limit: 100,
	}
	//for _, design := range config.Config().WidgetDesigns {
	//	body.Fields = append(body.Fields, strings.Split(design.ItemType, ".")[0])
	//}

	resp, err := client.NewIDataClient().Query(ctx, body)
	if err != nil {
		logs.CtxError(ctx, "[QueryDesignFullData]query remote failed: %s", err)
		return
	}

	logs.CtxInfo(ctx, "[QueryDesignFullData]first time success")
	for {
		syncData100(ctx, resp.Data)

		if len(resp.Data) < 100 {
			break
		}
		body.DataID, _ = resp.Data[99].Get("_id")
		resp, err = client.NewIDataClient().Query(ctx, body)
		if err != nil {
			logs.CtxError(ctx, "[QueryDesignFullData]query remote failed: %s", err)
			return
		}
		logs.CtxInfo(ctx, "[QueryDesignFullData]other time success")
	}
	return
}

func syncData100(ctx context.Context, items []client.EntryItem) {
	// 查找 uuid 对应的detail
	for _, item := range items {
		syncSingleData(ctx, item)
	}
}

func syncSingleData(ctx context.Context, item client.EntryItem) {
	id, ok := item.Get("_id")
	if !ok {
		logs.CtxError(ctx, "[HandleData]get data id failed")
		//c.JSON(http.StatusNotAcceptable, utils.H{
		//	"message": "get data id failed",
		//})
		return
	}
	logs.CtxInfo(ctx, "[syncSingleData]ready to sync: %s", id)
	// 获取boms
	set := utils2.NewSetSlice[string, *dao.DesignRecord]()
	// 获取完整 BOM 记录
	boms := map[string]*utils2.Pair[client.EntryItem, string]{}
	for _, design := range config.Config().WidgetDesigns {
		// 对于每一个子控件，获取其三元组，任何一个获取失败，均认为此条记录不存在
		children, ok := item.GetChildren(design.ItemChildren)
		if !ok {
			logs.CtxInfo(ctx, "[HandleData]get children failed: %s", design.ItemChildren)
			continue
		}
		for _, child := range children {
			ed, ok := dao.DesignRecord{}, false
			ed.Type, ok = child.Get(design.ItemType)
			if !ok {
				logs.CtxInfo(ctx, "[HandleData]get type failed: %s", design.ItemType)
				continue
			}
			ed.Code, ok = child.Get(design.ItemCode)
			if !ok {
				logs.CtxInfo(ctx, "[HandleData]get code failed: %s", design.ItemCode)
				continue
			}
			ed.Color, ok = child.Get(design.ItemColor)
			if !ok {
				logs.CtxInfo(ctx, "[HandleData]get color failed: %s", design.ItemColor)
				continue
			}
			ed.BaseUUID, ok = child.Get(design.ItemBaseUUID)
			if !ok {
				logs.CtxInfo(ctx, "[HandleData]get uuid failed: %s", design.ItemBaseUUID)
				continue
			}
			ed.Year, ok = item.Get(config.Config().WidgetDesignCommon.ItemYear)
			if !ok {
				logs.CtxInfo(ctx, "[HandleData]get year failed: %s", design.ItemYear)
				continue
			}
			ed.BaseID = id
			if ed.Code == "" && ed.Color == "" && (ed.Type == "" || ed.Type == "面料") {
				continue
			}
			set.Upsert(ed.BaseUUID, &ed)
			boms[ed.BaseUUID] = utils2.MakePair(child, design.ItemChildren)
		}
	}
	// 代表当前的三元组
	drs := set.GetSlice()
	logs.CtxInfo(ctx, "[HandleData]get eds: %s", utils2.MustMarshal(drs))
	// 对于每个boms，执行同步逻辑
	for uuid, bom := range boms {
		txID := fmt.Sprintf("%s-%d", id, time.Now().Unix())
		detail, err := dao.GetDetailByUUID(ctx, uuid)
		if err != nil {
			logs.CtxError(ctx, "[]get uuid %s failed: %s", uuid, err)
			continue
		}
		if detail == nil {
			logs.CtxInfo(ctx, "[]uuid %s is nil", uuid)
			continue
		}
		bid, ok := bom.First.Get("_id")
		if !ok {
			logs.CtxInfo(ctx, "no id for %s", uuid)
			continue
		}
		e := dao.MakeDetailEntry(ctx, item, bom.First, bom.Second, txID)
		err = dao.UpdateDetail(ctx, bid, e, txID)
		if err != nil {
			logs.CtxError(ctx, "[]update uuid %s failed: %s", uuid, err)
			continue
		}
	}

	for _, dr := range drs {
		task, err := dao.QueryTask(ctx, dr.ToEntryData())
		if err != nil {
			logs.CtxError(ctx, "[handleBase]query task failed: %s", err)
			continue
		}
		if task == nil {
			logs.CtxInfo(ctx, "no task")
			continue
		}
		tid, ok := task.Get("_id")
		if !ok {
			logs.CtxInfo(ctx, "[]no tid is task")
			continue
		}
		logs.CtxInfo(ctx, "[handleBase]task exist")
		bom := boms[dr.BaseUUID]
		nt, txID := dao.MakeTaskItem(ctx, dr.ToEntryData(), item, bom.First, bom.Second, true)
		if nt == nil {
			logs.CtxInfo(ctx, "[]make task failed: %s", dr.BaseUUID)
		}
		err = dao.UpdateTask(ctx, tid, nt, txID)
		if err != nil {
			logs.CtxInfo(ctx, "[]update task failed: %s", err)
			continue
		}
	}

}
