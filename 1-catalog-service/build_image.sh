docker rmi uid4oe/alva-catalog -f
docker buildx build --platform=linux/amd64 . -t uid4oe/alva-catalog:latest
docker push uid4oe/alva-catalog:latest