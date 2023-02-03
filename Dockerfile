FROM amd64/golang:1.19.2-bullseye

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . ./

## Run tinygo install script
RUN bash ./scripts/install-tinygo.sh
RUN bash ./scripts/install-wabt.sh

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /api

ENV PORT 8080
EXPOSE $PORT

CMD ["/api", "serve"]
