// Package dao
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/time/rate"
	"pdm-plugin.github.com/biz/dao/client"
	"pdm-plugin.github.com/biz/dao/filter"
	"pdm-plugin.github.com/config"
	"pdm-plugin.github.com/logs"
	"pdm-plugin.github.com/utils"
)

// QueryDesignExisted
// deprecate
func QueryDesignExisted(ctx context.Context, data *EntryData) (bool, error) {
	logs.CtxInfo(ctx, "[QueryDesignExisted]ready to exec: %s", utils.MustMarshal(data))
	items, err := QueryDesignFullData(ctx)
	if err != nil {
		logs.CtxError(ctx, "[QueryDesignExisted]query full data failed: %s", err)
		return false, err
	}
	logs.CtxInfo(ctx, "[QueryDesignExisted]query full data success: %d", len(items))
	// 存在完全相同的三元组则返回
	for _, item := range items {
		for _, design := range config.Config().WidgetDesigns {
			logs.CtxInfo(ctx, "[QueryDesignExisted]ready to handle: %s, design: %s", item.MGet("_id"), utils.MustMarshal(design))
			t, ok := item.Get(design.ItemType)
			if !ok {
				continue
			}
			code, ok := item.Get(design.ItemCode)
			if !ok {
				continue
			}
			color, ok := item.Get(design.ItemColor)
			if !ok {
				continue
			}
			if t == data.Type && code == data.Code && color == data.Color {
				logs.CtxInfo(ctx, "[QueryDesignExisted]return true")
				return true, nil
			}
		}
	}
	logs.CtxInfo(ctx, "[QueryDesignExisted]return false")
	return false, nil
}

func QueryDesignFullData(ctx context.Context) ([]client.EntryItem, error) {
	logs.CtxInfo(ctx, "[QueryDesignFullData]ready to query full data")
	body := &client.QueryParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryDesignID,
		},
		Limit: 100,
	}
	for _, design := range config.Config().WidgetDesigns {
		body.Fields = append(body.Fields, strings.Split(design.ItemType, ".")[0])
	}

	resp, err := client.NewIDataClient().Query(ctx, body)
	if err != nil {
		logs.CtxError(ctx, "[QueryDesignFullData]query remote failed: %s", err)
		return nil, err
	}
	ret := []client.EntryItem{}
	ret = append(ret, resp.Data...)

	logs.CtxInfo(ctx, "[QueryDesignFullData]first time success")
	for len(resp.Data) == 100 {
		body.DataID, _ = resp.Data[99].Get("_id")
		resp, err = client.NewIDataClient().Query(ctx, body)
		if err != nil {
			logs.CtxError(ctx, "[QueryDesignFullData]query remote failed: %s", err)
			return nil, err
		}
		ret = append(ret, resp.Data...)
		logs.CtxInfo(ctx, "[QueryDesignFullData]other time success")
	}
	logs.CtxInfo(ctx, "[QueryDesignFullData]final success: %d", len(ret))
	return ret, nil
}

func QueryRecordData(ctx context.Context, baseID string) ([]client.EntryItem, error) {
	logs.CtxInfo(ctx, "[QueryRecordData]ready to query record: %s", baseID)
	body := &client.QueryParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryRecordID,
		},
		Filter: filter.NewFilter().WithRel(filter.RelationAnd).
			WithCond(&filter.Cond{
				Field:  config.Config().WidgetRecord.ItemBaseID,
				Type:   "string",
				Method: filter.CondMethodEq,
				Value:  []any{baseID},
			}),
		Limit: 100,
	}
	resp, err := client.NewIDataClient().Query(ctx, body)
	if err != nil {
		logs.CtxError(ctx, "[QueryRecordData]query remote failed: %s", err)
		return nil, err
	}
	ret := []client.EntryItem{}
	ret = append(ret, resp.Data...)

	logs.CtxInfo(ctx, "[QueryRecordData]first time success")

	for len(resp.Data) == 100 {
		body.DataID, _ = resp.Data[99].Get("_id")
		resp, err = client.NewIDataClient().Query(ctx, body)
		if err != nil {
			logs.CtxError(ctx, "[QueryRecordData]query remote failed: %s", err)
			return nil, err
		}
		ret = append(ret, resp.Data...)
		logs.CtxInfo(ctx, "[QueryRecordData]other time success")
	}
	logs.CtxInfo(ctx, "[QueryRecordData]final success: %d", len(ret))
	return ret, nil
}

func CreateRecordDatas(ctx context.Context, datas []*DesignRecord) error {
	logs.CtxInfo(ctx, "[CreateRecordDatas]ready to create record: %s", utils.MustMarshal(datas))
	body := &client.BatchCreateParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryRecordID,
		},
	}

	items := []client.EntryItem{}
	for _, data := range datas {
		item := client.EntryItem(make(map[string]any))
		item[config.Config().WidgetRecord.ItemType] = client.NewEntryValue(data.Type)
		item[config.Config().WidgetRecord.ItemColor] = client.NewEntryValue(data.Color)
		item[config.Config().WidgetRecord.ItemCode] = client.NewEntryValue(data.Code)
		item[config.Config().WidgetRecord.ItemYear] = client.NewEntryValue(data.Year)
		item[config.Config().WidgetRecord.ItemBaseID] = client.NewEntryValue(data.BaseID)
		item[config.Config().WidgetRecord.ItemBaseUUID] = client.NewEntryValue(data.BaseUUID)
		items = append(items, item)
	}
	body.DataList = items
	logs.CtxInfo(ctx, "[CreateRecordDatas]make data success: %s", utils.MustMarshal(items))

	err := client.NewIDataClient().BatchCreate(ctx, body)
	logs.CtxInfo(ctx, "[CreateRecordDatas]return is: %s", err)
	return err
}

func DeleteRecordDatas(ctx context.Context, ids []string) error {
	logs.CtxInfo(ctx, "[DeleteRecordDatas]ready to delete record: %v", ids)
	body := &client.BatchDeleteParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryRecordID,
		},
		DataIDs: ids,
	}
	err := client.NewIDataClient().BatchDelete(ctx, body)
	logs.CtxInfo(ctx, "[DeleteRecordDatas]return is: %s", err)
	return err
}

func GetSingleDesignItem(ctx context.Context, id string) (client.EntryItem, error) {
	resp, err := client.NewIDataClient().Get(ctx, &client.GetParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryDesignID,
		},
		DataID: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func QueryTask(ctx context.Context, data *EntryData) (client.EntryItem, error) {
	logs.CtxInfo(ctx, "[QueryTask]ready to query task: %s", utils.MustMarshal(data))
	body := &client.QueryParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryTaskID,
		},
		Filter: filter.NewFilter().WithRel(filter.RelationAnd).
			WithCond(&filter.Cond{
				Field:  config.Config().WidgetTask.ItemType,
				Type:   "string",
				Method: filter.CondMethodEq,
				Value:  []any{data.Type},
			}).WithCond(&filter.Cond{
			Field:  config.Config().WidgetTask.ItemColor,
			Type:   "string",
			Method: filter.CondMethodEq,
			Value:  []any{data.Color},
		}).WithCond(&filter.Cond{
			Field:  config.Config().WidgetTask.ItemCode,
			Type:   "string",
			Method: filter.CondMethodEq,
			Value:  []any{data.Code},
		}).WithCond(&filter.Cond{
			Field:  config.Config().WidgetTask.ItemYear,
			Type:   "string",
			Method: filter.CondMethodEq,
			Value:  []any{data.Year},
		}),
		Limit: 10,
	}
	v, err := client.NewIDataClient().Query(ctx, body)
	if err != nil {
		logs.CtxError(ctx, "[QueryDesignExisted]query data %v failed: %s", data, err)
		return nil, err
	}
	logs.CtxInfo(ctx, "[QueryTask]get data: %s", utils.MustMarshal(v.Data))
	if len(v.Data) == 0 {
		logs.CtxInfo(ctx, "[QueryTask]return nil")
		return nil, nil
	}
	return v.Data[0], nil
}

func CreateTaskItem(ctx context.Context, data *EntryData, base, bom client.EntryItem, entryID string) error {
	logs.CtxInfo(ctx, "[CreateTaskItem]ready to create task: %s", utils.MustMarshal(data))
	txID := fmt.Sprintf("%s-%d", ctx.Value(logs.CtxKeyLogID), time.Now().Unix())
	body := &client.CreateParam{
		BaseEntry: client.BaseEntry{
			AppID:         config.Config().AppID,
			EntryID:       config.Config().EntryTaskID,
			TransactionID: txID,
		},
	}
	item := client.EntryItem(make(map[string]any))
	item[config.Config().WidgetTask.ItemType] = client.NewEntryValue(data.Type)
	item[config.Config().WidgetTask.ItemColor] = client.NewEntryValue(data.Color)
	item[config.Config().WidgetTask.ItemCode] = client.NewEntryValue(data.Code)
	item[config.Config().WidgetTask.ItemYear] = client.NewEntryValue(data.Year)
	item[config.Config().WidgetTask.ItemStatus] = client.NewEntryValue(TaskStatusNormal)
	// 查询【起订量要求、大货周期、打样周期】并设置
	sp, sc, iu, br, pr, le, gy, cf, gyswlbm := queryBase(ctx, data.Type, data.Code)
	logs.CtxInfo(ctx, "[CreateTaskItem]query base return:sp: %s, sc: %s, iu: %s, br: %d, pr: %d, le: %s, gy: %v, cf: %s, gyswlbm: %s", sp, sc, iu, br, pr, le, gy, cf, gyswlbm)
	if br != 0 {
		item[config.Config().WidgetTask.ItemBigRound] = client.NewEntryValue(br)
	}
	if sp != "" {
		item[config.Config().WidgetTask.ItemSupplier] = client.NewEntryValue(sp)
	}
	//if sc != "" {
	//	item[config.Config().WidgetTask.ItemSupplierCode] = client.NewEntryValue(sc)
	//}
	if iu != "" {
		item[config.Config().WidgetTask.ItemIfUse] = client.NewEntryValue(iu)
	}
	if pr != 0 {
		item[config.Config().WidgetTask.ItemProofRound] = client.NewEntryValue(pr)
	}
	if le != "" {
		item[config.Config().WidgetTask.ItemLeast] = client.NewEntryValue(le)
	}
	if gy != nil {
		item[config.Config().WidgetTask.ItemGY] = client.NewEntryValue(gy)
	}
	if cf != "" {
		item[config.Config().WidgetTask.ItemCF] = client.NewEntryValue(cf)
	}
	//if gyswlbm != "" {
	//	item[config.Config().WidgetTask.ItemGYSWLBM] = client.NewEntryValue(gyswlbm)
	//}

	ic := queryIfColor(ctx, data.Code)
	logs.CtxInfo(ctx, "[CreateTaskItem]query if color return:ic: %s", ic)
	if ic != "" {
		item[config.Config().WidgetTask.ItemIfColor] = client.NewEntryValue(ic)
	}

	var conf *config.WidgetDesignBOMConfig
	for _, bc := range config.Config().WidgetDesignBOMs {
		if bc.BOMID == entryID {
			logs.CtxInfo(ctx, "[MakeDetailEntry]find entry: %s", utils.MustMarshal(bc))
			conf = &bc
			break
		}
	}
	if conf == nil {
		logs.CtxInfo(ctx, "[MakeDetailEntry]no conf find: %s", entryID)
		return nil
	}
	// 处理other字段
	common := config.Config().WidgetDesignCommon
	other := config.Config().WidgetTaskOther

	//SafeSet(ctx, common.GY, base, item, other.GY)
	SafeSet(ctx, common.BD, base, item, other.BD)
	SafeSet(ctx, common.KH, base, item, other.KH)
	SafeSet(ctx, common.SJS, base, item, other.SJS)

	//SafeSet(ctx, conf.TP, bom, item, other.TP)
	SafeSetFile(ctx, conf.TP, bom, item, other.TP, txID, config.Config().EntryTaskID)
	SafeSet(ctx, conf.WLMC, bom, item, other.WLMC)
	SafeSet(ctx, conf.SH, bom, item, other.SH)
	SafeSet(ctx, conf.SYFK, bom, item, other.SYFK)
	SafeSet(ctx, conf.BDBFK, bom, item, other.BDBFK)
	SafeSet(ctx, conf.KZ, bom, item, other.KZ)
	SafeSet(ctx, conf.DW, bom, item, other.DW)
	SafeSet(ctx, conf.ZRR, bom, item, other.FZR)
	SafeSet(ctx, conf.YWLBM, bom, item, other.YWLBM)
	SafeSet(ctx, conf.SFCYWL, bom, item, other.SFCYWL)
	SafeSet(ctx, conf.MLGY, bom, item, other.MLGY)
	SafeSet(ctx, conf.WLJD, bom, item, other.WLJD)
	SafeSet(ctx, conf.TSYQ, bom, item, other.TSYQ)
	SafeSet(ctx, conf.BZ, bom, item, other.BZ)
	SafeSet(ctx, conf.HSDJ, bom, item, other.HSDJ)
	SafeSet(ctx, conf.GYSZH, bom, item, other.GYSZH)
	SafeSet(ctx, conf.SFYYS, bom, item, other.SFYYS)
	SafeSet(ctx, conf.FZZDBDSJ, bom, item, other.FZZDBDSJ)
	SafeSet(ctx, conf.GYSWLBM, bom, item, other.GYSWLBM)

	body.Data = item
	err := client.NewIDataClient().Create(ctx, body)
	logs.CtxInfo(ctx, "[CreateTaskItem]return is: %s", err)
	return err
}

func MakeTaskItem(ctx context.Context, data *EntryData, base, bom client.EntryItem,
	entryID string, isSync bool) client.EntryItem {
	txID := fmt.Sprintf("%s-%d", ctx.Value(logs.CtxKeyLogID), time.Now().Unix())
	item := client.EntryItem(make(map[string]any))
	item[config.Config().WidgetTask.ItemType] = client.NewEntryValue(data.Type)
	item[config.Config().WidgetTask.ItemColor] = client.NewEntryValue(data.Color)
	item[config.Config().WidgetTask.ItemCode] = client.NewEntryValue(data.Code)
	item[config.Config().WidgetTask.ItemYear] = client.NewEntryValue(data.Year)
	if !isSync {
		item[config.Config().WidgetTask.ItemStatus] = client.NewEntryValue(TaskStatusNormal)
	}
	// 查询【起订量要求、大货周期、打样周期】并设置
	sp, sc, iu, br, pr, le, gy, cf, gyswlbm := queryBase(ctx, data.Type, data.Code)
	logs.CtxInfo(ctx, "[CreateTaskItem]query base return:sp: %s, sc: %s, iu: %s, br: %d, pr: %d, le: %s, gy: %v, cf: %s, gyswlbm: %s", sp, sc, iu, br, pr, le, gy, cf, gyswlbm)
	if br != 0 {
		item[config.Config().WidgetTask.ItemBigRound] = client.NewEntryValue(br)
	}
	if sp != "" {
		item[config.Config().WidgetTask.ItemSupplier] = client.NewEntryValue(sp)
	}
	//if sc != "" {
	//	item[config.Config().WidgetTask.ItemSupplierCode] = client.NewEntryValue(sc)
	//}
	if iu != "" {
		item[config.Config().WidgetTask.ItemIfUse] = client.NewEntryValue(iu)
	}
	if pr != 0 {
		item[config.Config().WidgetTask.ItemProofRound] = client.NewEntryValue(pr)
	}
	if le != "" {
		item[config.Config().WidgetTask.ItemLeast] = client.NewEntryValue(le)
	}
	if gy != nil {
		item[config.Config().WidgetTask.ItemGY] = client.NewEntryValue(gy)
	}
	if cf != "" {
		item[config.Config().WidgetTask.ItemCF] = client.NewEntryValue(cf)
	}
	//if gyswlbm != "" {
	//	item[config.Config().WidgetTask.ItemGYSWLBM] = client.NewEntryValue(gyswlbm)
	//}

	ic := queryIfColor(ctx, data.Code)
	logs.CtxInfo(ctx, "[CreateTaskItem]query if color return:ic: %s", ic)
	if ic != "" {
		item[config.Config().WidgetTask.ItemIfColor] = client.NewEntryValue(ic)
	}

	var conf *config.WidgetDesignBOMConfig
	for _, bc := range config.Config().WidgetDesignBOMs {
		if bc.BOMID == entryID {
			logs.CtxInfo(ctx, "[MakeDetailEntry]find entry: %s", utils.MustMarshal(bc))
			conf = &bc
			break
		}
	}
	if conf == nil {
		logs.CtxInfo(ctx, "[MakeDetailEntry]no conf find: %s", entryID)
		return nil
	}
	// 处理other字段
	common := config.Config().WidgetDesignCommon
	other := config.Config().WidgetTaskOther

	//SafeSet(ctx, common.GY, base, item, other.GY)
	SafeSet(ctx, common.BD, base, item, other.BD)
	SafeSet(ctx, common.KH, base, item, other.KH)
	SafeSet(ctx, common.SJS, base, item, other.SJS)

	//SafeSet(ctx, conf.TP, bom, item, other.TP)
	SafeSetFile(ctx, conf.TP, bom, item, other.TP, txID, config.Config().EntryTaskID)
	SafeSet(ctx, conf.WLMC, bom, item, other.WLMC)
	SafeSet(ctx, conf.SH, bom, item, other.SH)
	SafeSet(ctx, conf.SYFK, bom, item, other.SYFK)
	SafeSet(ctx, conf.BDBFK, bom, item, other.BDBFK)
	SafeSet(ctx, conf.KZ, bom, item, other.KZ)
	SafeSet(ctx, conf.DW, bom, item, other.DW)
	SafeSet(ctx, conf.ZRR, bom, item, other.FZR)
	SafeSet(ctx, conf.YWLBM, bom, item, other.YWLBM)
	SafeSet(ctx, conf.SFCYWL, bom, item, other.SFCYWL)
	SafeSet(ctx, conf.MLGY, bom, item, other.MLGY)
	SafeSet(ctx, conf.WLJD, bom, item, other.WLJD)
	SafeSet(ctx, conf.TSYQ, bom, item, other.TSYQ)
	SafeSet(ctx, conf.BZ, bom, item, other.BZ)
	SafeSet(ctx, conf.HSDJ, bom, item, other.HSDJ)
	SafeSet(ctx, conf.GYSZH, bom, item, other.GYSZH)
	SafeSet(ctx, conf.SFYYS, bom, item, other.SFYYS)
	if !isSync {
		SafeSet(ctx, conf.FZZDBDSJ, bom, item, other.FZZDBDSJ)
	}
	SafeSet(ctx, conf.GYSWLBM, bom, item, other.GYSWLBM)

	return item
}

const (
	cateFabric = "面料"
	cateSub    = "辅料"
)

func queryBase(ctx context.Context, cate string, code string) (string, string, string, int, int, string, any, string, string) {
	body := &client.QueryParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: utils.TernaryForm(cate == cateFabric, config.Config().EntryFabric, config.Config().EntrySubsidiary),
		},
		Limit: 100,
	}
	widget := utils.TernaryForm(cate == cateFabric, config.Config().WidgetFabric, config.Config().WidgetSubsidiary)
	body.Filter = filter.NewFilter().WithRel(filter.RelationAnd).
		WithCond(&filter.Cond{
			Field:  widget.ItemCode,
			Type:   "string",
			Method: filter.CondMethodEq,
			Value:  []any{code},
		})
	resp, err := client.NewIDataClient().Query(ctx, body)
	if err != nil {
		return "", "", "", 0, 0, "", nil, "", ""
	}
	if len(resp.Data) == 0 {
		return "", "", "", 0, 0, "", nil, "", ""
	}
	data := resp.Data[0]
	return data.MGet(widget.ItemSupplier), data.MGet(widget.ItemSupplierCode), data.MGet(widget.ItemIfUse),
		int(utils.SafeAssert[float64](data.MGetAny(widget.ItemBigRound))), int(utils.SafeAssert[float64](data.MGetAny(widget.ItemProofRound))),
		data.MGet(widget.ItemLeast), data.MGetAny(widget.ItemGY), data.MGet(widget.ItemCF), data.MGet(widget.ItemGYSWLBM)
}

func queryIfColor(ctx context.Context, code string) string {
	body := &client.QueryParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryIfColor,
		},
		Limit: 100,
	}
	widget := config.Config().WidgetIfColor
	body.Filter = filter.NewFilter().WithRel(filter.RelationAnd).
		WithCond(&filter.Cond{
			Field:  widget.ItemCode,
			Type:   "string",
			Method: filter.CondMethodEq,
			Value:  []any{code},
		})
	resp, err := client.NewIDataClient().Query(ctx, body)
	if err != nil {
		return ""
	}
	if len(resp.Data) == 0 {
		return ""
	}
	data := resp.Data[0]
	return data.MGet(widget.ItemIfColor)
}

func UpdateTaskStatus(ctx context.Context, id string, status TaskStatus) error {
	logs.CtxInfo(ctx, "[UpdateTaskStatus]ready to update task: %s, %s", id, status)
	body := &client.UpdateParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryTaskID,
		},
		DataID: id,
	}
	item := client.EntryItem(make(map[string]any))
	item[config.Config().WidgetTask.ItemStatus] = client.NewEntryValue(status)
	body.Data = item
	err := client.NewIDataClient().Update(ctx, body)
	logs.CtxInfo(ctx, "[UpdateTaskStatus]return is: %s", err)
	return err
}

func UpdateTask(ctx context.Context, id string, item client.EntryItem) error {
	logs.CtxInfo(ctx, "[UpdateTaskStatus]ready to update task: %s, %s", id, utils.MustMarshal(item))
	body := &client.UpdateParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryTaskID,
		},
		DataID: id,
		Data:   item,
	}
	err := client.NewIDataClient().Update(ctx, body)
	logs.CtxInfo(ctx, "[UpdateTaskStatus]return is: %s", err)
	return err
}

func DeleteTasks(ctx context.Context, ids []string) error {
	logs.CtxInfo(ctx, "[DeleteRecordDatas]ready to delete record: %v", ids)
	body := &client.BatchDeleteParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryTaskID,
		},
		DataIDs: ids,
	}
	err := client.NewIDataClient().BatchDelete(ctx, body)
	logs.CtxInfo(ctx, "[DeleteRecordDatas]return is: %s", err)
	return err
}

var gdLimiter = rate.NewLimiter(30, 30)

func GetDetailByUUID(ctx context.Context, uuid string) (client.EntryItem, error) {
	logs.CtxInfo(ctx, "[GetDetailByUUID]start, uuid: %s", uuid)
	gdLimiter.Wait(ctx)
	body := &client.QueryParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryDetailID,
		},
		Filter: filter.NewFilter().WithRel(filter.RelationAnd).
			WithCond(&filter.Cond{
				Field:  config.Config().WidgetDetail.UUID,
				Type:   "string",
				Method: filter.CondMethodEq,
				Value:  []any{uuid},
			}),
		Limit: 100,
	}
	item, err := client.NewIDataClient().Query(ctx, body)
	if err != nil {
		logs.CtxError(ctx, "[GetDetailByUUID]query remote failed: %s", err)
		return nil, err
	}
	logs.CtxInfo(ctx, "[GetDetailByUUID]data len: %d", len(item.Data))
	if len(item.Data) == 0 {
		logs.CtxInfo(ctx, "[GetDetailByUUID]return both nil")
		return nil, nil
	}
	logs.CtxInfo(ctx, "[GetDetailByUUID]finish: %s", utils.MustMarshal(item.Data[0]))
	return item.Data[0], nil
}

func QueryDetail(ctx context.Context, data *EntryData) (bool, error) {
	logs.CtxInfo(ctx, "[QueryDetail]start, uuid: %s", utils.MustMarshal(data))
	gdLimiter.Wait(ctx)
	body := &client.QueryParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryDetailID,
		},
		Filter: filter.NewFilter().WithRel(filter.RelationAnd).
			WithCond(&filter.Cond{
				Field:  config.Config().WidgetDetail.WLBM,
				Type:   "string",
				Method: filter.CondMethodEq,
				Value:  []any{data.Code},
			}).WithCond(&filter.Cond{
			Field:  config.Config().WidgetDetail.QHJD,
			Type:   "string",
			Method: filter.CondMethodEq,
			Value:  []any{data.Year},
		}).WithCond(&filter.Cond{
			Field:  config.Config().WidgetDetail.YSZWMC,
			Type:   "string",
			Method: filter.CondMethodEq,
			Value:  []any{data.Color},
		}).WithCond(&filter.Cond{
			Field:  config.Config().WidgetDetail.WLMC,
			Type:   "string",
			Method: filter.CondMethodEq,
			Value:  []any{data.WLMC},
		}),
		Limit: 10,
	}
	item, err := client.NewIDataClient().Query(ctx, body)
	if err != nil {
		logs.CtxError(ctx, "[GetDetailByUUID]query remote failed: %s", err)
		return false, err
	}
	logs.CtxInfo(ctx, "[GetDetailByUUID]data len: %d", len(item.Data))
	if len(item.Data) == 0 {
		logs.CtxInfo(ctx, "[GetDetailByUUID]return both nil")
		return false, nil
	}
	logs.CtxInfo(ctx, "[GetDetailByUUID]finish: %s", utils.MustMarshal(item.Data[0]))
	return true, nil
}

func CreateDetailDatas(ctx context.Context, datas []client.EntryItem, txID string) error {
	logs.CtxInfo(ctx, "[CreateDetailDatas]start: %s", utils.MustMarshal(datas))
	body := &client.BatchCreateParam{
		BaseEntry: client.BaseEntry{
			AppID:         config.Config().AppID,
			EntryID:       config.Config().EntryDetailID,
			TransactionID: txID,
		},
		DataList: datas,
	}
	err := client.NewIDataClient().BatchCreate(ctx, body)
	logs.CtxInfo(ctx, "[CreateDetailDatas]return is: %s", err)
	return err
}

func UpdateDetailStatus(ctx context.Context, id string, status string) error {
	logs.CtxInfo(ctx, "[UpdateDetailStatus]start, id: %s, status: %s", id, status)
	body := &client.UpdateParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryDetailID,
		},
		DataID: id,
	}
	item := client.EntryItem(make(map[string]any))
	item[config.Config().WidgetDetail.RWZT] = client.NewEntryValue(status)
	body.Data = item
	err := client.NewIDataClient().Update(ctx, body)
	logs.CtxInfo(ctx, "[UpdateDetailStatus]return is: %s", err)
	return err
}

func UpdateDetail(ctx context.Context, id string, item client.EntryItem) error {
	logs.CtxInfo(ctx, "[UpdateDetailStatus]start, id: %s, status: %s", id, utils.MustMarshal(item))
	body := &client.UpdateParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryDetailID,
		},
		DataID: id,
	}
	body.Data = item
	err := client.NewIDataClient().Update(ctx, body)
	logs.CtxInfo(ctx, "[UpdateDetailStatus]return is: %s", err)
	return err
}

func BatchUpdateDetailStatus(ctx context.Context, ids []string, status string) error {
	logs.CtxInfo(ctx, "[BatchUpdateDetailStatus]start, ids: %s, status: %s", ids, status)
	if len(ids) == 0 {
		logs.CtxInfo(ctx, "[BatchUpdateDetailStatus]empty ids")
		return nil
	}
	body := &client.BatchUpdateParam{
		BaseEntry: client.BaseEntry{
			AppID:   config.Config().AppID,
			EntryID: config.Config().EntryDetailID,
		},
		DataIDs: ids,
	}
	item := client.EntryItem(make(map[string]any))
	item[config.Config().WidgetDetail.RWZT] = client.NewEntryValue(status)
	// 这里对其他字段也做增量修改
	body.Data = item
	err := client.NewIDataClient().BatchUpdate(ctx, body)
	logs.CtxInfo(ctx, "[BatchUpdateDetailStatus]return is: %s", err)
	return err
}
