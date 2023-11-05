// Package config
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/4
package config

import (
	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"pdm-plugin.github.com/utils"
)

type AppConfig struct {
	JianDaoHost string `yaml:"JianDaoHost"`
	AppSecret   string `yaml:"AppSecret"`
	AppID       string `yaml:"AppID"`

	EntryDesignID   string `yaml:"EntryDesignID"`
	EntryRecordID   string `yaml:"EntryRecordID"`
	EntryTaskID     string `yaml:"EntryTaskID"`
	EntryDetailID   string `yaml:"EntryDetailID"`
	EntryFabric     string `yaml:"EntryFabric"`
	EntrySubsidiary string `yaml:"EntrySubsidiary"`
	EntryIfColor    string `yaml:"EntryIfColor"`

	URLQueryData       string `yaml:"URLQueryData"`
	URLGetData         string `yaml:"URLGetData"`
	URLCreateData      string `yaml:"URLCreateData"`
	URLUpdateData      string `yaml:"URLUpdateData"`
	URLBatchUpdateData string `yaml:"URLBatchUpdateData"`
	URLBatchCreateData string `yaml:"URLBatchCreateData"`
	URLBatchDeleteData string `yaml:"URLBatchDeleteData"`
	URLGetUploadToken  string `yaml:"URLGetUploadToken"`

	WidgetTask         WidgetConfig             `yaml:"WidgetTask"`
	WidgetDesigns      []WidgetConfig           `yaml:"WidgetDesigns"`
	WidgetDesignCommon WidgetDesignCommonConfig `yaml:"WidgetDesignCommon"`
	WidgetRecord       WidgetConfig             `yaml:"WidgetRecord"`
	WidgetFabric       WidgetConfig             `yaml:"WidgetFabric"`
	WidgetSubsidiary   WidgetConfig             `yaml:"WidgetSubsidiary"`
	WidgetDesignBOMs   []WidgetDesignBOMConfig  `yaml:"WidgetDesignBOMs"`
	WidgetDetail       WidgetDetailConfig       `yaml:"WidgetDetail"`
	WidgetTaskOther    WidgetTaskOtherConfig    `yaml:"WidgetTaskOther"`
	WidgetIfColor      WidgetConfig             `yaml:"WidgetIfColor"`
}

type WidgetConfig struct {
	ItemChildren string `yaml:"ItemChildren"`
	ItemType     string `yaml:"ItemType"`
	ItemCode     string `yaml:"ItemCode"`
	ItemColor    string `yaml:"ItemColor"`
	ItemYear     string `yaml:"ItemYear"`

	ItemStatus string `yaml:"ItemStatus"`

	ItemBaseID   string `yaml:"ItemBaseID"`
	ItemBaseUUID string `yaml:"ItemBaseUUID"`

	ItemSupplierCode string `yaml:"ItemSupplierCode"`
	ItemSupplier     string `yaml:"ItemSupplier"`
	ItemBigRound     string `yaml:"ItemBigRound"`
	ItemProofRound   string `yaml:"ItemProofRound"`
	ItemLeast        string `yaml:"ItemLeast"`
	ItemIfUse        string `yaml:"ItemIfUse"`
	ItemIfColor      string `yaml:"ItemIfColor"`
	ItemGY           string `yaml:"ItemGY"`
	ItemCF           string `yaml:"ItemCF"`
	ItemGYSWLBM      string `yaml:"ItemGYSWLBM"`
}

// WidgetDesignCommonConfig
// Fuck低代码
type WidgetDesignCommonConfig struct {
	ItemYear string `yaml:"ItemYear"`
	NF       string `yaml:"NF"`
	RQSJ     string `yaml:"RQSJ"`
	PM       string `yaml:"PM"`
	CPCJ     string `yaml:"CPCJ"`
	PP       string `yaml:"PP"`
	MLXL     string `yaml:"MLXL"`
	DL       string `yaml:"DL"`
	XL       string `yaml:"XL"`
	JGDW     string `yaml:"JGDW"`
	DQZT     string `yaml:"DQZT"`
	KH       string `yaml:"KH"`
	JJ       string `yaml:"JJ"`
	SJLSH    string `yaml:"SJLSH"`
	FZZD     string `yaml:"FZZD"`
	QHJD     string `yaml:"QHJD"`
	SJS      string `yaml:"SJS"`
	SJSXM    string `yaml:"SJSXM"`
	SPULSH   string `yaml:"SPULSH"`
	SKCXGT   string `yaml:"SKCXGT"`
	ZTXL     string `yaml:"ZTXL"`
	BD       string `yaml:"BD"`
	XB       string `yaml:"XB"`
	ZL       string `yaml:"ZL"`
	JYDJ     string `yaml:"JYDJ"`
	CPMD     string `yaml:"CPMD"`
	GY       string `yaml:"GY"`
}

type WidgetDesignBOMConfig struct {
	BW     string `yaml:"BW"`
	GG     string `yaml:"GG"`
	YL     string `yaml:"YL"`
	BOMID  string `yaml:"BOMID"`
	WLBM   string `yaml:"WLBM"`
	BDBFK  string `yaml:"BDBFK"`
	UUID   string `yaml:"UUID"`
	WLMC   string `yaml:"WLMC"`
	YSYWMC string `yaml:"YSYWMC"`
	SH     string `yaml:"SH"`
	TP     string `yaml:"TP"`
	KZ     string `yaml:"KZ"`
	WLDL   string `yaml:"WLDL"`
	YSZWMC string `yaml:"YSZWMC"`
	SYFK   string `yaml:"SYFK"`
	DW     string `yaml:"DW"`
	ZRR    string `yaml:"ZRR"`
	MLBL   string `yaml:"MLBL"`

	RWZT string `yaml:"RWZT"`
}

type WidgetDetailConfig struct {
	GG       string `yaml:"GG"`
	NF       string `yaml:"NF"`
	RQSJ     string `yaml:"RQSJ"`
	RWZT     string `yaml:"RWZT"`
	YSZWMC   string `yaml:"YSZWMC"`
	SH       string `yaml:"SH"`
	TP       string `yaml:"TP"`
	BW       string `yaml:"BW"`
	KZ       string `yaml:"KZ"`
	FZZD     string `yaml:"FZZD"`
	KH       string `yaml:"KH"`
	PM       string `yaml:"PM"`
	XB       string `yaml:"XB"`
	ZL       string `yaml:"ZL"`
	JYDJ     string `yaml:"JYDJ"`
	DSRWSL   string `yaml:"DSRWSL"`
	SJLSH    string `yaml:"SJLSH"`
	PP       string `yaml:"PP"`
	SJSXM    string `yaml:"SJSXM"`
	MLXL     string `yaml:"MLXL"`
	DW       string `yaml:"DW"`
	CPCJ     string `yaml:"CPCJ"`
	YSYWMC   string `yaml:"YSYWMC"`
	YLDW     string `yaml:"YLDW"`
	SYFK     string `yaml:"SYFK"`
	QHJD     string `yaml:"QHJD"`
	DQZT     string `yaml:"DQZT"`
	SKCXGT   string `yaml:"SKCXGT"`
	UUID     string `yaml:"UUID"`
	SJS      string `yaml:"SJS"`
	BD       string `yaml:"BD"`
	WLBM     string `yaml:"WLBM"`
	YL       string `yaml:"YL"`
	YLZB     string `yaml:"YLZB"`
	CPMD     string `yaml:"CPMD"`
	WLDL     string `yaml:"WLDL"`
	SPULSH   string `yaml:"SPULSH"`
	DL       string `yaml:"DL"`
	JJ       string `yaml:"JJ"`
	ZTXL     string `yaml:"ZTXL"`
	XL       string `yaml:"XL"`
	JGDW     string `yaml:"JGDW"`
	WLMC     string `yaml:"WLMC"`
	BDBFK    string `yaml:"BDBFK"`
	WLXUJXFZ string `yaml:"WLXUJXFZ"`
}

type WidgetTaskOtherConfig struct {
	DW    string `yaml:"DW"`
	WLMC  string `yaml:"WLMC"`
	BD    string `yaml:"BD"`
	SH    string `yaml:"SH"`
	TP    string `yaml:"TP"`
	SYFK  string `yaml:"SYFK"`
	GY    string `yaml:"GY"`
	BDBFK string `yaml:"BDBFK"`
	KZ    string `yaml:"KZ"`
	FZR   string `yaml:"FZR"`
}

var conf *AppConfig

func Config() *AppConfig {
	return conf
}

func Init(path string) {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		panic("read conf failed: " + err.Error())
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic("unmarshal conf failed: " + err.Error())
	}
	logrus.Infof("read AppConfig: %s", utils.MustDo(any(conf), sonic.MarshalString))
}
