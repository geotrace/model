package model

import (
	"github.com/geotrace/uid"
	"gopkg.in/mgo.v2/bson"
)

// Device описывает информацию об устройстве.
//
// Каждое устройство имеет свой глобальный уникальный идентификатор, который не может повторяться.
// Плюс, устройство в каждый момент времени может быть привязано только к одной группе
// пользователей. Это позволяет устройству менять группу, блокируя доступ к старым данным, которые
// были собраны для другой группы.
//
// Устройству может быть назначен его тип. Это поле используется внутри сервиса для идентификации
// поддерживаемых устройством возможностей, формата данных и команд.
type Device struct {
	// глобальный уникальный идентификатор устройства
	ID string `bson:"_id" json:"id"`
	// уникальный идентификатор группы
	GroupID string `bson:"groupId,omitempty" json:"groupId,omitempty"`
	// отображаемое имя
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	// идентификатор типа устройства
	Type string `bson:"type,omitempty" json:"type,omitempty"`
	// хеш пароля для авторизации
	Password Password `bson:"password,omitempty" json:"-"`
}

// String возвращает строку с отображаемым именем устройства. Если для данного устройства
// определено имя, то возвращается именно оно. В противном случае возвращается уникальный
// идентификатор устройства.
func (d *Device) String() string {
	if d.Name != "" {
		return d.Name
	}
	return d.ID
}

// Get возвращает информацию о устройстве с указанным идентификатором, которое привязано
// к указанной группе.
func (db *DBDevices) Get(groupId, id string) (device *Device, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionDevices)
	device = new(Device)
	err = coll.Find(bson.M{"_id": id, "groupId": groupId}).
		Select(bson.M{"groupId": 0, "password": 0}).One(device)
	session.Close()
	return
}

// List возвращает список всех устройств, которые зарегистрированы для данной группы
// пользователей.
func (db *DBDevices) List(groupID string) (devices []*Device, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionDevices)
	devices = make([]*Device, 0)
	err = coll.Find(bson.M{"groupId": groupID}).
		Select(bson.M{"groupId": 0, "password": 0}).All(&devices)
	session.Close()
	return
}

// Create создает описание нового устройства, одновременно привязывая его к указанной группе.
func (db *DBDevices) Create(groupId string, device *Device) (err error) {
	if device.ID == "" {
		device.ID = uid.New()
	}
	device.GroupID = groupId
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionDevices)
	err = coll.Insert(device)
	session.Close()
	return
}

// Update обновляет описание устройства и привязывает его к указанной группе.
func (db *DBDevices) Update(groupId string, device *Device) (err error) {
	device.GroupID = groupId
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionDevices)
	err = coll.UpdateId(device.ID, device)
	session.Close()
	return
}

// Delete удаляет описание устройства.
func (db *DBDevices) Delete(groupId, id string) (err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionDevices)
	err = coll.Remove(bson.M{"_id": id, "groupId": groupId})
	session.Close()
	return
}
