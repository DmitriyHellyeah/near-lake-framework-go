package types

import (
	"encoding/json"
)

// ActionViewSimple if you use this type, you will need to skip decoding errors with string action
type ActionViewSimple map[string]interface{}

func (a *ActionViewSimple) IsDeployContract() bool {
	_, ok := (*a)["DeployContract"]
	return ok
}

func (a *ActionViewSimple) IsFunctionCall() bool {
	_, ok := (*a)["FunctionCall"]
	return ok
}


func (a *ActionViewSimple) IsTransfer() bool {
	_, ok := (*a)["Transfer"]
	return ok
}

func (a *ActionViewSimple) IsStake() bool {
	_, ok := (*a)["Stake"]
	return ok
}

func (a *ActionViewSimple) IsAddKey() (ok bool) {
	_, ok = (*a)["AddKey"]
	return
}


func (a *ActionViewSimple) GetFunctionCall() (fc *FunctionCall, err error) {
	if a.IsFunctionCall() {
		var data []byte
		data, err = json.Marshal((*a)["FunctionCall"])
		if err != nil {
			return
		}
		err = json.Unmarshal(data, &fc)
		if err != nil {
			return
		}
		return
	}
	return
}