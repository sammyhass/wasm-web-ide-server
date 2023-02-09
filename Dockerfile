FROM amd64/golang:1.19.2-bullseye

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . ./

## Run tinygo install script
RUN wget https://github.com/tinygo-org/tinygo/releases/download/v0.26.0/tinygo_0.26.0_amd64.deb && \
	dpkg -i tinygo_0.26.0_amd64.deb && \
	rm tinygo_0.26.0_amd64.deb && \
	tinygo version

RUN apt-get update && \
	apt-get install wabt

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o=./api


CMD ["./api", "serve"]
