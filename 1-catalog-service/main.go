package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/uid4oe/grpc-web-demo/1-catalog-service/catalogdb"
	"github.com/uid4oe/grpc-web-demo/1-catalog-service/catalogpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	catalogpb.UnimplementedCatalogServer
}

var (
	timeout = 2 * time.Second
)

func (*server) GetItem(ctx context.Context, req *catalogpb.GetItemRequest) (*catalogpb.GetItemResponse, error) {
	log.Println("Catalog Service - Called GetItem - ID:", req.Id)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	d, err := catalogdb.FindOne(c, req.Id)
	if err != nil {
		return nil, error_response(err)
	}

	return &catalogpb.GetItemResponse{Item: &catalogpb.CatalogItem{Id: d.Id, Name: d.Name, Details: d.Details,
		Availability: randomAvailability(), TotalOrders: d.Total_Orders, AverageCost: d.Average_Cost}}, nil
}

func (*server) UpdateItem(ctx context.Context, req *catalogpb.UpdateItemRequest) (*catalogpb.UpdateItemResponse, error) {
	log.Println("Catalog Service - Called UpdateItem - ID:", req.Id)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	catalog_item_from_db, err := catalogdb.FindOne(c, req.Id)
	if err != nil {
		return nil, error_response(err)
	}

	new_total_orders := catalog_item_from_db.Total_Orders + 1
	new_average_cost := (catalog_item_from_db.Average_Cost*catalog_item_from_db.Total_Orders + req.LatestCost) / new_total_orders

	err = catalogdb.UpdateOne(c, &catalogdb.Catalog{Id: req.Id, Total_Orders: new_total_orders, Average_Cost: new_average_cost})
	if err != nil {
		return nil, error_response(err)
	}

	return &catalogpb.UpdateItemResponse{}, nil
}

func (*server) GetItems(ctx context.Context, req *catalogpb.GetItemsRequest) (*catalogpb.GetItemsResponse, error) {
	log.Println("Catalog Service - Called GetItems")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := catalogdb.FindAll(c)
	if err != nil {
		return nil, error_response(err)
	}

	var res catalogpb.GetItemsResponse
	for _, d := range *data {
		res.Item = append(res.Item, &catalogpb.CatalogItem{Id: d.Id, Name: d.Name, Details: d.Details, Availability: randomAvailability(), TotalOrders: d.Total_Orders, AverageCost: d.Average_Cost})
	}

	return &res, nil
}

func randomAvailability() string {
	if rand.New(rand.NewSource(rand.Int63())).Float32() > 0.1 {
		return "online"
	}
	return "offline"
}

func error_response(err error) error {
	log.Println("Catalog Service - ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func main() {
	log.Println("Running Catalog Service")

	lis, err := net.Listen("tcp", "0.0.0.0:55050")
	if err != nil {
		log.Println("Catalog Service - ERROR:", err.Error())
	}

	catalogdb.Postgres_Client, err = catalogdb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer catalogdb.Postgres_Client.Close()

	s := grpc.NewServer()
	catalogpb.RegisterCatalogServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("Catalog Service - ERROR:", err.Error())
	}

}
