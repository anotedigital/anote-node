package main

import (
	"log"

	"github.com/anonutopia/gowaves"
)

func sendAINT(recipient string, amount int) {
	atr := &gowaves.AssetsTransferRequest{
		Recipient: recipient,
		Amount:    amount,
		Fee:       AintFee,
		Sender:    NodeAddress,
	}
	_, err := gowaves.WNC.AssetsTransfer(atr)
	if err != nil {
		log.Println(err.Error())
	}
}
