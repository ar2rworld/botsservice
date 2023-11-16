# BotsService

- Is an attempt to collect simple bots into one repo
- Use gRPC to transfer messages from clients to the main server

### Generate protobuf
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    app/messageservice/messageservice.proto
```
