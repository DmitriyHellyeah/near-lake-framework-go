package types

import (
	"fmt"
	clientTypes "github.com/DmitriyHellyeah/nearclient/types"
)

type BlockHeight int

type SignedTransaction struct {
	SignerId   string                   `json:"signer_id"`
	PublicKey  string                   `json:"public_key"`
	Nonce      uint64                   `json:"nonce"`
	ReceiverId string                   `json:"receiver_id"`
	Actions    []ActionView `json:"actions"`
	Signature  string                   `json:"signature"`
	Hash       string                   `json:"hash"`
}

type IndexerExecutionOutcomeWithReceipt struct {
	ExecutionOutcome ExecutionOutcomeWithId `json:"execution_outcome"`
	Receipt          ReceiptView            `json:"receipt"`
}

type CostGasUsed struct {
	CostCategory string      `json:"cost_category"`
	Cost         string      `json:"cost"`
	GasUsed      BigInt `json:"gas_used"`
}

type ExecutionMetadata struct {
	Version    int           `json:"version"`
	GasProfile []CostGasUsed `json:"gas_profile,omitempty"`
}

type ExecutionOutcome struct {
	Logs        []string          `json:"logs"`
	ReceiptIds  ReceiptIds          `json:"receipt_ids"`
	GasBurnt    uint64            `json:"gas_burnt"`
	TokensBurnt BigInt       `json:"tokens_burnt,string"`
	ExecutorId  string            `json:"executor_id"`
	Status      Status            `json:"status"`
	Metadata    ExecutionMetadata `json:"metadata"`
}

type ReceiptIds []string

func (rIds ReceiptIds) GetFirstItem() (id string, err error) {
	if len(rIds) == 0 {
		err = fmt.Errorf("`receipt_ids` must contain at least one Receipt Id")
		return
	}
	id = rIds[0]
	return
}

type Status map[string]interface{}

func (s Status) IsUnknown() bool {
	_, ok := s["Unknown"]
	return ok
}

func (s Status) IsFailure() bool {
	_, ok := s["Failure"]
	return ok
}

func (s Status) IsSuccess() bool {
	_, ok1 := s["SuccessValue"]
	_, ok2 := s["SuccessReceiptId"]

	return ok1 || ok2
}

func (s Status) SuccessReceiptId() *string {
	_, ok := s["SuccessReceiptId"]
	if s.IsSuccess() && ok {
		return s["SuccessReceiptId"].(*string)
	}
	return nil
}

func (s Status) SuccessValue() *string {
	_, ok := s["SuccessValue"]
	if s.IsSuccess() && ok {
		return s["SuccessValue"].(*string)
	}
	return nil
}

type ExecutionOutcomeWithId struct {
	Proof     interface{}      `json:"proof"`
	BlockHash string           `json:"block_hash"`
	Id        string           `json:"id"`
	Outcome   ExecutionOutcome `json:"outcome"`
}

type IndexerExecutionOutcomeWithOptionalReceipt struct {
	ExecutionOutcome ExecutionOutcomeWithId `json:"execution_outcome"`
	Receipt          ReceiptView            `json:"receipt,omitempty"`
}

type IndexerTransactionWithOutcome struct {
	Transaction SignedTransaction                          `json:"transaction"`
	Outcome     IndexerExecutionOutcomeWithOptionalReceipt `json:"outcome"`
}

type IndexerChunk struct {
	Author       string                          `json:"author"`
	Header       clientTypes.ChunkHeader         `json:"header"`
	Transactions []IndexerTransactionWithOutcome `json:"transactions"`
	Receipts     []ReceiptView                   `json:"receipts"`
}

type IndexerShard struct {
	ShardId                  int                                  `json:"shard_id"`
	Chunk                    *IndexerChunk                        `json:"chunk"`
	ReceiptExecutionOutcomes []IndexerExecutionOutcomeWithReceipt `json:"receipt_execution_outcomes"`
	StateChanges             []StateChange                        `json:"state_changes"`
}

type StateChange struct {
	Cause map[string]interface{} `json:"cause"`
	Value map[string]interface{} `json:"value"`
}

type StreamerMessage struct {
	Block  clientTypes.BlockDetails `json:"block"`
	Shards []*IndexerShard    `json:"shards"`
}

type StreamerReceipt struct {
	Block  clientTypes.BlockDetails `json:"block"`
	Receipts []ReceiptView    `json:"shards"`
}