package types

import (
	"encoding/json"
)

type ReceiptView struct {
	PredecessorId string  `json:"predecessor_id"`
	ReceiverId    string  `json:"receiver_id"`
	ReceiptId     string  `json:"receipt_id"`
	Receipt       Receipt `json:"receipt"`
}

type DataReceiverView struct {
	DataId     string `json:"data_id"`
	ReceiverId string `json:"receiver_id"`
}

type Data struct {
	DataId string
	Data   []uint8
}

type Action struct {
	SignerId            string             `json:"signer_id"`
	SignerPublicKey     string             `json:"signer_public_key"`
	GasPrice            *BigInt           `json:"gas_price"`
	OutputDataReceivers []DataReceiverView `json:"output_data_receivers"`
	InputDataIds        []string           `json:"input_data_ids"`
	Actions             []ActionView       `json:"actions"`
}

type Receipt map[string]interface{}


func (r *Receipt) IsAction() bool {
	_, ok := (*r)["Action"]
	return ok
}

func (r *Receipt) IsData() bool {
	_, ok := (*r)["Data"]
	return ok
}

func (r *Receipt) GetAction() (action *Action, err error) {
	if r.IsAction() {
		var data []byte
		data, err = json.Marshal((*r)["Action"])
		if err != nil {
			return
		}
		err = json.Unmarshal(data, &action)
		if err != nil {
			return
		}
		return
	}
	return
}

func (r *Receipt) GetData() *Data {
	if r.IsData() {
		return (*r)["Data"].(*Data)
	}
	return nil
}