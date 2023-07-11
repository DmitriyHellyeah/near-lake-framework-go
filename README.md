# near-lake-framework-go

It allows you to build your own indexer that subscribes to the stream of blocks
from the NEAR Lake data source and create your own logic to process
the NEAR Protocol data.

---

## How to install
```
go get github.com/Dmitriyhellyeah/near-lake-framework-go
```
## Example

### init config with credentials 
```golang
    config := nearlake.InitDefaultTestnetConfig(131163980, 10)
    config.AWSAccessKeyId = "AKIA2M26MSLIO56JVRY5"
    config.AWSSecretAccessKey = "4Z1Mb4/QSdNkBTV8QWfckwi35EKCkXGS5VEFLhjA"
    
    nlc, err := nearlake.NewClient(config)
    if err != nil {
        log.Fatalf("error", err)
    }

```

### init config without credentials 
```golang
    config := nearlake.InitDefaultTestnetConfig(131163980, 10)

    nlc, err := nearlake.NewClientWithoutCredentials(config)
    if err != nil {
        log.Fatalf("error", err)
    }

```

### Streamer
```golang

    config := nearlake.InitDefaultTestnetConfig(131163980, 10)
    
    nlc, err := nearlake.NewClientWithoutCredentials(config)
    
    if err != nil {
        log.Fatalf("error", err)
    }
    
    channel := nlc.Streamer()
    
    for messages := range channel {
	    for _, message := range messages {
	        log.Println("BLOCK HEIGHT", message.Block.Header.Height)
	    }
	}
```

### Transaction Watcher & arguments decoding

```golang
    type ExampleOracle struct {
        Prices []struct {
            AssetId string `json:"asset_id"`
            Price   struct {
                Multiplier string `json:"multiplier"`
                Decimals   int    `json:"decimals"`
            } `json:"price"`
        } `json:"prices"`
    }


    config := nearlake.InitDefaultTestnetConfig(131547225, 10)
    
    nlc, err := nearlake.NewClientWithoutCredentials(config)
    
    if err != nil {
        log.Fatalf("error", err)
    }
    
    chann := nlc.TxWatcher([]string{"priceoracle.testnet"})
    
    var exampleOracle ExampleOracle
    
    for msg := range chann {
        log.Println(msg)
        msg.DecodeArgs(&exampleOracle)
        
        ll, _ := json.Marshal(exampleOracle)
        log.Println(string(ll))
    }
```
