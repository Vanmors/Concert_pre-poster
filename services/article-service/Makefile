#!make

proto:
	protoc 	--go_out=./protos/ --go-grpc_out=./protos/ \
			--go_opt=paths=source_relative \
			--go-grpc_opt=paths=source_relative \
			--proto_path=./contract/ article.proto