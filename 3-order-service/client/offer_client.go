package client

import (
	"context"
	"errors"

	"github.com/uid4oe/grpc-web-demo/2-offer-service/offerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type OfferClient struct {
}

var (
	_                      = loadLocalEnv()
	offerGrpcService       = GetEnv("OFFER_GRPC_SERVICE")
	offerGrpcServiceClient offerpb.OfferClient
)

func prepareOfferGrpcClient(c *context.Context) error {

	conn, err := grpc.DialContext(*c, offerGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...)

	if err != nil {
		offerGrpcServiceClient = nil
		return errors.New("connection to offer gRPC service failed")
	}

	if offerGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	offerGrpcServiceClient = offerpb.NewOfferClient(conn)
	return nil
}

func (oc *OfferClient) GetOfferDetails(c *context.Context, id string) (*offerpb.GetOfferDetailsResponse, error) {

	if err := prepareOfferGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := offerGrpcServiceClient.GetOfferDetails(*c, &offerpb.GetOfferDetailsRequest{Id: id})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return res, nil
}

func (oc *OfferClient) UpdatePartnerTotalOrder(c *context.Context, id string) (*offerpb.UpdatePartnerTotalOrderResponse, error) {

	if err := prepareOfferGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := offerGrpcServiceClient.UpdatePartnerTotalOrder(*c, &offerpb.UpdatePartnerTotalOrderRequest{Id: id})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return res, nil
}

func (oc *OfferClient) UpdatePartnerRating(c *context.Context, id string, rating float32) (*offerpb.UpdatePartnerRatingResponse, error) {

	if err := prepareOfferGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := offerGrpcServiceClient.UpdatePartnerRating(*c, &offerpb.UpdatePartnerRatingRequest{Id: id, Rating: rating})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return res, nil
}
