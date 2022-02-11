docker rmi uid4oe/alva-catalog-envoy -f
docker build . -t uid4oe/alva-catalog-envoy:latest
docker push uid4oe/alva-catalog-envoy:latest