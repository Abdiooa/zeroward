FROM golang:1.21 as build
LABEL stage=intermediate
WORKDIR /app

COPY . .

RUN make build

FROM alpine:latest
COPY --from=build /app/bin/zeroward /bin/zeroward
CMD ["/bin/zeroward"]