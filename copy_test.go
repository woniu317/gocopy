package gocopy_test

import (
	"github.com/tidwall/gjson"
	"gocopy"
	"strconv"
	"testing"
)

type baseInfo struct {
	Name    string `json:"name"`
	Id      int    `json:"id"`
	Address string `json:"address"`
}

type personInfo struct {
	BaseInfo baseInfo `json:"base_info"`
	Company  string   `json:"company"`
	Money    float64  `json:"money"` // 单位元
	HorseId  string   `json:"horse_id"`
}

type databaseInfo struct {
	Name    string `cp:"base_info.name"`
	Id      int    `cp:"base_info.id"`
	Address string `cp:"base_info.address"`
	Company string `cp:"company"`
	Money   int64  `cp:"money"` // 单位分
	HorseId int    `cp:"horse_id"`
}

func moneyToInt(money gjson.Result) interface{} {
	return int64(money.Float() * 100)
}

func strToInt(str gjson.Result) interface{} {
	v, _ := strconv.Atoi(str.Str)
	return v
}

var person = personInfo{
	BaseInfo: baseInfo{
		Name:    "base",
		Id:      123,
		Address: "望京soho",
	},
	Company: "云账户（天津）",
	Money:   100000.25,
	HorseId: "401",
}

func TestCopy(t *testing.T) {

	var db databaseInfo
	err := gocopy.ConvertToTarget(person, &db,
		"cp", map[string]func(result gjson.Result) interface{}{
			"money":    moneyToInt,
			"horse_id": strToInt,
		})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(person.BaseInfo.Name == db.Name)
	t.Log(person.BaseInfo.Address == db.Address)
	t.Log(person.BaseInfo.Id == db.Id)
	t.Log(person.Company == db.Company)
	t.Log(int64(person.Money*100) == db.Money)
	t.Log(person.HorseId == strconv.Itoa(db.HorseId))
}

func TestErrorCopy(t *testing.T) {

	var db databaseInfo
	err := gocopy.ConvertToTarget(person, db, "cp", map[string]func(result gjson.Result) interface{}{
		"money": moneyToInt,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(person.BaseInfo.Name == db.Name)
	t.Log(person.BaseInfo.Address == db.Address)
	t.Log(person.BaseInfo.Id == db.Id)
	t.Log(person.Company == db.Company)
	t.Log(int64(person.Money*100) == db.Money)
}
