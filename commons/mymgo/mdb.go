package mymgo

import (
	"fmt"
	"strings"

	"dudu/config"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Field struct {
	Id         string
	Collection string
}

type MdbSession struct {
	session *mgo.Session
	db      string
}

func (mdb *MdbSession) Session() *mgo.Session {
	return mdb.session.New()
}

var (
	AutoIncIdCollection = "amb_config"
	field               = &Field{
		Id:         "seq",
		Collection: "_id",
	}
)

func GetDb(cfgs []*config.DBConfig) *MdbSession {
	var url = generateUrl(cfgs)
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err.Error())
	}
	return &MdbSession{
		session: session,
		db:      cfgs[0].Db,
	}
}

func generateUrl(cfgs []*config.DBConfig) string {
	if len(cfgs) == 0 {
		return ""
	}

	mongoUrl := "mongodb://"
	if cfgs[0].User != "" && cfgs[0].Passwd != "" {
		mongoUrl += cfgs[0].User + ":" + cfgs[0].Passwd + "@"
	}

	addrs := make([]string, 0, len(cfgs))
	for _, addr := range cfgs {
		addrs = append(addrs, fmt.Sprintf("%s:%d", addr.Host, addr.Port))
	}

	mongoUrl += strings.Join(addrs, ",")

	if cfgs[0].Db != "" {
		mongoUrl += "/" + cfgs[0].Db
	}

	if cfgs[0].Options != "" {
		mongoUrl += "?" + cfgs[0].Options
	}
	return mongoUrl
}

func (mdb *MdbSession) AutoIncId(name string) (id int) {
	s := mdb.Session()
	id, err := autoIncr(s.DB(mdb.db).C(AutoIncIdCollection), name)
	s.Close()
	if err != nil {
		panic("Get next id of [" + name + "] fail:" + err.Error())
	}
	return
}

func autoIncr(c *mgo.Collection, name string) (id int, err error) {
	return incr(c, name, 1)
}

func incr(c *mgo.Collection, name string, step int) (id int, err error) {
	result := make(map[string]interface{})
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{field.Id: step}},
		Upsert:    true,
		ReturnNew: true,
	}
	_, err = c.Find(bson.M{field.Collection: name}).Apply(change, result)
	if err != nil {
		return
	}
	id, ok := result[field.Id].(int)
	if ok {
		return
	}
	id64, ok := result[field.Id].(int64)
	if !ok {
		err = fmt.Errorf("%s is ont int or int64", field.Id)
		return
	}
	id = int(id64)
	return
}
