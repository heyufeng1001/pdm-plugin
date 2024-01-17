// Package dao
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package dao

import (
	"context"

	"pdm-plugin.github.com/biz/dao/client"
	"pdm-plugin.github.com/config"
	"pdm-plugin.github.com/logs"
	"pdm-plugin.github.com/utils"
)

type EntryData struct {
	Type   string
	Code   string
	Color  string
	Year   string
	Status string
	WLMC   string
}

func (e *EntryData) String() string {
	return e.Type + e.Code + e.Color + e.Year
}

type TaskStatus string

const (
	TaskStatusNormal   TaskStatus = "正常开发"
	TaskStatusCanceled TaskStatus = "已取消"
)

type DesignRecord struct {
	ID       string
	Type     string
	Code     string
	Color    string
	Year     string
	BaseID   string
	BaseUUID string
}

func (d *DesignRecord) String() string {
	return d.Type + d.Code + d.Color
}

func (d *DesignRecord) ToEntryData() *EntryData {
	return &EntryData{
		Type:  d.Type,
		Code:  d.Code,
		Color: d.Color,
		Year:  d.Year,
	}
}

func MakeDetailEntry(ctx context.Context, base, bom client.EntryItem, entryID, txID string) client.EntryItem {
	logs.CtxInfo(ctx, "[MakeDetailEntry]start, base: %s, bom: %s, entry id: %s", utils.MustMarshal(base), utils.MustMarshal(bom), entryID)
	ret := client.EntryItem(map[string]any{})
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
	common := config.Config().WidgetDesignCommon
	detail := config.Config().WidgetDetail

	SafeSet(ctx, conf.BW, bom, ret, detail.BW)
	SafeSet(ctx, conf.DW, bom, ret, detail.DW)
	SafeSet(ctx, conf.GG, bom, ret, detail.GG)
	SafeSet(ctx, conf.KZ, bom, ret, detail.KZ)
	SafeSet(ctx, conf.BDBFK, bom, ret, detail.BDBFK)
	SafeSet(ctx, conf.SH, bom, ret, detail.SH)
	SafeSet(ctx, conf.SYFK, bom, ret, detail.SYFK)
	SafeSet(ctx, conf.UUID, bom, ret, detail.UUID)
	SafeSet(ctx, conf.WLBM, bom, ret, detail.WLBM)
	SafeSet(ctx, conf.WLDL, bom, ret, detail.WLDL)
	SafeSet(ctx, conf.WLMC, bom, ret, detail.WLMC)
	SafeSet(ctx, conf.YL, bom, ret, detail.YL)
	SafeSet(ctx, conf.YSYWMC, bom, ret, detail.YSYWMC)
	SafeSet(ctx, conf.YSZWMC, bom, ret, detail.YSZWMC)
	SafeSet(ctx, conf.MLBL, bom, ret, detail.YLZB)
	SafeSet(ctx, conf.YWLBM, bom, ret, detail.YWLBM)
	SafeSet(ctx, conf.SFCYWL, bom, ret, detail.SFCYWL)
	SafeSet(ctx, conf.MLGY, bom, ret, detail.MLGY)
	SafeSet(ctx, conf.WLJD, bom, ret, detail.WLJD)
	SafeSet(ctx, conf.TSYQ, bom, ret, detail.TSYQ)
	SafeSet(ctx, conf.BZ, bom, ret, detail.BZ)
	SafeSet(ctx, conf.SH, bom, ret, detail.SH)
	SafeSet(ctx, conf.HSDJ, bom, ret, detail.HSDJ)
	SafeSet(ctx, conf.GYSJC, bom, ret, detail.GYSJC)
	SafeSet(ctx, conf.GYSZH, bom, ret, detail.GYSZH)
	SafeSet(ctx, conf.SFYYS, bom, ret, detail.SFYYS)
	SafeSet(ctx, conf.FZZDBDSJ, bom, ret, detail.FZZDBDSJ)
	SafeSet(ctx, conf.GYSWLBM, bom, ret, detail.GYSWLBM)
	SafeSet(ctx, conf.SZ, bom, ret, detail.SZ)
	SafeSet(ctx, conf.GYSMC, bom, ret, detail.GYSMC)
	SafeSet(ctx, conf.CGZRR, bom, ret, detail.CGZRR)
	SafeSet(ctx, conf.JCLX, bom, ret, detail.JCLX)
	SafeSet(ctx, conf.DSGD, bom, ret, detail.DSGD)
	SafeSetFile(ctx, conf.TP, bom, ret, detail.TP, txID, config.Config().EntryDetailID)

	SafeSet(ctx, common.BD, base, ret, detail.BD)
	SafeSet(ctx, common.CPCJ, base, ret, detail.CPCJ)
	SafeSet(ctx, common.DL, base, ret, detail.DL)
	SafeSet(ctx, common.CPMD, base, ret, detail.CPMD)
	SafeSet(ctx, common.DQZT, base, ret, detail.DQZT)
	SafeSet(ctx, common.FZZD, base, ret, detail.FZZD)
	SafeSet(ctx, common.JGDW, base, ret, detail.JGDW)
	SafeSet(ctx, common.JJ, base, ret, detail.JJ)
	SafeSet(ctx, common.JYDJ, base, ret, detail.JYDJ)
	SafeSet(ctx, common.KH, base, ret, detail.KH)
	SafeSet(ctx, common.MLXL, base, ret, detail.MLXL)
	SafeSet(ctx, common.NF, base, ret, detail.NF)
	SafeSet(ctx, common.PM, base, ret, detail.PM)
	SafeSet(ctx, common.PP, base, ret, detail.PP)
	SafeSet(ctx, common.QHJD, base, ret, detail.QHJD)
	SafeSet(ctx, common.RQSJ, base, ret, detail.RQSJ)
	SafeSet(ctx, common.SJLSH, base, ret, detail.SJLSH)
	SafeSet(ctx, common.SJS, base, ret, detail.SJS)
	SafeSet(ctx, common.SJSXM, base, ret, detail.SJSXM)
	SafeSet(ctx, common.SPULSH, base, ret, detail.SPULSH)
	SafeSet(ctx, common.XB, base, ret, detail.XB)
	SafeSet(ctx, common.XL, base, ret, detail.XL)
	SafeSet(ctx, common.ZL, base, ret, detail.ZL)
	SafeSet(ctx, common.ZTXL, base, ret, detail.ZTXL)
	SafeSetFile(ctx, common.SKCXGT, base, ret, detail.SKCXGT, txID, config.Config().EntryDetailID)

	ret[detail.DSRWSL] = client.NewEntryValue(1)

	logs.CtxInfo(ctx, "[MakeDetailEntry]make entry success: %s", utils.MustMarshal(ret))
	return ret
}

func SafeSet(ctx context.Context, key string, from, to client.EntryItem, keyTo string) {
	logs.CtxInfo(ctx, "[SafeSet]start: key: %s, from: %s, to: %s, keyTo: %s", key, utils.MustMarshal(from), utils.MustMarshal(to), keyTo)
	kk, ok := from.GetAny(key)
	if !ok {
		logs.CtxInfo(ctx, "[SafeSet]get any return false")
		return
	}
	if keyTo == config.Config().WidgetTaskOther.FZR || keyTo == config.Config().WidgetDetail.SJS || keyTo == config.Config().WidgetTaskOther.SJS ||
		keyTo == config.Config().WidgetTaskOther.GYSZH || keyTo == config.Config().WidgetDetail.GYSZH || keyTo == config.Config().WidgetDetail.CGZRR {
		logs.CtxInfo(ctx, "[SafeSet]into key fzr, kk is: %s", utils.MustMarshal(kk))
		people, ok := kk.(map[string]any)
		if !ok {
			logs.CtxInfo(ctx, "[SafeSet]people type assert failed")
			return
		}
		name, ok := people["username"].(string)
		if !ok {
			logs.CtxInfo(ctx, "[SafeSet]name type assert failed")
			return
		}
		logs.CtxInfo(ctx, "[SafeSet]%s success: %v", key, name)
		to[keyTo] = client.NewEntryValue(name)
		return
	}

	logs.CtxInfo(ctx, "[SafeSet]success: %v", kk)
	to[keyTo] = client.NewEntryValue(kk)
}

func SafeSetFile(ctx context.Context, key string, from, to client.EntryItem, keyTo, txID, entryID string) {
	logs.CtxInfo(ctx, "[SafeSetFile]ready to set file: key: %s, from: %s, to: %s, keyTo: %s, tx_id: %s", key, utils.MustMarshal(from), utils.MustMarshal(to), keyTo, txID)
	kk, ok := from.GetAny(key)
	if !ok {
		logs.CtxError(ctx, "[SafeSetFile]get key failed")
		return
	}
	files, ok := kk.([]any)
	if !ok {
		logs.CtxError(ctx, "[SafeSetFile]kk type assert failed")
		return
	}
	fs := []*utils.Pair[string, string]{}
	for idx, ff := range files {
		file, ok := ff.(map[string]any)
		if !ok {
			logs.CtxError(ctx, "[SafeSetFile]ff type assert failed")
			return
		}
		logs.CtxInfo(ctx, "[SafeSetFile]ready to handle: %d, %s", idx, utils.MustMarshal(file))
		u, ok := file["url"]
		if !ok {
			logs.CtxError(ctx, "[SafeSetFile]get url failed, idx: %d", idx)
			return
		}
		url, ok := u.(string)
		if !ok {
			logs.CtxError(ctx, "[SafeSetFile]url type assert failed, idx: %d", idx)
			return
		}
		n, ok := file["name"]
		if !ok {
			logs.CtxError(ctx, "[SafeSetFile]get name failed, idx: %d", idx)
			return
		}
		name, ok := n.(string)
		if !ok {
			logs.CtxError(ctx, "[SafeSetFile]name type assert failed, idx: %d", idx)
			return
		}
		fs = append(fs, utils.MakePair(url, name))
	}

	if len(fs) == 0 {
		logs.CtxInfo(ctx, "[SafeSetFile]fs is nil")
		return
	}

	keys, err := UploadImage(ctx, fs, txID, entryID)
	if err != nil {
		logs.CtxError(ctx, "[SafeSetFile]upload failed: %s", err)
		return
	}
	vs := []any{}
	for _, s := range keys {
		vs = append(vs, s)
	}
	logs.CtxInfo(ctx, "[SafeSetFile]%s success: %s", key, utils.MustMarshal(vs))
	to[keyTo] = client.NewEntryValue(vs)
}
