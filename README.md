# Setup guide

First of all you should create a geth node (https://github.com/ethereum/go-ethereum)

1) Init a node
    ```$xslt
    geth --identity "MyNode1" --rpc --rpcport "8080" --rpccorsdomain "*" --datadir data --port "30303" --nodiscover --rpcapi "db,eth,net,web3, miner, personal" --networkid 1999 init genesis.json
    ```
    
2) Run the node
    ```$xslt
    geth --identity "MyNode" --rpc --rpcport "8080" --rpccorsdomain "*" --datadir data --port "30303" --nodiscover --rpcapi "db,eth,net,web3, miner, personal, admin" --networkid 1999
    ```
    
3) Enter to the node by following command
    ```$xslt
    geth attach http://localhost:8080
    ```
    
4) Create account
    ```$xslt
    personal.newAccount("password")
    miner.setEtherbase(eth.accounts[0])
    ```
    
5) Start mining
    ```$xslt
    miner.start()
    ```

# Run the app
    go run *.go    

# Transfer eth
    go run tcp_client.go send <sender> <recipient> <decimal amount> <password>
    
# Get lath transactions
    go run tcp_client.go get-last