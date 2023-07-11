package nearlake

import (
	"encoding/json"
	"fmt"
	"github.com/DmitriyHellyeah/nearclient/types"
	lakeTypes "github.com/DmitriyHellyeah/near-lake-framework-go/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"strconv"
	"strings"
	"time"
)

func (c *Client) ListBlocks(bucketName string, startFromBlockHeight uint64) (blockHeightList []uint64, err error) {
	startAfter := fmt.Sprintf("%012d", startFromBlockHeight)
	res, err := c.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:       aws.String(bucketName),
		Delimiter:    aws.String(c.Config.Delimiter),
		MaxKeys:      aws.Int64(int64(c.Config.NumberOfBlocksRequested)),
		StartAfter:   aws.String(startAfter),
		Prefix:       aws.String(""),
		RequestPayer: aws.String(c.Config.RequestPayer),
	})

	if err != nil {
		return nil, err
	}

	for _, prefix := range res.CommonPrefixes {
		if prefix.Prefix != nil {
			p := *prefix.Prefix
			p = strings.TrimLeft(p[:len(p)-1], "0")
			i, err := strconv.ParseUint(p, 0, 64)
			if err != nil {
				return nil, err
			}

			blockHeightList = append(blockHeightList, i)
		}
	}

	return
}

func (c *Client) FetchStreamerMessage(bucketName string, blockHeight uint64) (*lakeTypes.StreamerMessage, error) {
	blockId := fmt.Sprintf("%012d/block.json", blockHeight)
	var block types.BlockDetails
	var streamer lakeTypes.StreamerMessage

	response, err := c.S3Client.GetObject(&s3.GetObjectInput{
		Bucket:       aws.String(bucketName),
		Key:          aws.String(blockId),
		RequestPayer: aws.String(c.Config.RequestPayer),
	})

	if err != nil {
		log.Printf("Can't get object from s3. Method [GetObject] %s", err)
		return nil, err
	}

	defer func() {
		err = response.Body.Close()
	}()

	decoder := json.NewDecoder(response.Body)

	if err = decoder.Decode(&block); err != nil {
		return nil, err
	}
	err = response.Body.Close()

	streamer.Block = block

	for _, shard := range block.Chunks {
		res, err := c.FetchShardOrRetry(bucketName, blockHeight, shard.ShardId)
		if err != nil {
			log.Printf("Can't get object from s3. Method [FetchShardOrRetry] %s", err)
		}

		streamer.Shards = append(streamer.Shards, res)
	}

	return &streamer, nil
}

func (c *Client) FetchShardOrRetry(bucketName string, blockHeight uint64, shardId int) (*lakeTypes.IndexerShard, error) {
	shard := fmt.Sprintf("%012d/shard_%d.json", blockHeight, shardId)
	var indexerShard lakeTypes.IndexerShard
	for {
		response, err := c.S3Client.GetObject(&s3.GetObjectInput{
			Bucket:       aws.String(bucketName),
			Key:          aws.String(shard),
			RequestPayer: aws.String("requester"),
		})

		// add delay coz shards can appears later then the blocks list
		if err != nil {
			log.Printf("Can't get object from s3. Shard [%s]. %s. Will sleep %d seconds an retry", shard, err, c.Config.ShardsWaitingTimeout)
			time.Sleep(time.Second * c.Config.ShardsWaitingTimeout)
			continue
		}

		decoder := json.NewDecoder(response.Body)

		if err = decoder.Decode(&indexerShard); err != nil {
			return nil, err
		}
		err = response.Body.Close()
		return &indexerShard, nil
	}

}
