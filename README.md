#  Projects IDE API

This repository contains the source code for the Go API utilised by the 'Projects IDE' and is part of my submi

##  Setup

You will need to have Go installed on your machine. You can download Go from [here](https://golang.org/dl/).

###  Environement Variables

You will be required to provide the following environment variables in a `.env` file in the root of the project:

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

### Serving the API

Once you have setup environment variables and performed the necessary migrations, you can run the API by running the following command:

```bash
go run main.go serve
```

This will start the API on the port specified in the environment variables.

##  Running Tests

To run the unit tests for the API, run the following command:

```bash
go test ./...
```
