docker rmi uid4oe/alva-order -f
docker buildx build --platform=linux/amd64 . -t uid4oe/alva-order:latest
docker push uid4oe/alva-order:latest