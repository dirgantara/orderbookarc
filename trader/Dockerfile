FROM golang:1.23.1

WORKDIR /app

# install postgres client
RUN apt-get update && apt-get install -y postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["sh", "entrypoint.sh"]