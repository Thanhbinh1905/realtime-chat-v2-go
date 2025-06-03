protoc \
  --proto_path=./auth-service/api/auth/v1 \
  --proto_path=./auth-service/third_party \
  --go_out=paths=source_relative:./auth-service/api/auth/v1 \
  --go-grpc_out=paths=source_relative:./auth-service/api/auth/v1 \
  --grpc-gateway_out=paths=source_relative:./auth-service/api/auth/v1 \
  ./auth-service/api/auth/v1/auth.proto