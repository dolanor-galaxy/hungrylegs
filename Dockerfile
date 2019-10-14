FROM golang:alpine as builder
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/therohans/HungryLegs
COPY . .
RUN make build

FROM golang:alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/therohans/HungryLegs/build/ ./
RUN ls -alFh
CMD ["./hungrylegs"]
