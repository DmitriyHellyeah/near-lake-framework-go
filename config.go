package nearlake

import "time"

const (
	MainnetBucketName = "near-lake-data-mainnet"
	TestnetBucketName = "near-lake-data-testnet"
	AwsRegion         = "eu-central-1"
)

type S3Config struct {
	AWSAccessKeyId          string
	AWSSecretAccessKey      string
	Region                  string
	Bucket                  string
	Delimiter               string
	MaxRetries              int
	StartAfter              uint64
	NumberOfBlocksRequested uint64
	RequestPayer            string
	BlocksCountWaiting      int           // To prevent requests spamming if data does not exist yet
	RequestWaitingTimeout   time.Duration // To prevent requests spamming if data does not exist yet
	ShardsWaitingTimeout    time.Duration // sometimes shards can appears later then the blocks list
}

func InitDefaultMainnetConfig(startAfter, numberOfBlockRequested uint64) S3Config {
	cfg := InitDefaultConfig()
	cfg.Bucket = MainnetBucketName
	cfg.StartAfter = startAfter
	cfg.NumberOfBlocksRequested = numberOfBlockRequested
	return cfg
}

func InitDefaultTestnetConfig(startAfter, numberOfBlockRequested uint64) S3Config {
	cfg := InitDefaultConfig()
	cfg.Bucket = TestnetBucketName
	cfg.StartAfter = startAfter
	cfg.NumberOfBlocksRequested = numberOfBlockRequested
	return cfg
}

func InitDefaultConfig() S3Config {
	return S3Config{
		Region: AwsRegion,
		Delimiter: "/",
		MaxRetries: 3,
		RequestPayer: "requester",
		BlocksCountWaiting: 10,
		RequestWaitingTimeout: 10,
		ShardsWaitingTimeout: 2,
	}
}
