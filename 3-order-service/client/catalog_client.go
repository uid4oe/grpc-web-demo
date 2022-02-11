package client

import (
	"context"
	"errors"

	"github.com/uid4oe/grpc-web-demo/1-catalog-service/catalogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type CatalogClient struct {
}

var (
	_                        = loadLocalEnv()
	catalogGrpcService       = GetEnv("CATALOG_GRPC_SERVICE")
	catalogGrpcServiceClient catalogpb.CatalogClient
)

func prepareCatalogGrpcClient(c *context.Context) error {

	conn, err := grpc.DialContext(*c, catalogGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...)

	if err != nil {
		catalogGrpcServiceClient = nil
		return errors.New("connection to offer gRPC service failed")
	}

	if catalogGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	catalogGrpcServiceClient = catalogpb.NewCatalogClient(conn)
	return nil
}

func (oc *CatalogClient) UpdateItem(c *context.Context, id string, latest_cost int64) (*catalogpb.UpdateItemResponse, error) {

	if err := prepareCatalogGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := catalogGrpcServiceClient.UpdateItem(*c, &catalogpb.UpdateItemRequest{Id: id, LatestCost: latest_cost})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return res, nil
}
