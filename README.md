# gRPC-Web Demo

<p align="center"> 
<img src="./diagram.png" width="500" />
</p>

[Deployed, Check it out!](https://alva-grpc.web.app/)

Please read the [Medium Story](https://medium.com/alva-labs/building-microapps-with-grpc-web-64b7cdf50313) first.

This repo contains microservices written in Go (gRPC Servers) and their envoy proxy configurations for establishing connection through gRPC-Web with the client (UI). gRPC Unary and gRPC Server-side Streaming examples are available in the services.

|            | gRPC Unary | gRPC Server-side Streaming |
| :--------: | :--------: | :------------------------: |
| Catalog μS |     ✅     |                            |
|  Offer μS  |            |             ✅             |
|  Order μS  |     ✅     |             ✅             |

![](./ui.gif)

## Installation

Let's start!

Please make sure Docker is running.

First, clone the repository

```bash
git clone https://github.com/uid4oe/grpc-web-demo.git
```

Then, open a terminal window from the root of the project.
We will just need to create network & execute the `docker-compose` command

```bash
docker network create grpc-web-demo-alva-net
docker-compose up -d
```

At this point everything should be up and running!

You can track service logs from containers: `catalog`, `offer`, `order`.

You can access to UI at

```bash
http://localhost:3000
```

Additionally, you can check [gRPC-Web Demo UI](https://github.com/uid4oe/grpc-web-demo-ui) which has the UI code for this demo.

## Local Development

For running services in your local environment, we will need MongoDB, PostgreSQL, and envoy proxies (gRPC-Web related). We can use `local.yml` file for setting up MongoDB, PostgreSQL, Catalog Envoy, Offer Envoy, Order Envoy and UI.

Open a terminal in the project root folder, then;

```bash
docker network create grpc-web-demo-alva-net
docker compose --file local.yml up -d
```

Now let's start microservices locally.

```bash
sh start-local.sh
```

That's great. Everything is set.
You can track service logs from your IDE/terminal.

You can use the app through UI (container: `ui`) at

```bash
http://localhost:3000
```

If you also want to run UI locally which I would prefer for experimenting full-stack, please stop the container: `ui` and then follow the readme at [gRPC-Web Demo UI](https://github.com/uid4oe/grpc-web-demo-ui).

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
