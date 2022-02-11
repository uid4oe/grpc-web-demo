docker rmi uid4oe/alva-offer -f
docker buildx build --platform=linux/amd64 . -t uid4oe/alva-offer:latest
docker push uid4oe/alva-offer:latest