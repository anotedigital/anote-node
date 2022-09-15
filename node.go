package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/subosito/gotenv"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func initSeedFile() {
	OwnerAddress = os.Getenv("ADDRESS")

	if _, err := os.Stat("seed"); errors.Is(err, os.ErrNotExist) {
		seedStr := ""
		seed, encoded := generateSeed()
		PublicKey, PrivateKey = generateKeys(seed)
		key, encKey := generateApiKey()
		ip := getIP()

		seedStr += fmt.Sprintf("export SEED='%s'\n", seed)
		seedStr += fmt.Sprintf("export ENCODED='%s'\n", encoded)
		seedStr += fmt.Sprintf("export KEY='%s'\n", key)
		seedStr += fmt.Sprintf("export KEYENCODED='%s'\n", encKey)
		seedStr += fmt.Sprintf("export PUBLICIP='%s'", ip)

		f, _ := os.Create("seed")
		defer f.Close()
		f.Write([]byte(seedStr))
	} else {
		gotenv.Load("seed")
		seed := os.Getenv("SEED")
		PublicKey, PrivateKey = generateKeys(seed)
	}
}

func ping() {
	url, err := joinUrl(MasterNodeUrl, fmt.Sprintf("/ping/%s/%s", OwnerAddress, NodeAddress))
	if err != nil {
		log.Println(err.Error())
	}

	res, err := http.Get(url.String())
	if err != nil {
		log.Println(err.Error())
	}
	res.Body.Close()
}

func waitForAnotes() {
	time.Sleep(time.Second * 60)
}

func setScript() error {
	var networkByte = byte(55)
	var nodeURL = AnoteNodeURL

	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58(PublicKey)
	if err != nil {
		log.Println(err)
		return err
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58(PrivateKey)
	if err != nil {
		log.Println(err)
		return err
	}

	// Current time in milliseconds
	ts := uint64(time.Now().Unix() * 1000)

	gs, _ := base64.StdEncoding.DecodeString(generatorScript)

	tr := proto.NewUnsignedSetScriptWithProofs(
		2,
		networkByte,
		sender,
		gs,
		AnoteFee*2,
		ts)

	err = tr.Sign(networkByte, sk)
	if err != nil {
		log.Println(err)
		return err
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	client, err := client.NewClient(client.Options{BaseUrl: nodeURL, Client: &http.Client{}})
	if err != nil {
		log.Println(err)
		return err
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// // Send the transaction to the network
	_, err = client.Transactions.Broadcast(ctx, tr)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func callScript() error {
	return nil
}
