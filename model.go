package main

import (
	"context"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IDSetter interface {
	SetID(primitive.ObjectID)
}

type Model struct {
	ID          primitive.ObjectID `json:"_id"`
	DateCreated time.Time          `json:"date_created" bson:"date_created"`
	DateUpdated time.Time          `json:"date_updated" bson:"date_updated"`
}

func (m *Model) SetID(id primitive.ObjectID) {
	m.ID = id
}

type Todo struct {
	Model
	Collection string `model:"todo"`
	Title      string `json:"title" bson:"title"`
	Completed  bool   `json:"completed" bson:"completed"`
}

func GetCollectionName(model interface{}, panicIfMissing bool) string {
	typeOf := reflect.TypeOf(model)
	// validate that this is struct or pointer to a struct
	switch typeOf.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		// If it's a pointer
		if typeOf.Elem().Kind() != reflect.Struct {
			log.Panic("Expected a struct type")
		}
		typeOf = typeOf.Elem()
	default:
		log.Panic("Expect a struct type")
	}

	field, found := typeOf.FieldByName("Collection")
	if !found && panicIfMissing {
		log.Panic("Missing field Collection")
	} else if !found {
		return ""
	}

	tag := field.Tag.Get("model")
	if len(tag) == 0 && panicIfMissing {
		log.Panic("The field `Collection` is missing the tag model")
	}

	return tag
}

func GetCollection(model interface{}) (*mongo.Collection, error) {
	db, err := GetDatabase()
	if err != nil {
		return nil, err
	}

	return db.Collection(GetCollectionName(model, true)), nil
}

func InsertOne(model interface{}) (interface{}, error) {
	collection, err := GetCollection(model)
	if err != nil {
		log.Fatal(err)
	}

	// update the date created and updated fields if it's of type model
	if model, ok := model.(*Model); ok {
		model.DateCreated = time.Now()
		model.DateUpdated = time.Now()
	}

	result, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		log.Fatal(err)
	}

	if m, ok := model.(IDSetter); ok {
		m.SetID(result.InsertedID.(primitive.ObjectID))
	}

	return model, nil
}
