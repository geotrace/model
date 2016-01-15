package model

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var (
	ErrBadObjectId = errors.New("bad object id")
	ErrNotFound    = mgo.ErrNotFound
)

// DB описывает хранилище данных и работу с ним.
type DB struct {
	session *mgo.Session // открытая сессия соединения с MongoDB
	name    string       // название базы данных
}

// InitDB инициализирует описание соединения с хранилищем и возвращает его.
func InitDB(session *mgo.Session, dbName string) *DB {
	return &DB{session, dbName}
}

// Close закрывает сессию соединения с MongoDB.
func (db *DB) Close() {
	db.session.Close()
}

// Названия коллекций в хранилище.
var (
	CollectionUsers   = "users"
	CollectionDevices = "devices"
	CollectionEvents  = "events"
	CollectionPlaces  = "places"
)
