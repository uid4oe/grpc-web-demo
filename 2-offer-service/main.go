package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/uid4oe/grpc-web-demo/2-offer-service/offerdb"
	"github.com/uid4oe/grpc-web-demo/2-offer-service/offerpb"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	offerpb.UnimplementedOfferServer
}

func (s *server) GetOfferStream(req *offerpb.GetOfferStreamRequest,
	stream offerpb.Offer_GetOfferStreamServer) error {
	log.Println("Offer Service - Started GetOfferStream - Flow ID:", req.FlowId)

	// mock data insertion
	insert_random_offers(req.FlowId, req.CatalogId)

	var wg sync.WaitGroup

	offers, err := offerdb.FindAllOffersByFlowId(context.Background(), req.FlowId)
	if err != nil {
		return error_response(err)
	}

	for i, offer := range *offers {
		wg.Add(1)
		go func(offer offerdb.Offer) {
			defer wg.Done()
			time.Sleep(800 * time.Millisecond)

			partner_id, err := primitive.ObjectIDFromHex(offer.Partner_Id)
			if err != nil {
				stream.SendMsg(error_response(err))
				return
			}

			partner, err := offerdb.FindPartnerById(context.Background(), partner_id)
			if err != nil {
				stream.SendMsg(error_response(err))
				return
			}
			partner_info := &offerpb.PartnerInfo{Id: partner.Id.Hex(), Name: partner.Name, Rating: partner.Rating, Orders: partner.Orders}

			err = stream.Send(&offerpb.GetOfferStreamResponse{Id: offer.Id.Hex(), Partner: partner_info, Amount: offer.Amount})
			if err != nil {
				stream.SendMsg(error_response(err))
				return
			}
			log.Printf("Offer %d has been delivered", i+1)
		}(offer)
		wg.Wait()
	}

	log.Println("Offer Service - Completed GetOfferStream - Flow ID:", req.FlowId)
	return nil
}

func insert_random_offers(flow_id string, catalog_id string) {
	partners, err := offerdb.FindAllPartners(context.Background())
	if err != nil {
		log.Printf("Error at finding all partners - reason: %v", err.Error())
	}

	for _, partner := range *partners {
		offerdb.InsertOneOffer(context.Background(),
			&offerdb.Offer{Flow_Id: flow_id, Catalog_Id: catalog_id, Partner_Id: partner.Id.Hex(),
				Amount: 100 + rand.Int63n(300)})
	}
}

func (*server) GetOfferDetails(ctx context.Context, req *offerpb.GetOfferDetailsRequest) (*offerpb.GetOfferDetailsResponse, error) {
	log.Println("Offer Service - Called GetOfferDetails ID:", req.Id)

	uid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, error_response(err)
	}

	res, err := offerdb.FindOfferById(context.Background(), uid)

	if err != nil {
		return nil, error_response(err)
	}

	return &offerpb.GetOfferDetailsResponse{Id: res.Id.Hex(), CatalogId: res.Catalog_Id, PartnerId: res.Partner_Id, Amount: res.Amount}, nil
}

func (*server) UpdatePartnerRating(ctx context.Context, req *offerpb.UpdatePartnerRatingRequest) (*offerpb.UpdatePartnerRatingResponse, error) {
	log.Println("Offer Service - Called UpdatePartnerRating ID:", req.Id)

	uid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, error_response(err)
	}

	partner_from_db, err := offerdb.FindPartnerById(ctx, uid)
	if err != nil {
		return nil, error_response(err)
	}

	new_rating := (partner_from_db.Rating*float32(partner_from_db.Orders-1) + req.Rating) / float32(partner_from_db.Orders)

	err = offerdb.UpdatePartnerRating(context.Background(), &offerdb.Partner{
		Id: partner_from_db.Id, Rating: new_rating})
	if err != nil {
		return nil, error_response(err)
	}

	return &offerpb.UpdatePartnerRatingResponse{Rating: new_rating}, nil
}

func (*server) UpdatePartnerTotalOrder(ctx context.Context, req *offerpb.UpdatePartnerTotalOrderRequest) (*offerpb.UpdatePartnerTotalOrderResponse, error) {
	log.Println("Offer Service - Called UpdatePartnerTotalOrder ID:", req.Id)

	uid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, error_response(err)
	}

	err = offerdb.UpdatePartnerTotalOrder(context.Background(), uid)

	if err != nil {
		return nil, error_response(err)
	}

	return &offerpb.UpdatePartnerTotalOrderResponse{}, nil
}

func error_response(err error) error {
	log.Println("Offer Service - ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func main() {
	log.Println("Running Offer Service")

	lis, err := net.Listen("tcp", "0.0.0.0:55051")
	if err != nil {
		log.Println("Offer Service - ERROR:", err.Error())
	}

	offerdb.Mongo_Client, err = offerdb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer offerdb.Mongo_Client.Disconnect(context.Background())

	s := grpc.NewServer()
	offerpb.RegisterOfferServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("Offer Service - ERROR:", err.Error())
	}

}
