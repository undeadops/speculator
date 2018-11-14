FROM golang:1.11-stretch AS build

RUN apt-get update && apt-get install -y --no-install-recommends \
		bash \
		ca-certificates \
		git \
 	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o speculator *.go



FROM debian:stretch-slim

RUN apt-get update && apt-get install  -y --no-install-recommends \
		ca-certificates \
		bash \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=build /app/speculator /bin/speculator

# Use CMD so it can be overriden
CMD ["/bin/speculator"]
