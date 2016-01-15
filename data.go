package model

//go:generate codecgen -u=true -o=codec.go data.go

import (
	"errors"
	"time"

	"github.com/geotrace/geo"
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
	Login string `bson:"_id" json:"id" codec:"id"`
	// уникальный идентификатор группы
	GroupID string `bson:"groupId,omitempty" json:"groupId,omitempty" codec:"groupId,omitempty"`
	// отображаемое имя
	Name string `bson:"name,omitempty" json:"name,omitempty" codec:"name,omitempty"`
	// хеш пароля пользователя
	Password Password `bson:"password" json:"-" codec:"-"`
}

// Device описывает информацию об устройстве.
//
// Каждое устройство имеет свой глобальный уникальный идентификатор, который не
// может повторяться. Плюс, устройство в каждый момент времени может быть
// привязано только к одной группе пользователей. Это позволяет устройству
// менять группу, блокируя доступ к старым данным, которые были собраны для
// другой группы.
//
// Устройству может быть назначен его тип. Это поле используется внутри сервиса
// для идентификации поддерживаемых устройством возможностей, формата данных и
// команд.
type Device struct {
	// глобальный уникальный идентификатор устройства
	ID string `bson:"_id" json:"id" codec:"id"`
	// уникальный идентификатор группы
	GroupID string `bson:"groupId,omitempty" json:"groupId,omitempty" codec:"groupId,omitempty"`
	// отображаемое имя
	Name string `bson:"name,omitempty" json:"name,omitempty" codec:"name,omitempty"`
	// идентификатор типа устройства
	Type string `bson:"type,omitempty" json:"type,omitempty" codec:"type,omitempty"`
	// хеш пароля для авторизации
	Password Password `bson:"password,omitempty" json:"-" codec:"-"`
}

// String возвращает строку с отображаемым именем устройства. Если для данного
// устройства определено имя, то возвращается именно оно. В противном случае
// возвращается уникальный идентификатор устройства.
func (d *Device) String() string {
	if d.Name != "" {
		return d.Name
	}
	return d.ID
}

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
	ID bson.ObjectId `bson:"_id" json:"id" codec:"id"`
	// уникальный идентификатор устройства
	DeviceID string `bson:"deviceId" json:"deviceId" codec:"deviceId"`
	// уникальный идентификатор группы
	GroupID string `bson:"groupId,omitempty" json:"groupId,omitempty" codec:"groupId,omitempty"`

	// временная метка
	Time time.Time `bson:"time" json:"time" codec:"time"`
	// тип события: Arrive, Leave, Travel, Check-in, Happen
	Type string `bson:"type,omitempty" json:"type,omitempty" codec:"type,omitempty"`
	// координаты точки
	Location *geo.Point `bson:"location,omitempty" json:"location,omitempty" codec:"location,omitempty"`
	// погрешность координат в метрах
	Accuracy float64 `bson:"accuracy,omitempty" json:"accuracy,omitempty" codec:"accuracy,omitempty"`
	// уровень заряда устройства на тот момент
	Power uint8 `bson:"power,omitempty" json:"power,omitempty" codec:"power,omitempty"`

	// иконка в виде эмодзи
	Emoji rune `bson:"emoji,omitempty" json:"emoji,omitempty" codec:"emoji,omitempty"`
	// текстовый комментарий к событию
	Comment string `bson:"comment,omitempty" json:"comment,omitempty" codec:"comment,omitempty"`
	// дополнительная именованная информация
	Data map[string]interface{} `bson:"data,omitempty,inline" json:"data,omitempty" json:"codec,omitempty"`
}

// Place описывает географическое место, задаваемое для группы пользователей.
// Такое место может быть описано либо в виде круга, задаваемого координатами
// центральной точки и радиусом в метрах, либо полигоном. Круг имеет более
// высокий приоритет, поэтому если задано и то, и другое, то используется именно
// описание круга.
//
// К сожалению, в формате GeoJSON, который использует для описание
// географических координат в MongoDB, нет возможности описать круг. Поэтому для
// работы с ним его приходится трансформировать его в некий многоугольник.
// Получившийся результат сохраняется в поле Geo и индексируется сервером баз
// данных. В том же случае, если задан полигон, то его описания просто
// копируется в это поле без каких-либо изменений.
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

// ErrBadPlaceData возвращается, если ни полигон, ни окружность не заданы в
// описании места.
var ErrBadPlaceData = errors.New("cyrcle or polygon is require in place")

// String возвращает строку с отображаемым именем описания места. Если для
// данного места задано имя, то возвращается именно оно. В противном случае
// возвращается его уникальный идентификатор.
func (p *Place) String() string {
	if p.Name != "" {
		return p.Name
	}
	return p.ID
}

// prepare осуществляет предварительную подготовку данных, создавая специальный
// объект для индекса.
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
