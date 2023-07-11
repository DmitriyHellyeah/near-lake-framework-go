package types

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
)

type ActionView struct {
	ActionObject *map[string]interface{} `json:"action_object,omitempty"`
	ActionString *string                 `json:"action_string,omitempty"`
}

func (a *ActionView) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case '"':
		if err := json.Unmarshal(data, &a.ActionString); err != nil {
			return err
		}
	case '{':
		if err := json.Unmarshal(data, &a.ActionObject); err != nil {
			return err
		}
	}
	return nil
}

type DeployContract struct {
	Code string `json:"code"`
}

type Transfer struct {
	Deposit *big.Int `json:"deposit"`
}

type Stake struct {
	Stake     *big.Int `json:"stake"`
	PublicKey string   `json:"public_key"`
}

type FunctionCall struct {
	MethodName string       `json:"method_name"`
	Args       string       `json:"args"`
	Gas        uint64 		`json:"gas"`
	Deposit    *BigInt      `json:"deposit"`
}

type PermissionView struct {
	PermissionObject *PermissionObjectView `json:"permission_object,omitempty"`
	PermissionString *string                 `json:"permission_string,omitempty"`
}

type PermissionObjectView struct {
	FunctionCall struct {
		Allowance   BigInt        `json:"allowance"`
		MethodNames []string `json:"method_names"`
		ReceiverId  string        `json:"receiver_id"`
	} `json:"FunctionCall"`
}

func (pv *PermissionView) IsFullAccess() (ok bool) {
	if pv.PermissionString != nil {
		return *pv.PermissionString == "FullAccess"
	}
	return
}

func (pv *PermissionView) IsFunctionAccess() (ok bool) {
	return pv.PermissionObject != nil
}

func (pv *PermissionView) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case '"':
		if err := json.Unmarshal(data, &pv.PermissionString); err != nil {
			return err
		}
	case '{':
		if err := json.Unmarshal(data, &pv.PermissionObject); err != nil {
			return err
		}
	}
	return nil
}

type AddKey struct {
	AccessKey struct {
		Nonce      int    `json:"nonce"`
		PermissionView PermissionView `json:"permission"`
	} `json:"access_key"`
	PublicKey string `json:"public_key"`
}

func (fc *FunctionCall) DecodeArgs(to interface{}) (err error) {
	rawArgs, err := base64.StdEncoding.DecodeString(fc.Args)
	if err != nil {
		return
	}

	err = json.Unmarshal(rawArgs, to)
	if err != nil {
		return
	}
	return
}

func (a *ActionView) IsDeployContract() (ok bool) {
	if (a.ActionObject) != nil {
		_, ok = (*a.ActionObject)["DeployContract"]
		return
	}
	return
}

func (a *ActionView) GetDeployContract() (dc *DeployContract, err error)  {
	if a.IsDeployContract() {
		var data []byte
		data, err = json.Marshal((*a.ActionObject)["DeployContract"])
		if err != nil {
			return
		}
		err = json.Unmarshal(data, &dc)
		if err != nil {
			return
		}
		return
	}
	return
}

func (a *ActionView) IsFunctionCall() (ok bool) {
	if (a.ActionObject) != nil {
		_, ok = (*a.ActionObject)["FunctionCall"]
		return
	}
	return
}

func (a *ActionView) GetFunctionCall() (fc *FunctionCall, err error) {
	if a.IsFunctionCall() {
		var data []byte
		data, err = json.Marshal((*a.ActionObject)["FunctionCall"])
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

func (a *ActionView) IsTransfer() (ok bool) {
	if (a.ActionObject) != nil {
		_, ok = (*a.ActionObject)["Transfer"]
		return
	}
	return
}

func (a *ActionView) IsStake() (ok bool) {
	if (a.ActionObject) != nil {
		_, ok = (*a.ActionObject)["Stake"]
		return
	}
	return
}

func (a *ActionView) IsAddKey() (ok bool) {
	if (a.ActionObject) != nil {
		_, ok = (*a.ActionObject)["AddKey"]
		return
	}
	return
}

func (a *ActionView) GetAddKey() (ak *AddKey, err error) {
	if a.IsAddKey() {
		var data []byte
		data, err = json.Marshal((*a.ActionObject)["AddKey"])
		if err != nil {
			return
		}
		err = json.Unmarshal(data, &ak)
		if err != nil {
			return
		}
		return
	}
	return
}


func (a *ActionView) IsCreateAccount() (ok bool)  {
	if a.ActionString != nil {
		return *a.ActionString == "CreateAccount"
	}
	return
}