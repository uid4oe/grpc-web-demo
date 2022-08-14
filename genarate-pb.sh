cd 1-catalog-service/catalogpb/
protoc --go_out=. --go-grpc_out=. catalog.proto
cd ..
cd ..
cd 2-offer-service/offerpb/
protoc --go_out=. --go-grpc_out=. offer.proto
cd ..
cd .. 
cd 3-order-service/orderpb/
protoc --go_out=. --go-grpc_out=. order.proto