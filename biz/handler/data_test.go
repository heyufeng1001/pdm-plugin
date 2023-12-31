// Package handler
// Author: hyphen
// Copyright 2023 hyphen. All rights reserved.
// Create-time: 2023/10/18
package handler

import (
	"context"
	"encoding/json"
	"testing"

	"pdm-plugin.github.com/config"
)

func Test111(t *testing.T) {
	config.Init("../../conf/conf.yaml")
	d := &CallbackData{}
	json.Unmarshal([]byte(s), &d)
	test(context.Background(), d)
}

var s = `{
  "Timestamp": "1697618646",
  "Nonce": "e2f839",
  "op": "data_update",
  "data": {
    "_id": "651fc1b9550d16a33b372bb7",
    "_widget_1687311922144": "2024",
    "_widget_1689772437724": "00135",
    "_widget_1689820883385": {
      "_id": "64ffe045479d570008d00f8d",
      "name": "Liena Li",
      "status": 1,
      "type": 0,
      "username": "fbge8ccc"
    },
    "_widget_1689821736839": {
      "id": "2cd770b864feeb7af07732fd"
    },
    "_widget_1689821813925": "服装",
    "_widget_1689821813926": "秋冬",
    "_widget_1689821813927": "法网圣殿",
    "_widget_1689821813928": "夏2波",
    "_widget_1689821813929": "中性",
    "_widget_1689821813930": "连体",
    "_widget_1689821813931": "连体裤",
    "_widget_1689821813932": "进阶款",
    "_widget_1689821813933": "中阶价",
    "_widget_1689821813935": 498,
    "_widget_1689821813937": "M485120",
    "_widget_1689822029731": "中性的角度讲连体裤",
    "_widget_1689822029732": {
      "_id": "64e3a5d0a7e4f10008878e55",
      "name": "Candy",
      "status": 1,
      "type": 0,
      "username": "3faa4687"
    },
    "_widget_1689822553756": {
      "_id": "64e3a5d0a7e4f10008878e55",
      "name": "Candy",
      "status": 1,
      "type": 0,
      "username": "3faa4687"
    },
    "_widget_1689822553757": null,
    "_widget_1689822553763": [],
    "_widget_1689825348641": [
      {
        "_id": "651fc1b9550d16a33b372bae",
        "_widget_1689825348643": "M241S016",
        "_widget_1689825348644": "test0920003",
        "_widget_1689825348645": "朗格伦黄",
        "_widget_1689825348648": "1.2",
        "_widget_1689825348649": "20",
        "_widget_1689825348651": "上身",
        "_widget_1689825348653": {
          "id": "650a96e0a4191b8c9b3cd05e"
        },
        "_widget_1692277461302": "个",
        "_widget_1692277461304": "Lenglen Yellow",
        "_widget_1692354341637": "18",
        "_widget_1692501608712": 2,
        "_widget_1692526054992": "",
        "_widget_1692526054993": [],
        "_widget_1692526055059": "面料",
        "_widget_1692526055061": {},
        "_widget_1692527336030": "",
        "_widget_1694171032467": {
          "id": "64e45fd322d50900079e8c46"
        },
        "_widget_1694855597949": {
          "id": "6507a356791f400008f96643"
        },
        "_widget_1695299267678": "A款色组-b7e875fa-983e-4070-b13e-d77515392789",
        "_widget_1695550458810": 0.7,
        "_widget_1695550458841": "1",
        "_widget_1695550458871": 1,
        "_widget_1696423628948": {
          "_id": "64e19bfc9075e50007bf02ae",
          "name": "moodytiger",
          "status": 1,
          "type": 0,
          "username": "#admin"
        },
        "_widget_1696423628951": "供应商B",
        "_widget_1696424835454": "10002"
      },
      {
        "_id": "651fc1b9550d16a33b372baf",
        "_widget_1689825348643": "M241S015",
        "_widget_1689825348644": "test0920002",
        "_widget_1689825348645": "朗格伦绿",
        "_widget_1689825348648": "",
        "_widget_1689825348649": "",
        "_widget_1689825348651": "小肩长",
        "_widget_1689825348653": {
          "id": "650a9610e8eedc202fcda1d2"
        },
        "_widget_1692277461302": "个",
        "_widget_1692277461304": "Lenglen Green",
        "_widget_1692354341637": "",
        "_widget_1692501608712": 2,
        "_widget_1692526054992": "",
        "_widget_1692526054993": [],
        "_widget_1692526055059": "面料",
        "_widget_1692526055061": {},
        "_widget_1692527336030": "",
        "_widget_1694171032467": {
          "id": "64e45fd322d50900079e8c44"
        },
        "_widget_1694855597949": {
          "id": "6514dce04cf6368fee3461c1"
        },
        "_widget_1695299267678": "A款色组-56e8c822-8429-4429-a649-3f54da5df4dc",
        "_widget_1695550458810": 0.2,
        "_widget_1695550458841": "2",
        "_widget_1695550458871": 1,
        "_widget_1696423628948": null,
        "_widget_1696423628951": "",
        "_widget_1696424835454": ""
      },
      {
        "_id": "651fc1b9550d16a33b372bb0",
        "_widget_1689825348643": "M241S014",
        "_widget_1689825348644": "test0920001",
        "_widget_1689825348645": "芬馥灰",
        "_widget_1689825348648": "",
        "_widget_1689825348649": "",
        "_widget_1689825348651": "脚围",
        "_widget_1689825348653": {
          "id": "650a94d7df1a2f4a0c087be8"
        },
        "_widget_1692277461302": "个",
        "_widget_1692277461304": "Fragrant Gray",
        "_widget_1692354341637": "",
        "_widget_1692501608712": 2,
        "_widget_1692526054992": "",
        "_widget_1692526054993": [],
        "_widget_1692526055059": "面料",
        "_widget_1692526055061": {},
        "_widget_1692527336030": "",
        "_widget_1694171032467": {
          "id": "64e45fd322d50900079e8c42"
        },
        "_widget_1694855597949": {
          "id": "6514dcf286cd3db3fcae9f91"
        },
        "_widget_1695299267678": "A款色组-827f138e-b37b-41e5-84ac-36fdd576a581",
        "_widget_1695550458810": 0.1,
        "_widget_1695550458841": "3",
        "_widget_1695550458871": 1,
        "_widget_1696423628948": null,
        "_widget_1696423628951": "",
        "_widget_1696424835454": ""
      },
      {
        "_id": "652f9ad5743ac3ca9b91aa9a",
        "_widget_1689825348643": "T23FWP001",
        "_widget_1689825348644": "IP系列立体硅胶烫标",
        "_widget_1689825348645": "暖流绿",
        "_widget_1689825348648": "",
        "_widget_1689825348649": "",
        "_widget_1689825348651": "后中长",
        "_widget_1689825348653": {},
        "_widget_1692277461302": "个",
        "_widget_1692277461304": "Current Green",
        "_widget_1692354341637": "",
        "_widget_1692501608712": 1,
        "_widget_1692526054992": "",
        "_widget_1692526054993": [],
        "_widget_1692526055059": "辅料",
        "_widget_1692526055061": {
          "id": "652496bd207b03de7b45c9c7"
        },
        "_widget_1692527336030": "35*35mm",
        "_widget_1694171032467": {
          "id": "64e45fd322d50900079e8bea"
        },
        "_widget_1694855597949": {
          "id": "6514dcd2c7d04664fb4f72e3"
        },
        "_widget_1695299267678": "A款色组-3d76dd64-4bda-4554-923f-2b3e27fa8968",
        "_widget_1695550458810": null,
        "_widget_1695550458841": "4",
        "_widget_1695550458871": 0,
        "_widget_1696423628948": null,
        "_widget_1696423628951": "东莞市冠荣商标织造有很公司",
        "_widget_1696424835454": ""
      }
    ],
    "_widget_1692016855698": "否",
    "_widget_1692021492114": "",
    "_widget_1692338476270": "2025-0148",
    "_widget_1692338476271": "辅助",
    "_widget_1692338476272": "moodytiger-2025年春夏企划案",
    "_widget_1692338999133": [
      {
        "_id": "651fc1b9550d16a33b372bad",
        "_widget_1692338999135": "（成人）浅灰白",
        "_widget_1692338999136": "Blanc de Blanc",
        "_widget_1692354341609": 1,
        "_widget_1692458015476": "正常",
        "_widget_1694855115708": "0914"
      }
    ],
    "_widget_1692340738228": "",
    "_widget_1692353365138": "45",
    "_widget_1692353365151": "",
    "_widget_1692354205554": [],
    "_widget_1692354205555": [],
    "_widget_1692354341623": 1,
    "_widget_1692354341730": "",
    "_widget_1692504409127": "供应商B",
    "_widget_1692527689426": [
      {
        "_id": "651fc1b9550d16a33b372bb1",
        "_widget_1692527689427": "面料",
        "_widget_1692527689428": {},
        "_widget_1692527689429": {},
        "_widget_1692527689431": "",
        "_widget_1692527689432": "",
        "_widget_1692527689433": "",
        "_widget_1692527689434": "",
        "_widget_1692527689435": "",
        "_widget_1692527689436": [],
        "_widget_1692527689437": "",
        "_widget_1692527689438": "",
        "_widget_1692527689439": "",
        "_widget_1692527689440": "",
        "_widget_1692527689441": "",
        "_widget_1692527689442": "",
        "_widget_1692527689445": null,
        "_widget_1694855597946": {},
        "_widget_1694855597950": {},
        "_widget_1695299267679": "",
        "_widget_1695550458872": null,
        "_widget_1695550458875": 1,
        "_widget_1695550458878": "",
        "_widget_1696423628956": "",
        "_widget_1696423628957": null,
        "_widget_1696424835455": ""
      }
    ],
    "_widget_1692856726767": [
      {
        "_id": "651fc1b9550d16a33b372bb2",
        "_widget_1692856726768": "面料",
        "_widget_1692856726769": {},
        "_widget_1692856726770": {},
        "_widget_1692856726772": "",
        "_widget_1692856726773": "",
        "_widget_1692856726774": "",
        "_widget_1692856726775": "",
        "_widget_1692856726776": "",
        "_widget_1692856726777": [],
        "_widget_1692856726778": "",
        "_widget_1692856726779": "",
        "_widget_1692856726780": "",
        "_widget_1692856726781": "",
        "_widget_1692856726782": "",
        "_widget_1692856726783": "",
        "_widget_1692856726786": null,
        "_widget_1694855597985": {},
        "_widget_1694855597986": {},
        "_widget_1695299267680": "",
        "_widget_1695550458937": 1,
        "_widget_1695550458938": null,
        "_widget_1695550458940": "",
        "_widget_1696423628958": "",
        "_widget_1696423628959": null,
        "_widget_1696424835456": ""
      }
    ],
    "_widget_1692856726792": [
      {
        "_id": "651fc1b9550d16a33b372bb3",
        "_widget_1692856726793": "面料",
        "_widget_1692856726794": {},
        "_widget_1692856726795": {},
        "_widget_1692856726797": "",
        "_widget_1692856726798": "",
        "_widget_1692856726799": "",
        "_widget_1692856726800": "",
        "_widget_1692856726801": "",
        "_widget_1692856726802": [],
        "_widget_1692856726803": "",
        "_widget_1692856726804": "",
        "_widget_1692856726805": "",
        "_widget_1692856726806": "",
        "_widget_1692856726807": "",
        "_widget_1692856726808": "",
        "_widget_1692856726811": null,
        "_widget_1695435901850": {},
        "_widget_1695435901852": {},
        "_widget_1695550458969": "",
        "_widget_1695550458972": "",
        "_widget_1695550458973": 1,
        "_widget_1695550458974": null,
        "_widget_1696423628960": "",
        "_widget_1696423628961": null,
        "_widget_1696424835457": ""
      }
    ],
    "_widget_1692856726817": [
      {
        "_id": "651fc1b9550d16a33b372bb4",
        "_widget_1692856726818": "面料",
        "_widget_1692856726819": {},
        "_widget_1692856726820": {},
        "_widget_1692856726822": "",
        "_widget_1692856726823": "",
        "_widget_1692856726824": "",
        "_widget_1692856726825": "",
        "_widget_1692856726826": "",
        "_widget_1692856726827": [],
        "_widget_1692856726828": "",
        "_widget_1692856726829": "",
        "_widget_1692856726830": "",
        "_widget_1692856726831": "",
        "_widget_1692856726832": "",
        "_widget_1692856726833": "",
        "_widget_1692856726836": null,
        "_widget_1695435901853": {},
        "_widget_1695435901861": {},
        "_widget_1695550458975": "",
        "_widget_1695550458976": 1,
        "_widget_1695550458977": null,
        "_widget_1695550458990": "",
        "_widget_1696423628962": "",
        "_widget_1696423628963": null,
        "_widget_1696424835458": ""
      }
    ],
    "_widget_1692856726842": [
      {
        "_id": "651fc1b9550d16a33b372bb5",
        "_widget_1692856726843": "面料",
        "_widget_1692856726844": {},
        "_widget_1692856726845": {},
        "_widget_1692856726847": "",
        "_widget_1692856726848": "",
        "_widget_1692856726849": "",
        "_widget_1692856726850": "",
        "_widget_1692856726851": "",
        "_widget_1692856726852": [],
        "_widget_1692856726853": "",
        "_widget_1692856726854": "",
        "_widget_1692856726855": "",
        "_widget_1692856726856": "",
        "_widget_1692856726857": "",
        "_widget_1692856726858": "",
        "_widget_1692856726861": null,
        "_widget_1695435901855": {},
        "_widget_1695435901856": {},
        "_widget_1695550458980": "",
        "_widget_1695550458981": 1,
        "_widget_1695550458982": null,
        "_widget_1695550458991": "",
        "_widget_1696423628964": "",
        "_widget_1696423628965": null,
        "_widget_1696424835459": ""
      }
    ],
    "_widget_1692856726867": [
      {
        "_id": "651fc1b9550d16a33b372bb6",
        "_widget_1692856726868": "面料",
        "_widget_1692856726869": {},
        "_widget_1692856726870": {},
        "_widget_1692856726872": "",
        "_widget_1692856726873": "",
        "_widget_1692856726874": "",
        "_widget_1692856726875": "",
        "_widget_1692856726876": "",
        "_widget_1692856726877": [],
        "_widget_1692856726878": "",
        "_widget_1692856726879": "",
        "_widget_1692856726880": "",
        "_widget_1692856726881": "",
        "_widget_1692856726882": "",
        "_widget_1692856726883": "",
        "_widget_1692856726886": null,
        "_widget_1695435901858": {},
        "_widget_1695435901860": {},
        "_widget_1695550458985": "",
        "_widget_1695550458986": 1,
        "_widget_1695550458987": null,
        "_widget_1695550458992": "",
        "_widget_1696423628966": "",
        "_widget_1696423628967": null,
        "_widget_1696424835460": ""
      }
    ],
    "_widget_1692917265214": "头版样衣已完成",
    "_widget_1692918710402": "",
    "_widget_1692919566102": [
      {
        "mime": "image/png",
        "name": "image.png",
        "size": 32293,
        "url": "https://files.jiandaoyun.com/78a63ef7-456d-4da5-a62f-c45386c2a4f7?attname=image.png&e=1698915599&token=bM7UwVPyBBdPaleBZt21SWKzMylqPUpn-05jZlas:Qqkzz4YvjZ5XbnhMCHEfyI2qIg8="
      }
    ],
    "_widget_1693488251032": "",
    "_widget_1693489286370": "8",
    "_widget_1693489286371": "5",
    "_widget_1693489286372": "1",
    "_widget_1693489286373": "20",
    "_widget_1693812368356": [
      {
        "_id": "651fc1b9550d16a33b372bac",
        "_widget_1693812368358": "",
        "_widget_1693812368361": "",
        "_widget_1693812368362": "",
        "_widget_1693812368363": 1,
        "_widget_1693812368364": "fb469883-c364-4b5e-9f96-80e9f0ae1529",
        "_widget_1693812368366": null,
        "_widget_1693840358497": [],
        "_widget_1693840358498": [],
        "_widget_1694855385428": {},
        "_widget_1694855597796": {}
      }
    ],
    "_widget_1693812368365": "fb469883-c364-4b5e-9f96-80e9f0ae1529",
    "_widget_1694405121124": {
      "_id": "64ee058587c4940008bd3909",
      "name": "Scrum",
      "status": 1,
      "type": 0,
      "username": "gece8b53"
    },
    "_widget_1694492853947": "供应商C",
    "_widget_1694493028787": null,
    "_widget_1694588675961": {
      "_id": "64ee058587c4940008bd3908",
      "name": "Frank Cao",
      "status": 1,
      "type": 0,
      "username": "aecf91e1"
    },
    "_widget_1694588675962": {
      "_id": "64ee91fe314eac0007be65e9",
      "name": "Damon Song",
      "status": 1,
      "type": 0,
      "username": "de2666d4"
    },
    "_widget_1694590177852": {
      "_id": "64f033cd80f03e00085415b2",
      "name": "rogge",
      "status": 1,
      "type": 0,
      "username": "6ge6ac85"
    },
    "_widget_1694855597871": {},
    "_widget_1695297446692": [],
    "_widget_1695550458812": "A款色组-b7e875fa-983e-4070-b13e-d77515392789A款色组-56e8c822-8429-4429-a649-3f54da5df4dcA款色组-827f138e-b37b-41e5-84ac-36fdd576a581A款色组-3d76dd64-4bda-4554-923f-2b3e27fa8968",
    "_widget_1695550458870": "0.7,0.2,0.1",
    "_widget_1695550458873": "",
    "_widget_1695550458874": "1,",
    "_widget_1695550458935": "",
    "_widget_1695550458936": "1,",
    "_widget_1695550458970": "",
    "_widget_1695550458971": "1,",
    "_widget_1695550458978": "",
    "_widget_1695550458979": "1,",
    "_widget_1695550458983": "",
    "_widget_1695550458984": "",
    "_widget_1695550458988": "",
    "_widget_1695550458989": "1,",
    "_widget_1696433218443": "",
    "_widget_1696668898618": "",
    "_widget_1696854858502": [
      "头版",
      "复版",
      "拍照样"
    ],
    "_widget_1696854858505": 1,
    "_widget_1696854858506": 1,
    "_widget_1696854858507": 1,
    "appId": "64e1f2d0d8b6100009e88053",
    "bdb_list": "20,,,",
    "bw_list": "上身,小肩长,脚围,后中长",
    "createTime": "2023-10-06T08:13:45.030Z",
    "creator": {
      "_id": "64ffe045479d570008d00f8d",
      "name": "Liena Li",
      "status": 1,
      "type": 0,
      "username": "fbge8ccc"
    },
    "deleteTime": null,
    "deleter": null,
    "design_file": [
      {
        "mime": "image/png",
        "name": "image.png",
        "size": 32293,
        "url": "https://files.jiandaoyun.com/7f1e7b4a-9236-4c91-a1de-631cc293653e?attname=image.png&e=1698915599&token=bM7UwVPyBBdPaleBZt21SWKzMylqPUpn-05jZlas:x_uYohhUoIF8TYQ6KucYJgaZnas="
      }
    ],
    "dw_list": "个,个,个,个",
    "entryId": "64b89ecaf7997d00085e492e",
    "formName": "版单设计",
    "gg_list": ",,,35*35mm",
    "ilmc_list": "test0920003,test0920002,test0920001,IP系列立体硅胶烫标",
    "kz_list": "1.2,,,",
    "sf_list": "18,,,",
    "style_number": "M48512045",
    "updateTime": "2023-10-18T08:44:04.724Z",
    "updater": {
      "_id": "64e19bfc9075e50007bf02ae",
      "name": "moodytiger",
      "status": 1,
      "type": 0,
      "username": "#admin"
    },
    "wild_list": "面料,面料,面料,辅料",
    "wl_list": "M241S016,M241S015,M241S014,T23FWP001"
  }
}`
