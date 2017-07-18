package internal

import "gopkg.in/mgo.v2/bson"

type UserData struct {
	Id bson.ObjectId `"_id"`
	Name string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

func (u *UserData)register()error  {
	db := mongoDB.Ref()
	defer mongoDB.UnRef(db)

	u.Id = bson.NewObjectId()
	err := db.DB(DBName).C(myusers).Insert(u)
	return err
}
