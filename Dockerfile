FROM golang:1.23-alpine AS build

WORKDIR /app/

COPY . .

RUN go build -o chord .

FROM golang:1.23-alpine AS runtime

WORKDIR /app/

COPY --from=build /app/chord .

ENTRYPOINT ["./chord"]