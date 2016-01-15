package model

import (
	"github.com/geotrace/uid"
	"gopkg.in/mgo.v2/bson"
)

type Places DB // для обращения к данным об описании мест

// Get возвращает описание места по его идентификатору. Кроме идентификатора
// места, который является уникальным, необходимо так же указывать идентификатор
// группы — это позволяет дополнительно ограничить даже случайный доступ
// пользователей к чужой информации.
func (db *Places) Get(groupId, id string) (place *Place, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionPlaces)
	place = new(Place)
	err = coll.Find(bson.M{"_id": id, "groupId": groupId}).
		Select(bson.M{"groupId": 0, "geo": 0}).One(place)
	session.Close()
	return
}

// List возвращает список всех мест, определенных в хранилище для данной группы
// пользователей.
func (db *Places) List(groupID string) (places []*Place, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionPlaces)
	places = make([]*Place, 0)
	err = coll.Find(bson.M{"groupId": groupID}).
		Select(bson.M{"groupId": 0, "geo": 0}).All(&places)
	session.Close()
	return
}

// Create добавляет в хранилище описание нового места для группы. Указание
// группы позволяет дополнительно защитить от ошибок переназначения места для
// другой группы.
func (db *Places) Create(groupId string, place *Place) (err error) {
	if err = place.prepare(); err != nil {
		return
	}
	if place.ID == "" {
		place.ID = uid.New()
	}
	place.GroupID = groupId
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionPlaces)
	err = coll.Insert(place)
	session.Close()
	return
}

// Update обновляет информацию о месте в хранилище. Указание группы позволяет
// дополнительно защитить от ошибок переназначения места для другой группы.
func (db *Places) Update(groupId string, place *Place) (err error) {
	if err = place.prepare(); err != nil {
		return
	}
	place.GroupID = groupId
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionPlaces)
	err = coll.UpdateId(place.ID, place)
	session.Close()
	return
}

// Delete удаляет описание места с указанным идентификатором из хранилища.
// Указание группы позволяет дополнительно защитить от ошибок доступа к чужой
// информации.
func (db *Places) Delete(groupId, id string) (err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionPlaces)
	err = coll.Remove(bson.M{"_id": id, "groupId": groupId})
	session.Close()
	return
}
