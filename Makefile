
rebuild-messages-event-proto:
	@echo "Rebuilding es/event.proto..."
	@protoc -I=. --go-grpc_out=. --go_out=. --go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative es/event.proto


rebuild-errors-proto:
	@echo "Rebuilding es/event.proto..."
	@protoc -I=. --go-grpc_out=. --go_out=. --go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative es/event.proto


add-module-replaces:
	@echo "Adding module replaces..."
	@go run internal/releasehelper/releasehelper.go addreplace


drop-module-replaces:
	@echo "Adding module replaces..."
	@go run internal/releasehelper/releasehelper.go dropreplace