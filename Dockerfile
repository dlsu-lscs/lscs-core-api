FROM golang:1.23.2 as build

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /bin/lscs-central-auth ./cmd/api/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build /bin/lscs-central-auth /bin/

EXPOSE 42069

CMD ["/bin/lscs-central-auth"]
