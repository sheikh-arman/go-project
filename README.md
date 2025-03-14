https://www.youtube.com/watch?v=W7HK2yD0_U0&list=PLQ9_95hffac8_0bj5oeCe4FdxeNZi0UJ2&index=2











go get -u github.com/spf13/cobra@latest

go install github.com/spf13/cobra-cli@latest

cobra-cli init

cobra-cli add start

cobra-cli add subcommand --parent startCmd

go build -o live-chat

./live-chat start

expose the websocket on ws://localhost:8080/ws


orderbook api on ws://localhost:8080/orderbook 
orderbook will fetch data in every 2 second












