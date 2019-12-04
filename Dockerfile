FROM golang:alpine as builder
RUN apk --no-cache add gcc g++ make ca-certificates git
WORKDIR /go/src/github.com/robrohan/HungryLegs
COPY . .
RUN make build.server

FROM golang:alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/robrohan/HungryLegs/build/ ./
RUN ls -alFh
CMD ["./hungrylegs"]
