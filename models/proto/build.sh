protoc --proto_path=. --go_out=paths=source_relative:. *.proto
#protoc-go-inject-tag -input="*.pb.go"
mv *.pb.go ../