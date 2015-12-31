package model

import (
	"testing"

	"gopkg.in/mgo.v2"
)

func TestDBType(t *testing.T) {
	mdi, err := mgo.ParseURL("mongodb://localhost/geotrace")
	if err != nil {
		t.Fatal(err)
	}
	session, err := mgo.DialWithInfo(mdi)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	db := &DB{session, mdi.Database}
	users := (*DBUsers)(db)
	// users.List("groupID")
	// pretty.Println(db)
	// pretty.Println(users)
}
