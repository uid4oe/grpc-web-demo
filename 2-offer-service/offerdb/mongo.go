package offerdb

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Offer struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Flow_Id    string             `bson:"flow_id,omitempty"`
	Catalog_Id string             `bson:"catalog_id,omitempty"`
	Partner_Id string             `bson:"partner_id,omitempty"`
	Amount     int64              `bson:"amount,omitempty"`
}

type Partner struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name,omitempty"`
	Rating float32            `bson:"rating,omitempty"`
	Orders int64              `bson:"orders,omitempty"`
}

var _ = loadLocalEnv()
var (
	db           = GetEnv("MONGO_INITDB_DATABASE")
	user         = GetEnv("MONGO_INITDB_USER")
	pwd          = GetEnv("MONGO_INITDB_PWD")
	partner_coll = GetEnv("PARTNER_MONGO_COLLECTION")
	offer_coll   = GetEnv("OFFER_MONGO_COLLECTION")
	addr         = GetEnv("MONGO_CONN")
)

var Mongo_Client *mongo.Client

func NewClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(addr).
			SetAuth(options.Credential{
				AuthSource: db,
				Username:   user,
				Password:   pwd,
			}))
	if err != nil {
		return nil, errors.New("invalid mongodb options")
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.New("cannot connect to mongodb instance")
	}
	return client, nil
}

func UpdatePartnerRating(ctx context.Context, partner *Partner) error {

	collection := Mongo_Client.Database(db).Collection(partner_coll)

	_, err := FindPartnerById(ctx, partner.Id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": partner.Id}
	update := bson.M{"$set": bson.M{"rating": partner.Rating}}
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func UpdatePartnerTotalOrder(ctx context.Context, id primitive.ObjectID) error {

	collection := Mongo_Client.Database(db).Collection(partner_coll)

	partner_from_db, err := FindPartnerById(ctx, id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"orders": partner_from_db.Orders + 1}}
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

func FindOfferById(ctx context.Context, id primitive.ObjectID) (*Offer, error) {

	collection := Mongo_Client.Database(db).Collection(offer_coll)

	var data Offer
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func FindPartnerById(ctx context.Context, id primitive.ObjectID) (*Partner, error) {

	collection := Mongo_Client.Database(db).Collection(partner_coll)

	var data Partner
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func FindAllPartners(ctx context.Context) (*[]Partner, error) {

	collection := Mongo_Client.Database(db).Collection(partner_coll)

	var data []Partner
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var partner Partner
		cursor.Decode(&partner)
		data = append(data, partner)
	}
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func InsertOneOffer(ctx context.Context, offer *Offer) error {

	collection := Mongo_Client.Database(db).Collection(offer_coll)

	offer.Id = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, offer)
	return err
}

func FindAllOffersByFlowId(ctx context.Context, flow_id string) (*[]Offer, error) {

	collection := Mongo_Client.Database(db).Collection(offer_coll)

	var data []Offer
	cursor, err := collection.Find(ctx, bson.M{"flow_id": flow_id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var offer Offer
		cursor.Decode(&offer)
		data = append(data, offer)
	}
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func loadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("MONGO_CONN"); !runningInContainer {
		err := godotenv.Load("../.env.local")
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal("Environment variable not found: ", key)
	}
	return value
}
