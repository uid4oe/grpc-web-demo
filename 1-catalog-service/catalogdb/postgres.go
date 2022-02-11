package catalogdb

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type Catalog struct {
	Id           string
	Name         string
	Details      string
	Total_Orders int64
	Average_Cost int64
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

func FindOne(ctx context.Context, id string) (*Catalog, error) {
	catalog := Catalog{}
	err := Postgres_Client.QueryRow(ctx, "select * from catalog where id=$1", id).Scan(&catalog.Id, &catalog.Name, &catalog.Details, &catalog.Total_Orders, &catalog.Average_Cost)
	if err != nil {
		return nil, err
	}
	return &catalog, nil
}

func UpdateOne(ctx context.Context, catalog *Catalog) error {
	_, err := Postgres_Client.Exec(ctx, "update catalog set total_orders=$1, average_cost=$2 where id=$3", catalog.Total_Orders, catalog.Average_Cost, catalog.Id)
	return err
}

func FindAll(ctx context.Context) (*[]Catalog, error) {

	rows, err := Postgres_Client.Query(ctx, "select id,name,details,total_orders,average_cost from catalog order by name")

	var catalog_list []Catalog

	for rows.Next() {
		var c Catalog
		err = rows.Scan(&c.Id, &c.Name, &c.Details, &c.Total_Orders, &c.Average_Cost)
		catalog_list = append(catalog_list, c)
	}

	if err != nil {
		return nil, err
	}

	return &catalog_list, nil
}

func loadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("CATALOG_GRPC_SERVICE"); !runningInContainer {
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
