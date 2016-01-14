package model

import (
	"github.com/geotrace/uid"
	"gopkg.in/mgo.v2/bson"
)

// User описывает информацию о пользователе.
//
// Логин пользователя является глобальным уникальным идентификатором
// пользователя и не может повторяться для разных пользователей. Поэтому, скорее
// всего, удобнее использовать в качестве такого идентификатора e-mail, что
// избавит от головной боли с уникальностью. Или любой другой идентификатор,
// который будет действительно глобально уникальным.
//
// Пользователи объединяются в группы, которые разделяют общие ресурсы: имеют
// доступ к трекам устройств той же группы, общие описания мест и так далее.
// Пользователь может состоять только в одной группе, но может ее сменить.
// Идентификатор группы генерируется непосредственно сервером.
//
// Пароль пользователя не хранится в системе, а вместо этого хранится хеш от
// него: этого вполне достаточно, чтобы иметь возможность проверить правильность
// введенного пароля, но не позволит его восстановить в исходном виде. В
// качестве алгоритма хеширования выбран bcrypt (Provos and Mazières's bcrypt
// adaptive hashing algorithm).
type User struct {
	// логин пользователя
	Login string `bson:"_id" json:"id"`
	// уникальный идентификатор группы
	GroupID string `bson:"groupId,omitempty" json:"groupId,omitempty"`
	// отображаемое имя
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	// хеш пароля пользователя
	Password Password `bson:"password" json:"-"`
}

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
