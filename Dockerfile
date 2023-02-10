
FROM amd64/golang:1.19.2-bullseye


RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - && \
	apt-get install -y nodejs &&\
	node -v && \
	npm -v && \
	npm install -g assemblyscript && \
	asc --version

## Run tinygo install script
RUN wget https://github.com/tinygo-org/tinygo/releases/download/v0.26.0/tinygo_0.26.0_amd64.deb && \
	dpkg -i tinygo_0.26.0_amd64.deb && \
	rm tinygo_0.26.0_amd64.deb && \
	tinygo version

RUN apt-get update && \
	apt-get install wabt


WORKDIR /app
COPY go.mod .
COPY go.sum .

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o=./api


CMD ["./api", "serve"]
