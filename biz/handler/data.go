// Package handler
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package handler

import (
	"context"
	"errors"
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

type CallbackData struct {
	Timestamp string           `query:"timestamp"`
	Nonce     string           `query:"nonce"`
	OP        string           `json:"op"`
	Data      client.EntryItem `json:"data"`
}

func HandleData(ctx context.Context, c *app.RequestContext) {
	data := &CallbackData{}
	err := c.Bind(data)
	if err != nil {
		logs.CtxError(ctx, "[HandleData]bind data failed: %s", err)
		c.JSON(http.StatusNotAcceptable, utils.H{
			"message": "invalid param",
		})
		return
	}
	logs.CtxInfo(ctx, "[HandleData]bind data success: %s", utils2.MustMarshal(data))
	go handleData(ctx, data)
	c.JSON(http.StatusOK, utils.H{"message": "success"})
}

func handleData(ctx context.Context, data *CallbackData) {
	logs.CtxInfo(ctx, "[HandleData]start async handle")

	// 尝试设置log id为data_id+timestamp
	id, ok := data.Data.Get("_id")
	if ok {
		ctx = context.WithValue(ctx, logs.CtxKeyLogID, fmt.Sprintf("%s|%d", id, time.Now().Unix()))
	} else {
		logs.CtxError(ctx, "[HandleData]get data id failed")
		//c.JSON(http.StatusNotAcceptable, utils.H{
		//	"message": "get data id failed",
		//})
		return
	}
	// 签名验证 TODO

	set := utils2.NewSetSlice[string, *dao.DesignRecord]()
	// 获取完整 BOM 记录
	boms := map[string]*utils2.Pair[client.EntryItem, string]{}
	for _, design := range config.Config().WidgetDesigns {
		// 对于每一个子控件，获取其三元组，任何一个获取失败，均认为此条记录不存在
		children, ok := data.Data.GetChildren(design.ItemChildren)
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
			ed.Year, ok = data.Data.Get(config.Config().WidgetDesignCommon.ItemYear)
			if !ok {
				logs.CtxInfo(ctx, "[HandleData]get year failed: %s", design.ItemYear)
				continue
			}
			ed.BaseID = id
			if ed.Code == "" && ed.Color == "" && (ed.Type == "" || ed.Type == "面料") {
				continue
			}
			set.Upsert(ed.String(), &ed)
			boms[ed.BaseUUID] = utils2.MakePair(child, design.ItemChildren)
		}
	}
	// 代表当前的三元组
	drs := set.GetSlice()
	logs.CtxInfo(ctx, "[HandleData]get eds: %s", utils2.MustMarshal(drs))

	// 查询记录的三元组
	records, err := dao.QueryRecordData(ctx, id)
	if err != nil {
		// 视为记录中不存在此数据
		logs.CtxError(ctx, "[HandleData]query record data failed: %s", err)
		//c.JSON(http.StatusBadRequest, utils.H{
		//	"message": err.Error(),
		//	"log_id":  ctx.Value(logs.CtxKeyLogID),
		//})
		return
	}
	if len(records) == 0 {
		h := &handler{Inc: drs, Data: data.Data, BOMs: boms}
		err := h.handleInc(ctx)
		if err != nil {
			logs.CtxError(ctx, "[HandleData]handle inc data failed: %s", err)
			//c.JSON(http.StatusBadRequest, utils.H{
			//	"message": err.Error(),
			//	"log_id":  ctx.Value(logs.CtxKeyLogID),
			//})
			return
		}
		logs.CtxInfo(ctx, "[HandleData]Handle Success")
		//c.JSON(200, utils.H{"message": "success"})
		return
	}

	// 代表记录中的三元组
	rset := utils2.NewSetSlice[string, *dao.DesignRecord]()
	for _, record := range records {
		dr, ok := dao.DesignRecord{}, false
		dr.Type, ok = record.Get(config.Config().WidgetRecord.ItemType)
		if !ok {
			logs.CtxInfo(ctx, "[HandleData]2get type failed: %s", config.Config().WidgetRecord.ItemType)
			continue
		}
		dr.Code, ok = record.Get(config.Config().WidgetRecord.ItemCode)
		if !ok {
			logs.CtxInfo(ctx, "[HandleData]2get code failed: %s", config.Config().WidgetRecord.ItemCode)
			continue
		}
		dr.Color, ok = record.Get(config.Config().WidgetRecord.ItemColor)
		if !ok {
			logs.CtxInfo(ctx, "[HandleData]2get color failed: %s", config.Config().WidgetRecord.ItemColor)
			continue
		}
		dr.BaseID, ok = record.Get(config.Config().WidgetRecord.ItemBaseID)
		if !ok {
			logs.CtxInfo(ctx, "[HandleData]2get base id failed: %s", config.Config().WidgetRecord.ItemBaseID)
			continue
		}
		dr.BaseUUID, ok = record.Get(config.Config().WidgetRecord.ItemBaseUUID)
		if !ok {
			logs.CtxInfo(ctx, "[HandleData]2get uuid failed: %s", config.Config().WidgetRecord.ItemBaseUUID)
			continue
		}
		dr.Year, ok = record.Get(config.Config().WidgetRecord.ItemYear)
		if !ok {
			logs.CtxInfo(ctx, "[HandleData]2get year failed: %s", config.Config().WidgetRecord.ItemYear)
			continue
		}
		dr.ID, ok = record.Get("_id")
		if !ok {
			logs.CtxInfo(ctx, "[HandleData]2get id failed")
			continue
		}
		if dr.Code == "" && dr.Color == "" && dr.Type == "" {
			continue
		}
		rset.Upsert(dr.BaseUUID, &dr)
	}
	rdrs := rset.GetSlice()
	// diff 新旧三元组，存在三种可能：不变，新增，删除
	// 如果记录中没有此_id，在record表中新增全量eds，并将eds全部设置为新增
	hdl := &handler{
		Old:  rdrs,
		New:  drs,
		Data: data.Data,
		BOMs: boms,
	}

	hdl.diff(ctx)

	if len(hdl.Dec) != 0 {
		err = hdl.handleDec(ctx)
		if err != nil {
			logs.CtxError(ctx, "[HandleData]handle dec failed: %s", err)
			//c.JSON(http.StatusBadRequest, utils.H{
			//	"message": err.Error(),
			//	"log_id":  ctx.Value(logs.CtxKeyLogID),
			//})
			return
		}
	}

	if len(hdl.Inc) != 0 {
		err = hdl.handleInc(ctx)
		if err != nil {
			logs.CtxError(ctx, "[HandleData]handle inc failed: %s", err)
			//c.JSON(http.StatusBadRequest, utils.H{
			//	"message": err.Error(),
			//	"log_id":  ctx.Value(logs.CtxKeyLogID),
			//})
			return
		}
	}

	if len(hdl.Both) != 0 {
		err = hdl.handleBoth(ctx)
		if err != nil {
			logs.CtxError(ctx, "[HandleData]handle both failed: %s", err)
			//c.JSON(http.StatusBadRequest, utils.H{
			//	"message": err.Error(),
			//	"log_id":  ctx.Value(logs.CtxKeyLogID),
			//})
			return
		}
	}

	//
	logs.CtxInfo(ctx, "[HandleData]Handle Success")
	//c.JSON(200, utils.H{"message": "success"})
}

type handler struct {
	Old, New, Inc, Dec, Both []*dao.DesignRecord
	Data                     client.EntryItem
	BOMs                     map[string]*utils2.Pair[client.EntryItem, string]
}

func (h *handler) diff(ctx context.Context) {
	logs.CtxInfo(ctx, "[diff]start to diff, old: %s; new: %s", utils2.MustMarshal(h.Old), utils2.MustMarshal(h.New))
	var both, inc, dec []*dao.DesignRecord
	oldMap, newMap := map[string]*dao.DesignRecord{}, map[string]*dao.DesignRecord{}
	for _, record := range h.Old {
		oldMap[record.BaseUUID] = record
	}
	for _, record := range h.New {
		newMap[record.BaseUUID] = record
	}
	// 对于new，如果存在于old，则为both，否则为inc
	for _, record := range h.New {
		_, ok := oldMap[record.BaseUUID]
		if !ok {
			inc = append(inc, record)
		} else {
			both = append(both, record)
		}
	}

	// 对与old，如果不存在于new，则为dec
	for _, record := range h.Old {
		_, ok := newMap[record.BaseUUID]
		if !ok {
			dec = append(dec, record)
		}
	}
	logs.CtxInfo(ctx, "[diff]finish to diff, both: %s; inc: %s, dec: %s", utils2.MustMarshal(both), utils2.MustMarshal(inc), utils2.MustMarshal(dec))
	h.Both, h.Inc, h.Dec = both, inc, dec
	return
}

// record表不动，task表: 存在：更新状态为正常，不存在：新增
// detail表：数据全量覆盖
func (h *handler) handleBoth(ctx context.Context) error {
	return h.handleBase(ctx, h.Both)
}

// record表新增，task表: 存在：更新状态为正常，不存在：新增
// detail表：新增
func (h *handler) handleInc(ctx context.Context) error {
	err := dao.CreateRecordDatas(ctx, h.Inc)
	if err != nil {
		logs.CtxError(ctx, "[handleInc]crate record failed: %s", err)
		return err
	}

	return h.handleBase(ctx, h.Inc)
}

func (h *handler) handleBase(ctx context.Context, drs []*dao.DesignRecord) error {
	txID := fmt.Sprintf("%s-%d", ctx.Value(logs.CtxKeyLogID), time.Now().Unix())
	entries := []client.EntryItem{}
	buIDs := []string{}
	for _, record := range drs {
		uuid := record.BaseUUID
		logs.CtxInfo(ctx, "[handleBase]ready to handle: %s", uuid)
		detail, err := dao.GetDetailByUUID(ctx, uuid)
		if err != nil {
			logs.CtxError(ctx, "[handleBase]get detail failed: %s", err)
			return err
		}
		if detail == nil {
			// 新增detail
			logs.CtxInfo(ctx, "[handleBase]detail is nil, will create: %s", uuid)
			p, ok := h.BOMs[uuid]
			if !ok {
				logs.CtxInfo(ctx, "[handleBase]not in BOMs, jump")
				continue
			}
			entry := dao.MakeDetailEntry(ctx, h.Data, p.First, p.Second, txID)
			if entry == nil {
				logs.CtxInfo(ctx, "[handleBase]make entry return nil")
				continue
			}
			entry[config.Config().WidgetDetail.RWZT] = client.NewEntryValue("正常")
			logs.CtxInfo(ctx, "[handleBase]add entry: %s", utils2.MustMarshal(entry))
			entries = append(entries, entry)
			continue
		}
		id, ok := detail.Get("_id")
		if !ok {
			logs.CtxInfo(ctx, "[handleBase]get id failed: %s", utils2.MustMarshal(detail))
			continue
		}
		logs.CtxInfo(ctx, "[handleBase]ready to append id: %s", id)
		buIDs = append(buIDs, id)
	}
	logs.CtxInfo(ctx, "[handleBase]ready to batch update status")
	err := dao.BatchUpdateDetailStatus(ctx, buIDs, "正常")
	if err != nil {
		logs.CtxError(ctx, "[handleBase]batch update status failed: %s", err)
		return err
	}
	logs.CtxInfo(ctx, "[handleBase]entries len: %d", len(entries))
	if len(entries) != 0 {
		err := dao.CreateDetailDatas(ctx, entries, txID)
		if err != nil {
			logs.CtxError(ctx, "[handleBase]create detail failed: %s", err)
			return err
		}
	}

	for _, dr := range drs {
		// 查询task表中是否存在
		task, err := dao.QueryTask(ctx, dr.ToEntryData())
		if err != nil {
			logs.CtxError(ctx, "[handleBase]query task failed: %s", err)
			return err
		}
		if task == nil {
			logs.CtxInfo(ctx, "[handleBase]task not exist")
			p, ok := h.BOMs[dr.BaseUUID]
			if !ok {
				logs.CtxInfo(ctx, "[handleBase]not in BOMs, jump")
				continue
			}
			err = dao.CreateTaskItem(ctx, dr.ToEntryData(), h.Data, p.First, p.Second)
			if err != nil {
				logs.CtxError(ctx, "[handleBase]create task failed: %s", err)
				return err
			}
		} else {
			logs.CtxInfo(ctx, "[handleBase]task exist")
			id, ok := task.Get("_id")
			if !ok {
				logs.CtxError(ctx, "[handleBase]get task item id failed")
				return errors.New("get task item id failed")
			}
			logs.CtxInfo(ctx, "[handleBase]ready to update task item: %s", id)
			err = dao.UpdateTaskItem(ctx, id, dao.TaskStatusNormal)
			if err != nil {
				logs.CtxError(ctx, "[handleBase]failed to update task item: %s", err)
				return err
			}
		}
	}
	return nil
}

// record表删除，task表: 存在：更新状态为已取消，不存在：无变更
// detail表：更新UUID对应的数据的状态。
func (h *handler) handleDec(ctx context.Context) error {
	ids := []string{}
	uuids := []string{}
	for _, dr := range h.Dec {
		ids = append(ids, dr.ID)
		uuids = append(uuids, dr.BaseUUID)
	}
	err := dao.DeleteRecordDatas(ctx, ids)
	if err != nil {
		logs.CtxError(ctx, "[handleDec]delete record failed: %s", err)
		return err
	}

	logs.CtxInfo(ctx, "[handleDec]ready to handle details: %v", uuids)
	for _, uuid := range uuids {
		// 查询是否在Detail表中
		logs.CtxInfo(ctx, "[handleDec]ready to handle detail: %s", uuid)
		detail, err := dao.GetDetailByUUID(ctx, uuid)
		if err != nil {
			logs.CtxError(ctx, "[handleDec]get detail for %s failed: %s", uuid, err)
			return err
		}
		if detail == nil {
			logs.CtxInfo(ctx, "[handleDec]get nil detail for %s", uuid)
			continue
		}
		id, ok := detail.Get("_id")
		if !ok {
			logs.CtxInfo(ctx, "[handleDec]get detail id failed: %s", utils2.MustMarshal(detail))
			continue
		}
		err = dao.UpdateDetailStatus(ctx, id, "已取消")
		if err != nil {
			logs.CtxError(ctx, "[handleDec]update detail status failed: %s, %s", id, uuid)
			return err
		}
		logs.CtxInfo(ctx, "[handleDec]handle detail for: %s success", uuid)
	}

	for _, dr := range h.Dec {
		task, err := dao.QueryTask(ctx, dr.ToEntryData())
		if err != nil {
			logs.CtxError(ctx, "[handleDec]query task failed: %s", err)
			return err
		}
		if task != nil {
			logs.CtxInfo(ctx, "[handleDec]design not exist, task exist")
			id, ok := task.Get("_id")
			if !ok {
				logs.CtxError(ctx, "[handleDec]get task item id failed")
				return errors.New("get task item id failed")
			}
			logs.CtxInfo(ctx, "[handleDec]ready to update task item: %s", id)
			//当前进度不为空时,如果版单BOM表的[企划季度-编码-名称-颜色]全删掉了,开发状态变更为【已取消】
			//当前进度为空时,如果版单BOM表单[企划季度-编码-名名称-颜色]全删掉了,则删除打色任务的该条数据
			pace, ok := task.GetAny("_widget_1696865063325")
			if !ok {
				logs.CtxError(ctx, "[handleDec]get task item pace failed")
				return errors.New("get task item pace failed")
			}
			list, ok := pace.([]any)
			if !ok {
				logs.CtxError(ctx, "[handleDec]type assert list failed")
				return errors.New("type assert list failed")
			}
			if len(list) != 0 {
				err = dao.UpdateTaskItem(ctx, id, dao.TaskStatusCanceled)
				if err != nil {
					logs.CtxError(ctx, "[handleDec]failed to update task item: %s", err)
					return err
				}
			} else {
				err = dao.DeleteTasks(ctx, []string{id})
				if err != nil {
					logs.CtxError(ctx, "[handleDec]failed to delete task item: %s", err)
					return err
				}
			}
		}
	}
	return nil
}

func test(ctx context.Context, data *CallbackData) {
	set := utils2.NewSetSlice[string, *dao.DesignRecord]()
	// 获取完整 BOM 记录
	//boms := map[string]*utils2.Pair[client.EntryItem, string]{}
	for _, design := range config.Config().WidgetDesigns {
		// 对于每一个子控件，获取其三元组，任何一个获取失败，均认为此条记录不存在
		children, ok := data.Data.GetChildren(design.ItemChildren)
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
			ed.Year, ok = data.Data.Get(config.Config().WidgetDesignCommon.ItemYear)
			if !ok {
				logs.CtxInfo(ctx, "[HandleData]get year failed: %s", design.ItemYear)
				continue
			}
			ed.BaseID = ""
			if ed.Code == "" && ed.Color == "" && (ed.Type == "" || ed.Type == "面料") {
				continue
			}
			set.Upsert(ed.String(), &ed)
			//boms[ed.BaseUUID] = utils2.MakePair(child, design.ItemChildren)
		}
	}
	// 代表当前的三元组
	drs := set.GetSlice()
	fmt.Printf("[HandleData]get eds: %s\n", utils2.MustMarshal(drs))
}
