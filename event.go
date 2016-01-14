package model

import (
	"time"

	"github.com/geotrace/geo"
	"gopkg.in/mgo.v2/bson"
)

// Event обычно описывает место, время и событие, которое в нем случилось.
//
// Каждое событие получает свой уникальный идентификатор, назначаемый
// непосредственно системой. Кроме того, событие привязано к конкретному
// идентификатору устройства и группе пользователей. Группа пользователей здесь
// представлена отдельным свойством, не смотря на то, что ее можно достаточно
// легко получить и из связи с устройством. Это сделано намеренно, чтобы в тех
// случаях, когда устройство меняет владельца (группу), старые данные о событиях
// не становились автоматически доступны новым пользователям.
//
// Каждое событие в обязательном порядке характеризуется временем, когда оно
// произошло. Если при создании описания события время было опущено, то будет
// автоматически добавлено текущее время сервера.
//
// Тип события задает один из предопределенных типов события. Если не указано,
// то считается, что тип события не определен.
//
// Каждое событие обычно характеризуется координатами географической точки, в
// которой оно случилось и дополнительным параметром, указывающим возможный
// радиус погрешности вычисления данной точки.
//
// Дополнительно, каждое событие может иметь свое описание в текстовом виде и
// иконку, характеризующую его в некотором визуальном виде. Но с последним
// обычно тяжело: кто и сколько таких иконок нарисует? Поэтому было принято
// решения вместо иконки использовать пиктограмму из стандартного набора эмодзи.
//
// И, наконец, последний элемент: именованные поля с произвольным содержимым,
// позволяющим описать любую дополнительную информацию. В частности, думаю,
// значения датчиков и сенсоров хорошо и удобно сохранять именно в таком виде.
// Плюс, всегда можно добавить что-то дополнительно практически в любом удобном
// формате. Главное, чтобы приложение знало, что потом с этим делать.
type Event struct {
	// уникальный идентификатор записи
	ID bson.ObjectId `bson:"_id" json:"id"`
	// уникальный идентификатор устройства
	DeviceID string `bson:"deviceId" json:"deviceId"`
	// уникальный идентификатор группы
	GroupID string `bson:"groupId,omitempty" json:"groupId,omitempty"`

	// временная метка
	Time time.Time `bson:"time" json:"time"`
	// тип события: Arrive, Leave, Travel, Check-in, Happen
	Type string `bson:"type,omitempty" json:"type,omitempty"`
	// координаты точки
	Location *geo.Point `bson:"location,omitempty" json:"location,omitempty"`
	// погрешность координат в метрах
	Accuracy float64 `bson:"accuracy,omitempty" json:"accuracy,omitempty"`
	// уровень заряда устройства на тот момент
	Power uint8 `bson:"power,omitempty" json:"power,omitempty"`

	// иконка в виде эмодзи
	Emoji rune `bson:"emoji,omitempty" json:"emoji,omitempty"`
	// текстовый комментарий к событию
	Comment string `bson:"comment,omitempty" json:"comment,omitempty"`
	// дополнительная именованная информация
	Data map[string]interface{} `bson:"data,omitempty,inline" json:"data,omitempty"`
}

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
