FROM golang:latest
RUN mkdir /app
ADD . /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go build -o main ./cmd/web*
EXPOSE 4000 4000
CMD ["/app/main"]