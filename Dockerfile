FROM golang:latest AS compile

WORKDIR /src
RUN mkdir /api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o /api/server

FROM scratch
COPY --from=compile /api/server /api/
EXPOSE 8080
CMD ["/api/server"]
