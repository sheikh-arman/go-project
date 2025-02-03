go get -u github.com/spf13/cobra@latest

go install github.com/spf13/cobra-cli@latest

cobra-cli init

cobra-cli add start

go build -o live-chat

./live-chat start

cobra-cli add sub2 --parent start














