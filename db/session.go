package db

import (
	"encoding/json"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _session *mgo.Session = nil

func SetSession(s *mgo.Session) {
	_session = s
}

func Insert(ctx sdk.Context, collectionName string, data interface{}, subData map[string]interface{}) {
	if ctx.IsCheckTx() == false && _session != nil {
		go insertRoutine(ctx, collectionName, data, subData)
	}
}

func insertRoutine(ctx sdk.Context, collectionName string, data interface{}, subData map[string]interface{}) {
	if ctx.IsCheckTx() == false && _session != nil {
		defer func() {
			if r := recover(); r != nil {
				session := _session.Clone()
				collection := session.DB("").C("errors")
				message := fmt.Sprintf("Fail to insert in %s: %s", collectionName, r)
				collection.Insert(map[string]interface{}{"message": message, "block-height": ctx.BlockHeight()})
			}
		}()

		session := _session.Clone()
		collection := session.DB("").C(collectionName)
		var d interface{}
		jsonBytes, _ := json.Marshal(data)
		json.Unmarshal(jsonBytes, &d)
		mapData := d.(map[string]interface{})
		if subData != nil {
			for key, value := range subData {
				mapData[key] = value
			}
		}
		mapData["createTimestamp"] = ctx.BlockHeader().Time
		mapData["lastUpdateTimestamp"] = ctx.BlockHeader().Time
		err := collection.Insert(mapData)
		if err != nil {
			panic(err)
		}
	}
}

func UpdateSilently(ctx sdk.Context, collectionName string, query map[string]interface{}, data interface{}, subData map[string]interface{}) {
	if ctx.IsCheckTx() == false && _session != nil {
		go updateSilentlyRoutine(ctx, collectionName, query, data, subData)
	}
}

func updateSilentlyRoutine(ctx sdk.Context, collectionName string, query map[string]interface{}, data interface{}, subData map[string]interface{}) {
	if ctx.IsCheckTx() == false && _session != nil {
		defer func() {
			if r := recover(); r != nil {
				session := _session.Clone()
				collection := session.DB("").C("errors")
				message := fmt.Sprintf("Fail to upsert in %s: %s", collectionName, r)
				collection.Insert(map[string]interface{}{"message": message, "block-height": ctx.BlockHeight()})
			}
		}()

		session := _session.Clone()
		collection := session.DB("").C(collectionName)
		var d interface{}
		jsonBytes, _ := json.Marshal(data)
		json.Unmarshal(jsonBytes, &d)
		mapData := d.(map[string]interface{})
		if subData != nil {
			for key, value := range subData {
				mapData[key] = value
			}
		}

		var incData map[string]interface{} = nil
		incDataI, incOk := mapData["$inc"]
		if incOk {
			delete(mapData, "$inc")
			_incData, ok := incDataI.(map[string]interface{})
			if ok {
				incData = _incData
			}
		}

		_mapData := make(map[string]interface{})
		if len(mapData) > 0 {
			_mapData["$set"] = mapData
		}
		if incData != nil {
			_mapData["$inc"] = incData
		}
		err := collection.Update(query, _mapData)
		if err != nil {
			panic(err)
		}
	}
}

func Update(ctx sdk.Context, collectionName string, query map[string]interface{}, data interface{}, subData map[string]interface{}) {
	if ctx.IsCheckTx() == false && _session != nil {
		go updateRoutine(ctx, collectionName, query, data, subData)
	}
}

func updateRoutine(ctx sdk.Context, collectionName string, query map[string]interface{}, data interface{}, subData map[string]interface{}) {
	if ctx.IsCheckTx() == false && _session != nil {
		defer func() {
			if r := recover(); r != nil {
				session := _session.Clone()
				collection := session.DB("").C("errors")
				message := fmt.Sprintf("Fail to upsert in %s: %s", collectionName, r)
				collection.Insert(map[string]interface{}{"message": message, "block-height": ctx.BlockHeight()})
			}
		}()

		session := _session.Clone()
		collection := session.DB("").C(collectionName)
		var d interface{}
		jsonBytes, _ := json.Marshal(data)
		json.Unmarshal(jsonBytes, &d)
		mapData := d.(map[string]interface{})
		if subData != nil {
			for key, value := range subData {
				mapData[key] = value
			}
		}
		mapData["lastUpdateTimestamp"] = ctx.BlockHeader().Time

		var incData map[string]interface{} = nil
		incDataI, incOk := mapData["$inc"]
		if incOk {
			delete(mapData, "$inc")
			_incData, ok := incDataI.(map[string]interface{})
			if ok {
				incData = _incData
			}
		}

		_mapData := make(map[string]interface{})
		if len(mapData) > 0 {
			_mapData["$set"] = mapData
		}
		if incData != nil {
			_mapData["$inc"] = incData
		}
		err := collection.Update(query, _mapData)
		if err != nil {
			panic(err)
		}
	}
}

func Upsert(ctx sdk.Context, collectionName string, query map[string]interface{}, data interface{}, subData map[string]interface{}) {
	if ctx.IsCheckTx() == false && _session != nil {
		go upsertRotine(ctx, collectionName, query, data, subData)
	}
}

func upsertRotine(ctx sdk.Context, collectionName string, query map[string]interface{}, data interface{}, subData map[string]interface{}) {
	if ctx.IsCheckTx() == false && _session != nil {
		defer func() {
			if r := recover(); r != nil {
				session := _session.Clone()
				collection := session.DB("").C("errors")
				message := fmt.Sprintf("Fail to upsert in %s: %s", collectionName, r)
				collection.Insert(map[string]interface{}{"message": message, "block-height": ctx.BlockHeight()})
			}
		}()

		session := _session.Clone()
		collection := session.DB("").C(collectionName)
		var d interface{}
		jsonBytes, _ := json.Marshal(data)
		json.Unmarshal(jsonBytes, &d)
		mapData := d.(map[string]interface{})
		if subData != nil {
			for key, value := range subData {
				mapData[key] = value
			}
		}
		mapData["lastUpdateTimestamp"] = ctx.BlockHeader().Time

		var incData map[string]interface{} = nil
		incDataI, incOk := mapData["$inc"]
		if incOk {
			delete(mapData, "$inc")
			_incData, ok := incDataI.(map[string]interface{})
			if ok {
				incData = _incData
			}
		}

		_mapData := make(map[string]interface{})
		_mapData["$setOnInsert"] = bson.M{"createTimestamp": ctx.BlockHeader().Time}
		if len(mapData) > 0 {
			_mapData["$set"] = mapData
		}
		if incData != nil {
			_mapData["$inc"] = incData
		}
		_, err := collection.Upsert(query, _mapData)
		if err != nil {
			panic(err)
		}
	}
}
