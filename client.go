package nearlake

import (
	"github.com/DmitriyHellyeah/near-lake-framework-go/types"
	clientTypes "github.com/DmitriyHellyeah/nearclient/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"sort"
	"sync"
	"time"
)


type Client struct {
	S3Client *s3.S3
	Config   S3Config
}

func NewClient(config S3Config) (client Client, err error) {
	return ClientInit(config, true)
}

func NewClientWithoutCredentials(config S3Config) (client Client, err error) {
	return ClientInit(config, false)
}

func ClientInit(config S3Config, cred bool) (client Client, err error) {
	cfg := aws.Config{
		Region:      aws.String(config.Region),
		MaxRetries:  aws.Int(config.MaxRetries),
	}
	if cred {
		cfg.Credentials = credentials.NewStaticCredentials(config.AWSAccessKeyId, config.AWSSecretAccessKey, "")
	}

	sess, err := session.NewSession(&cfg)

	if err != nil {
		log.Fatalf("Can't init s3 credentials. %s", err)
		return
	}

	client.S3Client = s3.New(sess)
	client.Config = config

	if err != nil {
		return
	}

	return
}

func (c *Client) Streamer() chan []types.StreamerMessage {
	messageChannel := make(chan *types.StreamerMessage)
	sortedMessages := make(chan []types.StreamerMessage)

	go c.start(c.Config, messageChannel)
	go c.sortMessages(messageChannel, sortedMessages)

	return sortedMessages
}

func (c *Client) start(config S3Config, messageChannel chan *types.StreamerMessage)  {
	var startFromBlockHeight = config.StartAfter

	for {
		// get list blocks from s3
		blockHeightsPrefixes, err := c.ListBlocks(
			config.Bucket,
			startFromBlockHeight,
		)

		if err != nil {
			log.Printf("Can't fetch data from s3 Method[ListBlocks], retry in %s  ... - %s",time.Second * c.Config.RequestWaitingTimeout, err)
			time.Sleep(time.Second * c.Config.RequestWaitingTimeout)
			continue
		}

		// if list of blocks less than c.Config.BlocksCountWaiting, waiting until they appeared in s3
		// delay applied to avoid requests spamming
		if len(blockHeightsPrefixes) < c.Config.BlocksCountWaiting {
			log.Printf("No new blocks on S3, retry in %s  ...", time.Second * c.Config.RequestWaitingTimeout)
			time.Sleep(time.Second * c.Config.RequestWaitingTimeout)
			continue
		}

		log.Printf("Received [%d] blocks from s3 %v", len(blockHeightsPrefixes), blockHeightsPrefixes)

		// get block details for each block from s3 and send it to channel
		var wg sync.WaitGroup
		startFromBlockHeight = blockHeightsPrefixes[len(blockHeightsPrefixes)-1] + 1
		for _, blockHeight := range blockHeightsPrefixes {
			wg.Add(1)
			go func(blockHeight uint64) {
				defer func() { wg.Done() }()
				streamerMessage, err1 := c.FetchStreamerMessage(
					c.Config.Bucket,
					blockHeight)
				if err1 != nil {
					log.Printf("Can't fetch data from s3. Method[FetchStreamerMessage]. %s", err1)
				}
				messageChannel <- streamerMessage
			}(blockHeight)
		}
		wg.Wait()
	}
}

func (c *Client) sortMessages(messageChannel chan *types.StreamerMessage, sortedMessages chan []types.StreamerMessage) {
	var messages []types.StreamerMessage

	for {
		select {
		case message := <-messageChannel:
			messages = append(messages, *message)
			if uint64(len(messages)) < c.Config.NumberOfBlocksRequested {
				continue
			}
			sort.Slice(messages, func(i, j int) bool {
				return messages[i].Block.Header.Height < messages[j].Block.Header.Height
			})
			sortedMessages <- messages
			messages = make([]types.StreamerMessage, 0)
		}
	}
}

func (c *Client) TxWatcher(watchingList []string) chan *types.FunctionCall {
	functionCallChannel := make(chan *types.FunctionCall)
	channel := c.Streamer()
	var wantedReceiptIds []string

	go func() {
		defer close(channel)
		defer close(functionCallChannel)

		for messages := range channel {
			for _, message := range messages {
				for _, shard := range message.Shards {
					if shard.Chunk == nil {continue}
					retrieveReceiptIdFromTx(shard.Chunk.Transactions, &wantedReceiptIds, watchingList)
					retrieveFunctionCallActionFromExecutionOutcome(shard.ReceiptExecutionOutcomes, &wantedReceiptIds, functionCallChannel)
				}
			}
		}
	}()
	return functionCallChannel
}

type ReceiptViewWithBlockHeaderChannel struct {
	types.ReceiptView
	BlockHeader clientTypes.BlockHeader
}

func (c *Client) TxWatcherReceiptWithBlockHeader(watchingList []string) chan *ReceiptViewWithBlockHeaderChannel {
	receiptViewWithBlockHeaderChannel := make(chan *ReceiptViewWithBlockHeaderChannel)
	channel := c.Streamer()
	var wantedReceiptIds []string

	go func() {
		defer close(channel)
		defer close(receiptViewWithBlockHeaderChannel)

		for messages := range channel {
			for _, message := range messages {
				for _, shard := range message.Shards {
					if shard.Chunk == nil {continue}
					retrieveReceiptIdFromTx(shard.Chunk.Transactions, &wantedReceiptIds, watchingList)
					retrieveActionFromExecutionOutcome(shard.ReceiptExecutionOutcomes, &wantedReceiptIds, receiptViewWithBlockHeaderChannel, message.Block.Header)
				}
			}
		}
	}()
	return receiptViewWithBlockHeaderChannel
}

func retrieveReceiptIdFromTx(txs []types.IndexerTransactionWithOutcome, wantedReceiptIds *[]string, watchingList []string) {
	for _, tx := range txs {
		if isTxReceiverWatched(tx, watchingList) {
			receiptId, err := tx.Outcome.ExecutionOutcome.Outcome.ReceiptIds.GetFirstItem()
			if err != nil {
				log.Printf("receipt_ids is empty. %s", err)
				continue
			}
			*wantedReceiptIds = append(*wantedReceiptIds, receiptId)
		}
	}
}

func retrieveActionFromExecutionOutcome(outcome []types.IndexerExecutionOutcomeWithReceipt, wantedReceiptIds *[]string, channel chan *ReceiptViewWithBlockHeaderChannel, blockHeader clientTypes.BlockHeader) {
	for _, executionOutcome := range outcome {
		if contains(*wantedReceiptIds, executionOutcome.Receipt.ReceiptId) {
			channel <- &ReceiptViewWithBlockHeaderChannel{executionOutcome.Receipt, blockHeader}
			// remove receiptId from wantedReceiptIds list
			remove(wantedReceiptIds, executionOutcome.Receipt.ReceiptId)
		}
	}
}

func retrieveFunctionCallActionFromExecutionOutcome(outcome []types.IndexerExecutionOutcomeWithReceipt, wantedReceiptIds *[]string, channel chan *types.FunctionCall) {
	for _, executionOutcome := range outcome {
		if contains(*wantedReceiptIds, executionOutcome.Receipt.ReceiptId) {
			receipt := executionOutcome.Receipt.Receipt

			actionList, err := receipt.GetAction()
			if err != nil {
				log.Printf("error when get action from receipt [%s]. %s", executionOutcome.Receipt.ReceiptId, err)
				continue
			}
			for _, action := range actionList.Actions {
				fc, err := action.GetFunctionCall()
				if err != nil {
					log.Printf("error when get FunctionCall from action in receiptId [%s]. %s", executionOutcome.Receipt.ReceiptId, err)
					continue
				}
				channel <- fc
			}
			// remove receiptId from wantedReceiptIds list
			remove(wantedReceiptIds, executionOutcome.Receipt.ReceiptId)
		}
	}
}