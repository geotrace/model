package model

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var ErrBadObjectId = errors.New("bad object id")

// Названия коллекций в хранилище.
var (
	CollectionUsers   = "users"
	CollectionEvents  = "events"
	CollectionPlaces  = "places"
	CollectionDevices = "devices"
)

// DB описывает хранилище данных и работу с ним.
type DB struct {
	name    string       // название базы данных
	session *mgo.Session // открытая сессия соединения с MongoDB
}
