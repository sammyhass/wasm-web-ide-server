FROM amd64/golang:1.19.2-bullseye

ENV PORT 8080

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . ./

## Run tinygo install script
RUN bash ./scripts/install-tinygo.sh

RUN go build -o /ide-server

EXPOSE $PORT

CMD [ "/ide-server" ]