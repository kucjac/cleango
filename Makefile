
rebuild-messages-event-proto:
	@echo "Rebuilding es/event.proto..."
	@protoc -I=. --go-grpc_out=. --go_out=. --go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative es/event.proto


rebuild-errors-proto:
	@echo "Rebuilding es/event.proto..."
	@protoc -I=. --go-grpc_out=. --go_out=. --go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative es/event.proto