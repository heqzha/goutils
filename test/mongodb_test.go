package test

import (
	"fmt"
	"github.com/heqzha/goutils/date"
	"github.com/heqzha/goutils/db"
	"testing"
)

type ATestData struct {
	ID        string      `bson:"_id"`
	Name      string      `bson:"name"`
	Age       int32       `bson:"age"`
	Data      interface{} `bson:"data"`
	Address   string      `bson:address`
	CreatedTs int64       `bson:"created_ts"`
}

func TestMongoDBHandler(t *testing.T) {
	dbName := "test"
	tName := "a_test_data"
	h, err := db.MongoDBNewHandler("", "", dbName, "127.0.0.1:27017")
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer h.Close()

	data := ATestData{
		ID:   h.NewID(),
		Name: "TestName",
		Age:  10,
		Data: map[string]interface{}{
			"sex":    true,
			"height": 150.0,
			"weight": 60.0,
			"other":  "other data",
		},
		Address:   "1234567",
		CreatedTs: date.DateNowSecond(),
	}
	if err := h.Insert(dbName, tName, data); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("Insert data: %v\n", data))

	result := ATestData{}
	if err := h.Find(dbName, tName, db.BsonM{
		"name": "TestName",
	}, &result); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("Find result: %v\n", result))

	data.Age++
	data.Address = ""
	numUpsert, err := h.Upsert(dbName, tName, db.BsonM{
		"_id": data.ID,
	}, data)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("Upsert num: %d\n", numUpsert))

	result = ATestData{}
	t.Log(data)
	if err := h.FindByID(dbName, tName, data.ID, &result); err != nil {
		t.Error(err.Error())
		return
	}
	result = ATestData{}
	if err := h.Inc(dbName, tName, db.BsonM{
		"_id": data.ID,
	}, "age", 1, &result); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("Inc result: %v\n", result))

	if err := h.IncByID(dbName, tName, data.ID, "age", -1, &result); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("IncByID result: %v\n", result))

	result = ATestData{}
	if err := h.Set(dbName, tName, db.BsonM{
		"_id": data.ID,
	}, "data.other", "aaaa", &result); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("Set result: %v\n", result))

	if err := h.SetByID(dbName, tName, data.ID, "data.sex", false, &result); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("SetByID result: %v\n", result))

	exist, err := h.Exist(dbName, tName, db.BsonM{
		"_id": data.ID,
	})
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("Exist result: %v\n", exist))

	if err := h.EnsureIndex(dbName, tName, false, false, false, false, "name", "age"); err != nil {
		t.Error(err.Error())
		return
	}
	indexes, err := h.Indexes(dbName, tName)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("Indexes: %v\n", indexes))

	if err := h.DropIndex(dbName, tName, "name", "age"); err != nil {
		t.Error(err.Error())
		return
	}

	if err := h.Remove(dbName, tName, db.BsonM{
		"age": db.BsonM{
			"$gt": 10,
		},
	}); err != nil {
		t.Error(err.Error())
		return
	}

	numRemoved, err := h.RemoveAll(dbName, tName, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(fmt.Sprintf("Removed num: %d\n", numRemoved))
}
