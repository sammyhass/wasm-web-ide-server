#  Projects IDE API

This repository contains the source code for the Go API utilised by the 'Projects IDE' and is part of my submi

##  Local Development Setup

You will need to have Go installed on your machine. You can download Go from [here](https://golang.org/dl/).

###  Environment Variables

You will be required to provide the following environment variables, these can be provided in a `.env` file in the root of the project:

* `PORT` - The port that the API will run on (defaults to 8080)
* `POSTGRES_PORT` - The port that the Postgres database is running on
* `POSTGRES_HOST` - The host that the Postgres database is running on
* `POSTGRES_USER` - The username for the Postgres database
* `POSTGRES_PASSWORD` - The password for the Postgres database
* `POSTGRES_DB` - The name of the Postgres database
* `S3_ACCESS_KEY_ID` - AWS Access Key ID so that the API can access the S3 bucket
* `S3_SECRET_ACCESS_KEY` - AWS Secret Access Key so that the API can access the S3 bucket
* `S3_BUCKET` - The name of the S3 bucket that the API will use to store projects
* `JWT_SECRET` - The secret used to sign the JWT tokens
* `CORS_ALLOW_ORIGIN` - The origin that the API will allow CORS requests from
* `DEPLOY_URL` - The URL that the API will be deployed to

### Performing Migrations

Prior to starting the server, you will need to perform a database migration so that the postgres database is setup correctly. To do this, run the following command:

```bash
go run main.go migrate
```

###  Installing WebAssembly Related Dependencies

The server makes use of the following tools as part of its WebAssembly compilation pipeline:
* [AssemblyScript](https://www.assemblyscript.org/) - A TypeScript-like language that compiles to WebAssembly
* [WebAssembly Binary Toolkit (wabt)](https://github.com/WebAssembly/wabt) - A toolkit for working with WebAssembly binaries and text formats
* [TinyGo](https://tinygo.org/) - A Go compiler for WebAssembly

You will need to ensure each of these are installed on your machine. Scripts for installing each of these dependencies are provided in the [scripts](scripts) directory. Run all of these scripts from the root of the project. Note that these scripts expect a Debian environment so for different environment it may be required to install these dependencies using other operating-system specific approaches.

### Serving the API

Once you have setup environment variables and performed the necessary migrations, you can run the API by running the following command:

```bash
go run main.go serve
```

This will start the API on the port specified in the environment variables.

## Docker

The provided [Dockerfile](Dockerfile) can be used to build a Docker image for the API and is used to deploy the API in a production environment. Of course, you will need to provide the [required environment](#environment-variables) variables to the Docker container when running it.

##  Running Tests

To run the unit tests for the API, run the following command:

```bash
go test ./...
```
