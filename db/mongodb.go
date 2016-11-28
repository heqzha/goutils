package db

import(
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBHandler struct{
	se *mgo.Session
	url string
}

type BsonM bson.M

func MongoDBNewHandler(username, password, db string, url ...string)(*MongoDBHandler, error){
	jurl := strings.Join(url, ",")
	var furl string
	if db == ""{
		db = "default"
	}

	if username != "" {
		furl = fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, jurl, db)
	}else {
		furl = fmt.Sprintf("mongodb://%s/%s", jurl, db)
	}

	se, err := mgo.Dial(furl)
	if err != nil{
		return nil, err
	}
	return &MongoDBHandler{
		se: se,
		url: furl,
	}, nil
}

func (h *MongoDBHandler)Renew()(*MongoDBHandler, error){
	se, err := mgo.Dial(h.url)
	if err != nil{
		return nil, err
	}
	return &MongoDBHandler{
		se: se,
		url: h.url,
	}, nil
}

func (h *MongoDBHandler)Close(){
	h.se.Close()
}

func (h *MongoDBHandler)NewID()string{
	return bson.NewObjectId().Hex()
}

func (h *MongoDBHandler)CheckID(id string)bool{
	return bson.IsObjectIdHex(id)
}

func (h *MongoDBHandler)ToObjectID(id string)bson.ObjectId{
	return bson.ObjectIdHex(id)
}

func (h *MongoDBHandler)EnsureIndex(db, cName string, unique, dropDups, background, sparse bool, keys ...string) error{
	c := h.se.DB(db).C(cName)
	return c.EnsureIndex(mgo.Index{
		Key: keys,
		Unique: unique,
		DropDups: dropDups,
		Background: background,
		Sparse: sparse,
	})
}

func (h *MongoDBHandler)DropIndex(db, cName string, keys ...string) error{
	c := h.se.DB(db).C(cName)
	return c.DropIndex(keys...)
}

func (h *MongoDBHandler)Indexes(db, cName string)([]map[string]interface{}, error){
	c := h.se.DB(db).C(cName)
	indexes, err := c.Indexes()
	if err != nil{
		return nil, err
	}

	m := []map[string]interface{}{}
	for _, index := range indexes{
		mIdx := structs.Map(index)
		m = append(m, mIdx)
	}
	return m, nil
}

func (h *MongoDBHandler)find(db, cName string, selector BsonM) (*mgo.Query){
	c := h.se.DB(db).C(cName)
	return c.Find(selector)
}

func (h *MongoDBHandler)findByID(db, cName string, id string) (*mgo.Query){
	c := h.se.DB(db).C(cName)
	return c.FindId(bson.ObjectIdHex(id))
}

func (h *MongoDBHandler)FindAll(db, cName string, offset, limit int, sort []string, results interface{}) error{
	if len(sort) != 0{
		return h.find(db, cName, nil).Sort(sort...).Skip(offset).Limit(limit).All(results)
	}
	return h.find(db, cName, nil).Skip(offset).Limit(limit).All(results)
}

func (h *MongoDBHandler)FindOne(db, cName string, selector BsonM, offset, limit int, sort []string, result interface{}) error{
	if len(sort) != 0{
		return h.find(db, cName, selector).Sort(sort...).Skip(offset).Limit(limit).One(result)
	}
	return h.find(db, cName, selector).Skip(offset).Limit(limit).One(result)
}

func (h *MongoDBHandler)Find(db, cName string, selector BsonM, result interface{}) error{
	return h.find(db, cName, selector).One(result)
}

func (h *MongoDBHandler)FindByID(db, cName, id string, result interface{}) error{
	return h.findByID(db, cName, id).One(result)
}

func (h *MongoDBHandler)Insert(db, cName string, cObjects...interface{}) error{
	c := h.se.DB(db).C(cName)
	return c.Insert(cObjects...)
}

func (h *MongoDBHandler)Upsert(db, cName string, selector BsonM, cObject interface{}) (int, error){
	c := h.se.DB(db).C(cName)
	info, err := c.Upsert(selector, cObject)
	if err != nil{
		return 0, err
	}
	return info.Matched, nil
}

func (h *MongoDBHandler)UpsertedId(db, cName, id string, cObject interface{})(string, error){
	c := h.se.DB(db).C(cName)
	info, err := c.UpsertId(bson.ObjectIdHex(id), cObject)
	if err != nil{
		return "", err
	}
	return string(info.UpsertedId.(bson.ObjectId)), nil
}

func (h *MongoDBHandler)Remove(db, cName string, selector BsonM) error{
	c := h.se.DB(db).C(cName)
	return c.Remove(selector)
}

func (h *MongoDBHandler)RemoveAll(db, cName string, selector BsonM) (int, error){
	c := h.se.DB(db).C(cName)
	info, err := c.RemoveAll(selector)
	if err != nil{
		return 0, err
	}
	return info.Removed, nil
}

func (h *MongoDBHandler)RemoveByID(db, cName, id string) error{
	c := h.se.DB(db).C(cName)
	return c.RemoveId(bson.ObjectIdHex(id))
}
