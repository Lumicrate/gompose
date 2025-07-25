package mongodb

import (
	"errors"
	"fmt"
	"github.com/Lumicrate/gompose/db"
	"github.com/Lumicrate/gompose/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type MongoAdapter struct {
	client   *mongo.Client
	database *mongo.Database
	ctx      context.Context

	uri    string
	dbName string
}

func New(uri string, dbName string) *MongoAdapter {
	return &MongoAdapter{
		uri:    uri,
		dbName: dbName,
	}
}

func (m *MongoAdapter) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.uri))
	if err != nil {
		return err
	}

	m.client = client
	m.database = client.Database(m.dbName)
	m.ctx = context.TODO()

	return nil
}

func (m *MongoAdapter) Migrate(entities []any) error {
	return nil
}

func (m *MongoAdapter) Create(entity any) error {
	collection := m.collectionFor(entity)

	_, err := collection.InsertOne(m.ctx, entity)
	return err
}

func (m *MongoAdapter) Update(entity any) error {
	collection := m.collectionFor(entity)

	idValue, err := getEntityID(entity)
	if err != nil {
		return err
	}

	// Exclude `_id` from the update document
	updateDoc, err := toBsonDWithoutID(entity)
	if err != nil {
		return err
	}

	elemType := getElemType(entity)
	typedID, err := getTypedId(idValue, elemType)

	filter := bson.M{"id": typedID}
	update := bson.M{"$set": updateDoc}

	res, err := collection.UpdateOne(m.ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no document found with id = %v", idValue)
	}
	return nil
}

func (m *MongoAdapter) Delete(id string, entity any) error {
	collection := m.collectionFor(entity)

	// Determine the correct ID type from the entity
	elemType := getElemType(entity)
	typedID, err := getTypedId(id, elemType)
	_, err = collection.DeleteOne(m.ctx, bson.M{"id": typedID})
	return err
}

func (m *MongoAdapter) FindAll(entity any, filters map[string]any, pagination db.Pagination, sort []db.Sort) (any, error) {
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	sliceType := reflect.SliceOf(entityType)
	slicePtr := reflect.New(sliceType).Interface()

	collection := m.collectionFor(entity)

	findOptions := options.Find()
	if pagination.Limit > 0 {
		findOptions.SetLimit(int64(pagination.Limit))
	}
	if pagination.Offset > 0 {
		findOptions.SetSkip(int64(pagination.Offset))
	}
	if len(sort) > 0 {
		sortDoc := bson.D{}
		for _, s := range sort {
			dir := 1
			if s.Direction == "desc" {
				dir = -1
			}
			sortDoc = append(sortDoc, bson.E{Key: s.Field, Value: dir})
		}
		findOptions.SetSort(sortDoc)
	}

	cursor, err := collection.Find(m.ctx, filters, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(m.ctx)

	if err := cursor.All(m.ctx, slicePtr); err != nil {
		return nil, err
	}

	return reflect.ValueOf(slicePtr).Elem().Interface(), nil
}

func (m *MongoAdapter) FindByID(id string, entity any) (any, error) {
	collection := m.collectionFor(entity)

	elemType := getElemType(entity)
	typedID, err := getTypedId(id, elemType)
	result := reflect.New(elemType).Interface()
	err = collection.FindOne(m.ctx, bson.M{"id": typedID}).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MongoAdapter) collectionFor(entity any) *mongo.Collection {
	t := reflect.TypeOf(entity)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return m.database.Collection(strings.ToLower(utils.Pluralize(t.Name())))
}

func getEntityID(entity any) (string, error) {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", errors.New("ID field not found")
	}

	if idField.Kind() == reflect.String {
		return idField.String(), nil
	}

	if idField.Kind() == reflect.Int || idField.Kind() == reflect.Int64 || idField.Kind() == reflect.Int32 {
		return fmt.Sprintf("%d", idField.Int()), nil
	}

	return "", errors.New("unsupported ID type")
}

func toBsonDWithoutID(entity any) (bson.D, error) {
	data, err := bson.Marshal(entity)
	if err != nil {
		return nil, err
	}
	var doc bson.D
	if err := bson.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	// Remove _id if present
	var cleaned bson.D
	for _, elem := range doc {
		if elem.Key != "_id" {
			cleaned = append(cleaned, elem)
		}
	}
	return cleaned, nil
}

func getElemType(entity any) reflect.Type {
	elemType := reflect.TypeOf(entity)
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	return elemType
}

func getTypedId(id string, elemType reflect.Type) (any, error) {
	idField, ok := elemType.FieldByName("ID")
	if !ok {
		return nil, fmt.Errorf("entity does not have an ID field")
	}

	// Convert string ID to correct type
	var typedID any
	switch idField.Type.Kind() {
	case reflect.Int, reflect.Int64:
		if intVal, err := strconv.Atoi(id); err == nil {
			typedID = intVal
		} else {
			return nil, fmt.Errorf("invalid int ID: %v", err)
		}
	case reflect.Uint, reflect.Uint64:
		if uintVal, err := strconv.ParseUint(id, 10, 64); err == nil {
			typedID = uintVal
		} else {
			return nil, fmt.Errorf("invalid uint ID: %v", err)
		}
	case reflect.String:
		typedID = id
	default:
		return nil, fmt.Errorf("unsupported ID type: %s", idField.Type.Kind())
	}

	return typedID, nil
}
