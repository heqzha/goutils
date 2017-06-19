package db

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBHandler struct {
	se          *mgo.Session
	url         string
	ErrNotFound error
}

type BsonM bson.M

func MongoDBNewHandler(username, password, db string, url ...string) (*MongoDBHandler, error) {
	jurl := strings.Join(url, ",")
	var furl string
	if db == "" {
		db = "default"
	}

	if username != "" {
		furl = fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, jurl, db)
	} else {
		furl = fmt.Sprintf("mongodb://%s/%s", jurl, db)
	}

	se, err := mgo.Dial(furl)
	if err != nil {
		return nil, err
	}
	return &MongoDBHandler{
		se:          se,
		url:         furl,
		ErrNotFound: mgo.ErrNotFound,
	}, nil
}

func (h *MongoDBHandler) Renew() (*MongoDBHandler, error) {
	se, err := mgo.Dial(h.url)
	if err != nil {
		return nil, err
	}
	return &MongoDBHandler{
		se:  se,
		url: h.url,
	}, nil
}

func (h *MongoDBHandler) Close() {
	h.se.Close()
}

func (h *MongoDBHandler) NewID() string {
	return bson.NewObjectId().Hex()
}

func (h *MongoDBHandler) CheckID(id string) bool {
	return bson.IsObjectIdHex(id)
}

func (h *MongoDBHandler) ToObjectID(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}

func (h *MongoDBHandler) EnsureIndex(db, cName string, unique, dropDups, background, sparse bool, keys ...string) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.EnsureIndex(mgo.Index{
		Key:        keys,
		Unique:     unique,
		DropDups:   dropDups,
		Background: background,
		Sparse:     sparse,
	})
}

func (h *MongoDBHandler) DropIndex(db, cName string, keys ...string) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.DropIndex(keys...)
}

func (h *MongoDBHandler) Indexes(db, cName string) ([]map[string]interface{}, error) {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	indexes, err := c.Indexes()
	if err != nil {
		return nil, err
	}

	m := []map[string]interface{}{}
	for _, index := range indexes {
		mIdx := structs.Map(index)
		m = append(m, mIdx)
	}
	return m, nil
}

func (h *MongoDBHandler) FindAll(db, cName string, selector BsonM, offset, limit int, sort []string, results interface{}) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	if len(sort) != 0 {
		return c.Find(selector).Sort(sort...).Skip(offset).Limit(limit).All(results)
	}
	if limit > 0 {
		return c.Find(selector).Skip(offset).Limit(limit).All(results)
	}
	return c.Find(selector).Skip(offset).All(results)
}

func (h *MongoDBHandler) FindOne(db, cName string, selector BsonM, offset, limit int, sort []string, result interface{}) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	if len(sort) != 0 {
		return c.Find(selector).Sort(sort...).Skip(offset).Limit(limit).One(result)
	}
	return c.Find(selector).Skip(offset).Limit(limit).One(result)
}

func (h *MongoDBHandler) Find(db, cName string, selector BsonM, result interface{}) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.Find(selector).One(result)
}

func (h *MongoDBHandler) Exist(db, cName string, selector BsonM) (bool, error) {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	cnt, err := c.Find(selector).Count()
	if err != nil {
		return false, err
	} else if cnt > 0 {
		return true, nil
	}
	return false, nil
}

func (h *MongoDBHandler) FindByID(db, cName, id string, result interface{}) error {
	if !h.CheckID(id) {
		return fmt.Errorf("Invalid ID %s", id)
	}
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.FindId(bson.ObjectIdHex(id).Hex()).One(result)
}

func (h *MongoDBHandler) CountByID(db, cName, id string) (int, error) {
	if !h.CheckID(id) {
		return 0, fmt.Errorf("Invalid ID %s", id)
	}
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.FindId(bson.ObjectIdHex(id).Hex()).Count()
}

func (h *MongoDBHandler) Count(db, cName string, selector BsonM) (int, error) {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.Find(selector).Count()
}

func (h *MongoDBHandler) Distinct(db, cName, key string, selector BsonM, results interface{}) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.Find(selector).Distinct(key, results)
}

func (h *MongoDBHandler) Insert(db, cName string, cObjects ...interface{}) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.Insert(cObjects...)
}

func (h *MongoDBHandler) Update(db, cName string, selector BsonM, cObject interface{}) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.Update(selector, cObject)
}

func (h *MongoDBHandler) UpdateByID(db, cName, id string, cObject interface{}) error {
	if !h.CheckID(id) {
		return fmt.Errorf("Invalid ID %s", id)
	}
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.UpdateId(bson.ObjectIdHex(id).Hex(), cObject)
}

func (h *MongoDBHandler) Upsert(db, cName string, selector BsonM, cObject interface{}) (int, error) {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	info, err := c.Upsert(selector, cObject)
	if err != nil {
		return 0, err
	}
	return info.Matched, nil
}

func (h *MongoDBHandler) UpsertedId(db, cName, id string, cObject interface{}) (int, error) {
	if !h.CheckID(id) {
		return 0, fmt.Errorf("Invalid ID %s", id)
	}
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	info, err := c.UpsertId(bson.ObjectIdHex(id).Hex(), cObject)
	if err != nil {
		return 0, err
	}
	return info.Matched, nil
}

func (h *MongoDBHandler) Remove(db, cName string, selector BsonM) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.Remove(selector)
}

func (h *MongoDBHandler) RemoveAll(db, cName string, selector BsonM) (int, error) {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	info, err := c.RemoveAll(selector)
	if err != nil {
		return 0, err
	}
	return info.Removed, nil
}

func (h *MongoDBHandler) RemoveByID(db, cName, id string) error {
	if !h.CheckID(id) {
		return fmt.Errorf("Invalid ID %s", id)
	}
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	return c.RemoveId(bson.ObjectIdHex(id).Hex())
}

func (h *MongoDBHandler) findAndModifyByID(cmd, db, cName, id string, updater BsonM, result interface{}) error {
	if !h.CheckID(id) {
		return fmt.Errorf("Invalid ID %s", id)
	}
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	change := mgo.Change{
		Update:    bson.M{cmd: updater},
		ReturnNew: true,
	}
	_, err := c.FindId(bson.ObjectIdHex(id).Hex()).Apply(change, result)
	return err
}

func (h *MongoDBHandler) findAndModify(cmd, db, cName string, selector BsonM, updater BsonM, result interface{}) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	change := mgo.Change{
		Update:    bson.M{cmd: updater},
		ReturnNew: true,
	}
	_, err := c.Find(selector).Apply(change, result)
	return err
}

func (h *MongoDBHandler) Inc(db, cName string, selector BsonM, key string, value int, result interface{}) error {
	return h.findAndModify("$inc", db, cName, selector, BsonM{key: value}, result)
}

func (h *MongoDBHandler) Set(db, cName string, selector BsonM, key string, value interface{}, result interface{}) error {
	return h.findAndModify("$set", db, cName, selector, BsonM{key: value}, result)
}

func (h *MongoDBHandler) IncByID(db, cName, id, key string, value int, result interface{}) error {
	return h.findAndModifyByID("$inc", db, cName, id, BsonM{key: value}, result)
}

func (h *MongoDBHandler) SetByID(db, cName, id, key string, value interface{}, result interface{}) error {
	return h.findAndModifyByID("$set", db, cName, id, BsonM{key: value}, result)
}

func (h *MongoDBHandler) Pipe(db, cName string, selectors []BsonM, results interface{}) error {
	se := h.se.Copy()
	defer se.Close()
	c := se.DB(db).C(cName)

	if err := c.Pipe(selectors).All(results); err != nil {
		return err
	}
	return nil
}
