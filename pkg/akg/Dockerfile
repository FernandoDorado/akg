FROM golang:1.16-alpine
WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
RUN go build

FROM alpine:latest
COPY --from=0 /src/akg .
COPY templates/ templates/
CMD ["./akg"]
