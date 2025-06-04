// Config protoc auth-service

protoc \
  --proto_path=./auth-service/api/auth/v1 \
  --proto_path=./shared/third_party \
  --go_out=paths=source_relative:./auth-service/api/auth/v1 \
  --go-grpc_out=paths=source_relative:./auth-service/api/auth/v1 \
  --grpc-gateway_out=paths=source_relative:./auth-service/api/auth/v1 \
  ./auth-service/api/auth/v1/auth.proto

// Config protoc user-service

protoc \
  --proto_path=./user-service/api/auth/v1 \
  --proto_path=./shared/third_party \
  --go_out=paths=source_relative:./user-service/api/auth/v1 \
  --go-grpc_out=paths=source_relative:./user-service/api/auth/v1 \
  --grpc-gateway_out=paths=source_relative:./user-service/api/auth/v1 \
  ./user-service/api/auth/v1/user.proto
