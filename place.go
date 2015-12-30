package model

import (
	"errors"

	"github.com/geotrace/geo"
	"github.com/geotrace/uid"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrBadPlaceData = errors.New("bad place data: cyrcle or polygon is require")
)

// Place описывает географическое место, задаваемое для группы пользователей. Такое место может
// быть описано либо в виде круга, задаваемого координатами центральной точки и радиусом в метрах,
// либо полигоном. Круг имеет более высокий приоритет, поэтому если задано и то, и другое, то
// используется именно описание круга.
//
// К сожалению, в формате GeoJSON, который использует для описание географических координат в
// MongoDB, нет возможности описать круг. Поэтому для работы с ним его приходится трансформировать
// его в некий многоугольник. Получившийся результат сохраняется в поле Geo и индексируется
// сервером баз данных. В том же случае, если задан полигон, то его описания просто копируется в
// это поле без каких-либо изменений.
type Place struct {
	// уникальный идентификатор описания места
	ID string `bson:"_id,omitempty" json:"id"`
	// уникальный идентификатор группы
	GroupID string `bson:"groupId,omitempty" json:"groupId,omitempty"`
	// отображаемое имя
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	// географическое описание места как круга
	Circle *geo.Circle `bson:"circle,omitempty" json:"circle,omitempty"`
	// географическое описание места в виде полигона
	Polygon *geo.Polygon `bson:"polygon,omitempty" json:"polygon,omitempty"`
	// описание в формате GeoJSON для поиска
	Geo interface{} `bson:"geo" json:"-"`
}

// String возвращает строку с отображаемым именем описания места. Если для данного места задано
// имя, то возвращается именно оно. В противном случае возвращается его уникальный идентификатор.
func (p *Place) String() string {
	if p.Name != "" {
		return p.Name
	}
	return p.ID
}

// prepare осуществляет предварительную подготовку данных, создавая специальный объект для индекса.
func (p *Place) prepare() (err error) {
	// анализируем описание места и формируем данные для индексации
	if p.Circle != nil {
		p.Polygon = nil
		p.Geo = p.Circle.Geo()
	} else if p.Polygon != nil {
		p.Circle = nil
		p.Geo = p.Polygon.Geo()
	} else {
		err = ErrBadPlaceData
	}
	return
}

// PlaceGet возвращает описание места по его идентификатору. Кроме идентификатора места, который
// является уникальным, необходимо так же указывать идентификатор группы — это позволяет
// дополнительно ограничить даже случайный доступ пользователей к чужой информации.
func (db *DB) PlaceGet(groupId, id string) (place *Place, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionPlaces)
	place = new(Place)
	err = coll.Find(bson.M{"_id": id, "groupId": groupId}).Select(bson.M{"groupId": 0, "geo": 0}).One(place)
	session.Close()
	return
}

// PlaceCreate добавляет в хранилище описание нового места для группы. Указание группы позволяет
// дополнительно защитить от ошибок переназначения места для другой группы.
func (db *DB) PlaceCreate(groupId string, place *Place) (err error) {
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

// PlaceUpdate обновляет информацию о месте в хранилище. Указание группы позволяет
// дополнительно защитить от ошибок переназначения места для другой группы.
func (db *DB) PlaceUpdate(groupId string, place *Place) (err error) {
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

// PlaceDelete удаляет описание места с указанным идентификатором из хранилища. Указание группы
// позволяет дополнительно защитить от ошибок доступа к чужой информации.
func (db *DB) PlaceDelete(groupId, id string) (err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionPlaces)
	err = coll.Remove(bson.M{"_id": id, "groupId": groupId})
	session.Close()
	return
}

// PlaceList возвращает список всех мест, определенных в хранилище для данной группы пользователей.
func (db *DB) PlaceList(groupID string) (places []*Place, err error) {
	session := db.session.Copy()
	coll := session.DB(db.name).C(CollectionPlaces)
	places = make([]*Place, 0)
	err = coll.Find(bson.M{"groupId": groupID}).Select(bson.M{"groupId": 0, "geo": 0}).All(&places)
	session.Close()
	return
}
