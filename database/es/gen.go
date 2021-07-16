package es

// Generate event proto file.
//go:generate protoc -I=. --go_out=. event.proto --go_opt=paths=source_relative

// Generate event state proto file.
//go:generate protoc -I=. --go_out=. eventstate/eventstate.proto --go_opt=paths=source_relative
