FROM golang:1.24

RUN apt update -y && apt install -y protobuf-compiler

RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.14.0 \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.14.0

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
