

# create internal/proto directory if it doesn't exist
if [ ! -d "internal/proto" ]; then
    mkdir internal/proto
fi

protoc --go_out=internal/proto --go_opt=paths=source_relative \
    --go-grpc_out=internal/proto --go-grpc_opt=paths=source_relative \
    runit.proto