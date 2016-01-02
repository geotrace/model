package model

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var ErrBadObjectId = errors.New("bad object id")

// DB описывает хранилище данных и работу с ним.
type DB struct {
	session *mgo.Session // открытая сессия соединения с MongoDB
	name    string       // название базы данных
}

// InitDB инициализирует описание соединения с хранилищем и возвращает его.
func InitDB(session *mgo.Session, dbName string) *DB {
	return &DB{session, dbName}
}

// Специализированные объекты для доступа к разным типам данных в хранилище.
type (
	DBUsers   DB // для обращения к данным о зарегистрированных пользователях
	DBDevices DB // для обращения к данным об устройствах
	DBEvents  DB // для обращения к данным о событиях
	DBPlaces  DB // для обращения к данным об описании мест
)

// Названия коллекций в хранилище.
var (
	CollectionUsers   = "users"
	CollectionDevices = "devices"
	CollectionEvents  = "events"
	CollectionPlaces  = "places"
)
