https://www.youtube.com/watch?v=W7HK2yD0_U0&list=PLQ9_95hffac8_0bj5oeCe4FdxeNZi0UJ2&index=2











go get -u github.com/spf13/cobra@latest

go install github.com/spf13/cobra-cli@latest

cobra-cli init

cobra-cli add start


grpc:

protoc --go_out=. --go-grpc_out=. --proto_path=./pkg/proto test.proto 

./go-project server

./go-project client













