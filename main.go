package main

import (
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo/bson/primitive"
)

type Model struct {
	ID          primitive.ObjectID `json:"_id"`
	DateCreated time.Time          `json:"date_created" bson:"date_created"`
	DateUpdated time.Time          `json:"date_updated" bson:"date_updated"`
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
