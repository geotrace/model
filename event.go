package model

import "gopkg.in/mgo.v2/bson"

type Events DB // для обращения к данным о событиях

// Get возвращает описание события с указанным идентификатором для конкретного
// устройства из хранилища.
func (db *Events) Get(groupId, deviceId, id string) (event *Event, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionEvents)
	if !bson.IsObjectIdHex(id) {
		err = ErrBadObjectId
		return
	}
	objID := bson.ObjectIdHex(id)
	event = new(Event)
	err = coll.Find(bson.M{"_id": objID, "groupId": groupId, "deviceId": deviceId}).
		Select(bson.M{"groupId": 0, "deviceId": 0}).One(event)
	session.Close()
	return
}

// List возвращает список всех событий, зарегистрированных для указанного
// устройства.
func (db *Events) List(groupID, deviceId string) (events []*Event, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionEvents)
	events = make([]*Event, 0)
	err = coll.Find(bson.M{"groupId": groupID, "deviceId": deviceId}).
		Select(bson.M{"groupId": 0, "deviceId": 0}).All(&events)
	session.Close()
	return
}

// Devices возвращает список идентификаторов устройств, данные о которых есть в
// коллекции событий для данной группы пользователей.
func (db *Events) Devices(groupID string) (deviceIds []string, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionEvents)
	deviceIds = make([]string, 0)
	err = coll.Find(bson.M{"groupId": groupID}).Distinct("deviceID", &deviceIds)
	session.Close()
	return
}

// Create добавляет в хранилище описание новых событий с привязкой к устройству.
func (db *Events) Create(groupId, deviceId string, events ...*Event) (err error) {
	objs := make([]interface{}, len(events))
	for i, event := range events {
		if !event.ID.Valid() {
			event.ID = bson.NewObjectId()
		}
		event.GroupID = groupId
		event.DeviceID = deviceId
		objs[i] = event
	}
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionEvents)
	err = coll.Insert(objs...)
	session.Close()
	return
}

// Update обновляет описание события в хранилище.
func (db *Events) Update(groupId, deviceId string, event *Event) (err error) {
	event.GroupID = groupId
	event.DeviceID = deviceId
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionEvents)
	err = coll.UpdateId(event.ID, event)
	session.Close()
	return
}

// Delete удаляет описание события из хранилища.
func (db *Events) Delete(groupId, deviceId, id string) (err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionEvents)
	err = coll.Remove(bson.M{"_id": id, "groupId": groupId, "deviceId": deviceId})
	session.Close()
	return
}
