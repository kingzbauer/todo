package main

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	IDSetter interface {
		SetID(primitive.ObjectID)
	}

	FilterID interface {
		FilterID() bson.D
	}

	Model struct {
		DateCreated time.Time `json:"date_created" bson:"date_created"`
		DateUpdated time.Time `json:"date_updated" bson:"date_updated"`
	}

	Todo struct {
		*Model
		ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		Collection string             `model:"todo" bson:"-" json:"-"`
		Title      string             `json:"title" bson:"title" validate:"required"`
		Completed  bool               `json:"completed" bson:"completed"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (c *CustomValidator) Validate(i interface{}) error {
	return c.validator.Struct(i)
}

func (m *Todo) SetID(id primitive.ObjectID) {
	m.ID = id
}

func (t *Todo) FilterID() bson.D {
	return bson.D{{"_id", t.ID}}
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

	if modelField.IsNil() {
		modelField.Set(reflect.New(reflect.TypeOf(Model{})))
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

	location, _ := time.LoadLocation("Africa/Nairobi")

	concreteValue.DateCreated = time.Now().In(location)
	concreteValue.DateUpdated = time.Now().In(location)
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

func GetID(model interface{}) primitive.ObjectID {
	if !isStruct(model) {
		log.Fatal("Requires a struct")
	}

	valueOf := reflect.ValueOf(model)
	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}

	idField := valueOf.FieldByName("ID")

	return idField.Interface().(primitive.ObjectID)
}

func Update(model *Todo) error {
	collection, err := GetCollection(model)
	if err != nil {
		log.Fatal(err)
	}

	update := bson.D{{"$set", bson.D{{"title", model.Title}, {"completed", model.Completed}}}}
	_, err = collection.UpdateOne(context.TODO(), model.FilterID(), update)
	return err
}

func Delete(model FilterID) error {
	collection, err := GetCollection(model)
	if err != nil {
		log.Fatal(err)
	}

	_, err = collection.DeleteOne(context.TODO(), model.FilterID())
	return err
}

func Retrieve(filter bson.D, typ interface{}) (interface{}, error) {
	typeOf := reflect.TypeOf(typ)
	if typeOf.Kind() != reflect.Struct {
		log.Panic("Requires a struct")
	}

	col, err := GetCollection(typ)
	if err != nil {
		log.Fatal(err)
	}

	result := col.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}

	value := reflect.New(typeOf).Interface()
	result.Decode(value)

	return value, nil
}
