package main

import (
	"context"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IDSetter interface {
	SetID(primitive.ObjectID)
}

type Model struct {
	DateCreated time.Time `json:"date_created" bson:"date_created"`
	DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
}

func (m *Todo) SetID(id primitive.ObjectID) {
	m.ID = id
}

type Todo struct {
	*Model
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Collection string             `model:"todo" bson:"-"`
	Title      string             `json:"title" bson:"title"`
	Completed  bool               `json:"completed" bson:"completed"`
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
	UpdateCreateDates(model)

	result, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		log.Fatal(err)
	}

	if m, ok := model.(IDSetter); ok {
		m.SetID(result.InsertedID.(primitive.ObjectID))
	}

	return model, nil
}

func UpdateCreateDates(model interface{}) {
	valueOf := reflect.ValueOf(model)

	if !isStruct(model) {
		return
	}

	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}

	// check if Model is part the instance definition
	modelField := valueOf.FieldByName("Model")
	if !modelField.IsValid() {
		return
	}

	// Model should be a struct
	modelValue := modelField.Interface()
	if !isStruct(modelValue) {
		return
	}

	concreteValue, ok := modelValue.(*Model)
	if !ok {
		return
	}

	concreteValue.DateCreated = time.Now()
	concreteValue.DateUpdated = time.Now()
}

func isStruct(value interface{}) bool {
	typeOf := reflect.TypeOf(value)

	switch typeOf.Kind() {
	case reflect.Struct:
		return true
	case reflect.Ptr:
		if typeOf.Elem().Kind() == reflect.Struct {
			return true
		}
	default:
		return false
	}

	return false
}

func List() []*Todo {
	todos := make([]*Todo, 0)
	collection, err := GetCollection(Todo{})

	if err != nil {
		log.Panic(err)
	}

	find := options.Find().SetShowRecordID(true)
	cursor, err := collection.Find(
		context.TODO(), bson.D{}, find)
	if err != nil {
		log.Panic(err)
	}

	cursor.All(context.TODO(), &todos)

	return todos
}
