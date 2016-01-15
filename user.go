package model

import (
	"github.com/geotrace/uid"
	"gopkg.in/mgo.v2/bson"
)

type Users DB // для обращения к данным о зарегистрированных пользователях

// Login возвращает информацию о пользователе по его логину.
func (db *Users) Login(userID string) (user *User, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionUsers)
	err = coll.FindId(userID).One(&user)
	session.Close()
	return
}

// List возвращает список всех пользователей, зарегистрированных в указанной
// группе.
func (db *Users) List(groupID string) (users []User, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionUsers)
	users = make([]User, 0)
	err = coll.Find(bson.M{"groupId": groupID}).
		Select(bson.M{"password": 0, "groupId": 0}).All(&users)
	session.Close()
	return
}

// Create создает нового пользователя по его описанию. Поле Login должно быть
// уникальным, в противном случае возвращается ошибка.
func (db *Users) Create(user *User) (err error) {
	if user.Login == "" {
		user.Login = uid.New()
	}
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionUsers)
	err = coll.Insert(user)
	session.Close()
	return
}

// Update обновляет информацию о пользователе в хранилище.
func (db *Users) Update(user User) (err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionUsers)
	err = coll.UpdateId(user.Login, user)
	session.Close()
	return
}

// Delete удаляет пользователя с указанным логином из хранилища.
func (db *Users) Delete(login string) (err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionUsers)
	err = coll.RemoveId(login)
	session.Close()
	return
}
