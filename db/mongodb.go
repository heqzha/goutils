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
	db string
}

func MongoDBNewHandler(username, password, dbName string, url ...string)(*MongoDBHandler, error){
	jurl := strings.Join(url, ",")
	var furl string
	if username != "" {
		furl := fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, jurl, dbName)
	}else if dbName != ""{
		furl := fmt.Sprintf("mongodb://%s/%s", jurl, dbName)
	}else{
		dbName = "default"
		furl := fmt.Sprintf("mongodb://%s/%s", jurl, dbName)
	}
	se, err := mgo.Dial(furl)
	if err != nil{
		return nil, err
	}
	return &MongoDBHandler{
		se: se,
		url: furl,
		db: dbName,
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
		db: h.db,
	}, nil
}

func (h *MongoDBHandler)Close(){
	h.se.Close
	mgo.Monotonic
}

func (h *MongoDBHandler)SelectDB(dbName string){
	h.db = dbName
}

func (h *MongoDBHandler)NewID()string{
	return string(bson.NewObjectId())
}

func (h *MongoDBHandler)EnsureIndex(cName string, unique, dropDups, background, sparse bool, keys ...string) error{
	c := h.se.DB(h.db).C(cName)
	return c.EnsureIndex(Index{
		Key: keys,
		Unique: unique,
		DropDups: dropDups,
		Background: background,
		Sparse: sparse,
	})
}

func (h *MongoDBHandler)DropIndex(cName string, keys ...string) error{
	c := h.se.DB(h.db).C(cName)
	return c.DropIndex(keys)
}

func (h *MongoDBHandler)Indexes(cName string)([]map[string]interface{}, error){
	c := h.se.DB(h.db).C(cName)
	indexes, err := c.Indexes()
	if err != nil{
		return nil, err
	}

	m := []map[string]interface{}{}
	for _, index := range indexes{
		mIdx := structs.Map(indexes)
		m = append(m, mIdx)
	}
	return m
}

func (h *MongoDBHandler)find(cName string, selector map[string]interface{}) (*mgo.Query){
	c := h.se.DB(h.db).C(cName)
	return c.Find(bson.M(selector))
}

func (h *MongoDBHandler)findByID(cName string, id string) (*mgo.Query){
	c := h.se.DB(h.db).C(cName)
	return c.FindId(bson.ObjectIdHex(id))
}

func (h *MongoDBHandler)FindAll(cName string, offset, limit int, sort map[string]interface{}, results interface{}) error{
	if len(sort) != 0{
		return h.find(cName, nil).Sort(sort).Skip(offset).Limit(limit).All(&results)
	}
	return h.find(cName, nil).Skip(offset).Limit(limit).All(&results)
}

func (h *MongoDBHandler)FindOne(cName string, offset, limit int, sort map[string]interface{}, result interface{}) error{
	if len(sort) != 0{
		return h.find(cName, nil).Sort(sort).Skip(offset).Limit(limit).One(&result)
	}
	return h.find(cName, nil).Skip(offset).Limit(limit).One(&result)
}

func (h *MongoDBHandler)FindByID(cName, id string, result interface{}) error{
	return h.findByID(cName, id).One(&result)
}

func (h *MongoDBHandler)Insert(cName string, cObjects...interface{}) error{
	c := h.se.DB(h.db).C(cName)
	return c.Insert(cObjects...)
}

func (h *MongoDBHandler)Upsert(cName string, selector map[string]interface{}, cObject interface{}) (string, error){
	c := h.se.DB(h.db).C(cName)
	info, err := c.Upsert(bson.M(selector), cObject)
	if err != nil{
		return "", err
	}
	return string(info.UpsertedId.(bson.ObjectId)), nil
}

func (h *MongoDBHandler)UpsertedId(cName, id string, cObject interface{})(string, error){
	c := h.se.DB(h.db).C(cName)
	info, err := c.UpsertId(bson.ObjectIdHex(id), cObject)
	if err != nil{
		return "", err
	}
	return string(info.UpsertedId.(bson.ObjectId)), nil
}

func (h *MongoDBHandler)Remove(cName string, selector map[string]interface{}) error{
	c := h.se.DB(h.db).C(cName)
	return c.Remove(bson.M(selector))
}

func (h *MongoDBHandler)RemoveAll(cName string, selector map[string]interface{}) (int, error){
	c := h.se.DB(h.db).C(cName)
	info, err := c.RemoveAll(bson.M(selector))
	if err != nil{
		return 0, err
	}
	return info.Removed(), err
}

func (h *MongoDBHandler)RemoveByID(cName, id string) error{
	c := h.se.DB(h.db).C(cName)
	return c.RemoveId(bson.ObjectIdHex(id))
}
