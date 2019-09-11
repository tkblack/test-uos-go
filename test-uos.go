package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/lialvin/uos-go"
)

const (
	//chain_id   = "689f5b7d2a4ee8777e42d775cd3edf7a2be85845462ef7a6cb229d3f8b2c21ba"
	contract    = "uosio.token"
	accountFrom = "testaccount1"
	accountTo   = "testaccount2"
	udisk       = "testaccountb"
	perm        = "active"
	privatekey  = "5JoC5N2oTCyQrcG3R9riGbg8qVWPyDgUa6PnhSTidgRUbns6HTS"
	privatekey1 = "5K89DFV1yjqEqcZsbnWgyjoSmcmne4ZapEcu8Ch4nncRLiKm7x6"
)

//
func pushTransactionUOS() {
	url := "https://testrpc1.uosio.org:20580"
	api := uos.New(url)

	keyBag := uos.NewKeyBag()
	keyBag.ImportPrivateKey(privatekey1)
	keyBag.ImportPrivateKey(privatekey)
	api.SetSigner(keyBag)

	quantity, err := uos.NewUOSAssetFromString("1.0000 UOS")
	if err != nil {
		panic(fmt.Errorf("invalid quantity: %s", err))
	}
	memo := ""

	txOpts := &uos.TxOptions{}
	if err := txOpts.FillFromChain(api); err != nil {
		panic(fmt.Errorf("filling tx opts: %s", err))
	}

	//uosio.token合约transfer的参数列表
	type Transfer struct {
		From     uos.AccountName `json:"from"`
		To       uos.AccountName `json:"to"`
		Quantity uos.Asset       `json:"quantity"`
		Memo     string          `json:"memo"`
	}

	//testaccountb合约setstorage的参数列表
	type Setstorage struct {
		User    uos.AccountName `json:"user"`
		Storage int64           `json:"storage"`
	}

	//构建action---transfer
	action := uos.Action{
		Account: uos.AccountName("uosio.token"),
		Name:    uos.ActionName("transfer"),
		Authorization: []uos.PermissionLevel{
			{Actor: uos.AccountName(accountFrom), Permission: uos.PermissionName("active")},
		},
		ActionData: uos.NewActionData(Transfer{
			From:     uos.AccountName(accountFrom),
			To:       uos.AccountName(accountTo),
			Quantity: quantity,
			Memo:     memo,
		}),
	}

	//构建action---setstorage
	action1 := uos.Action{
		Account: uos.AccountName(udisk),
		Name:    uos.ActionName("setstorage"),
		Authorization: []uos.PermissionLevel{
			{Actor: uos.AccountName(udisk), Permission: uos.PermissionName("active")},
		},
		ActionData: uos.NewActionData(Setstorage{
			User:    uos.AccountName("111111111113 "), //纯数字账户，注意增加空格（注意在后面加，前面加会报错，正常情况前后加都可以）
			Storage: 1,
		}),
	}
	tx := uos.NewTransaction([]*uos.Action{&action1, &action}, txOpts)
	//tx := uos.NewTransaction([]*uos.Action{token.NewTransfer(uos.AccountName(accountFrom), uos.AccountName(accountTo), quantity, memo)}, txOpts)

	signedTx, packedTx, err := api.SignTransaction(tx, txOpts.ChainID, uos.CompressionNone)
	if err != nil {
		panic(fmt.Errorf("sign transaction: %s", err))
	}

	content, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		panic(fmt.Errorf("json marshalling transaction: %s", err))
	}

	fmt.Println(string(content))
	fmt.Println()

	response, err := api.PushTransaction(packedTx)
	if err != nil {
		panic(fmt.Errorf("push transaction: %s", err))
	}

	fmt.Printf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
}

func main() {
	pushTransactionUOS()
}
