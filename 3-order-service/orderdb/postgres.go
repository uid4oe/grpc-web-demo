package orderdb

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type Order struct {
	Id       string
	Offer_Id string
	Rating   float32
}

type Order_Detail struct {
	Step       int32
	Catalog_ID string
	Detail     string
}

var _ = loadLocalEnv()
var (
	db       = GetEnv("POSTGRES_DB")
	username = GetEnv("POSTGRES_USER")
	password = GetEnv("POSTGRES_PASSWORD")
	host     = GetEnv("POSTGRES_HOST")
)

var Postgres_Client *pgxpool.Pool

func NewClient(ctx context.Context) (*pgxpool.Pool, error) {
	url := "postgres://" + username + ":" + password + "@" + host + "/" + db
	client, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, errors.New("cannot connect to postgres instance")
	}
	return client, nil
}

func FindOneOrder(ctx context.Context, id string) (*Order, error) {
	order := Order{}
	err := Postgres_Client.QueryRow(ctx, "select * from orders where id=$1", id).Scan(&order.Id, &order.Offer_Id, &order.Rating)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func CreateOneOrder(ctx context.Context, order *Order) error {
	_, err := Postgres_Client.Exec(ctx, "insert into orders(id,offer_id,rating) values($1,$2,0.0)", order.Id, order.Offer_Id)
	return err
}

func UpdateOneOrder(ctx context.Context, order *Order) error {
	_, err := Postgres_Client.Exec(ctx, "update orders set rating=$1 where id=$2", order.Rating, order.Id)
	return err
}

func FindAllDetailsByCatalogId(ctx context.Context, order_type string) (*[]Order_Detail, error) {

	rows, err := Postgres_Client.Query(ctx, "select * from order_details where catalog_id=$1 order by step", order_type)

	var order_detail_list []Order_Detail

	for rows.Next() {
		var od Order_Detail
		err = rows.Scan(&od.Step, &od.Catalog_ID, &od.Detail)
		order_detail_list = append(order_detail_list, od)
	}

	if err != nil {
		return nil, err
	}

	return &order_detail_list, nil
}

func loadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("ORDER_GRPC_SERVICE"); !runningInContainer {
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
