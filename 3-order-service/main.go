package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/uid4oe/grpc-web-demo/3-order-service/client"
	"github.com/uid4oe/grpc-web-demo/3-order-service/orderdb"
	"github.com/uid4oe/grpc-web-demo/3-order-service/orderpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	orderpb.UnimplementedOrderServer
}

var (
	offer_client   client.OfferClient
	catalog_client client.CatalogClient
)

func (*server) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	log.Println("Order Service - Called CreateOrder")

	order_id := uuid.New().String()

	offer_details, err := offer_client.GetOfferDetails(&ctx, req.OfferId)
	if err != nil {
		return nil, error_response(err)
	}

	err = orderdb.CreateOneOrder(ctx, &orderdb.Order{
		Id: order_id, Offer_Id: offer_details.Id})

	if err != nil {
		return nil, error_response(err)
	}

	// call catalog client to update stats
	_, err = catalog_client.UpdateItem(&ctx, offer_details.CatalogId, offer_details.Amount)
	if err != nil {
		return nil, error_response(err)
	}

	// call offer client to update stats
	_, err = offer_client.UpdatePartnerTotalOrder(&ctx, offer_details.PartnerId)
	if err != nil {
		return nil, error_response(err)
	}

	return &orderpb.CreateOrderResponse{Id: order_id}, nil
}

func (*server) HandleOrderCompletion(ctx context.Context, req *orderpb.HandleOrderCompletionRequest) (*orderpb.HandleOrderCompletionResponse, error) {
	log.Println("Order Service - Called HandleOrderCompletion ID:", req.Id)

	order_details, err := orderdb.FindOneOrder(ctx, req.Id)
	if err != nil {
		return nil, error_response(err)
	}

	err = orderdb.UpdateOneOrder(context.Background(), &orderdb.Order{
		Id: order_details.Id, Rating: req.Rating})
	if err != nil {
		return nil, error_response(err)
	}

	offer_details, err := offer_client.GetOfferDetails(&ctx, order_details.Offer_Id)
	if err != nil {
		return nil, error_response(err)
	}

	res, err := offer_client.UpdatePartnerRating(&ctx, offer_details.PartnerId, req.Rating)
	if err != nil {
		return nil, error_response(err)
	}

	return &orderpb.HandleOrderCompletionResponse{Rating: res.Rating}, nil
}

func (s *server) GetOrderDetailStream(req *orderpb.GetOrderDetailsRequest,
	stream orderpb.Order_GetOrderDetailStreamServer) error {
	log.Println("Order Service - Started GetOrderDetailStream - Order ID:", req.Id)

	var wg sync.WaitGroup

	order, err := orderdb.FindOneOrder(context.Background(), req.Id)
	if err != nil {
		return error_response(err)
	}

	ctx := context.Background()

	offer_details, err := offer_client.GetOfferDetails(&ctx, order.Offer_Id)
	if err != nil {
		return error_response(err)
	}

	order_details, err := orderdb.FindAllDetailsByCatalogId(context.Background(), offer_details.CatalogId)
	if err != nil {
		return error_response(err)
	}

	for _, order_detail := range *order_details {
		wg.Add(1)
		go func(order_detail orderdb.Order_Detail) {
			defer wg.Done()
			time.Sleep(1000 * time.Millisecond)

			err = stream.Send(&orderpb.GetOrderDetailsResponse{Step: order_detail.Step, Detail: order_detail.Detail})
			if err != nil {
				stream.SendMsg(error_response(err))
				return
			}
			log.Printf("Order detail %d has been delivered", order_detail.Step)
		}(order_detail)
		wg.Wait()
	}
	stream.Context().Done()
	log.Println("Order Service - Completed GetOrderDetailStream - Order ID:", req.Id)
	return nil
}

func error_response(err error) error {
	log.Println("Order Service - ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func main() {
	log.Println("Running Order Service")

	lis, err := net.Listen("tcp", "0.0.0.0:55052")
	if err != nil {
		log.Println("Order Service - ERROR:", err.Error())
	}

	orderdb.Postgres_Client, err = orderdb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer orderdb.Postgres_Client.Close()

	s := grpc.NewServer()
	orderpb.RegisterOrderServer(s, &server{})

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("Order Service - ERROR:", err.Error())
	}

}
