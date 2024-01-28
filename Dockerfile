FROM golang:1.21 as build
LABEL stage=intermediate
WORKDIR /app

COPY . .

RUN make build

FROM scratch as scratch
COPY --from=build /app/bin/zeroward /bin/zeroward
CMD ["./zeroward"]