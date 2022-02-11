docker rmi uid4oe/alva-order-envoy -f
docker build . -t uid4oe/alva-order-envoy:latest
docker push uid4oe/alva-order-envoy:latest