package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/p2p/netutil"
)

var charRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()-_+=")

var defaultChainDir = "ethbs\\private-chain"

var genesisTemplate = `{
	"config": {
		"chainId": 0,
		"homesteadBlock": 0,
		"eip155Block": 0,
		"eip158Block": 0
	},
	"alloc": {
		"[ACCOUNT_ADDR_1]": {
			"balance": "222222222"
		}
	},
	"coinbase"   : "0x0000000000000000000000000000000000000000",
	"difficulty" : "0x400",
	"extraData"  : "",
	"gasLimit"   : "0x2fefd8",
	"nonce"      : "[NONCE]",
	"mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
	"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
	"timestamp"  : "0x00"
}`

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := createChainDir(defaultChainDir); err != nil {
		return
	}

	if err := startBootNode("ethbs\\boot.key"); err != nil {
		return
	}

	address, err := createAccount()
	if err != nil {
		return
	}

	var genesisJSON = genesisTemplate
	genesisJSON = strings.Replace(genesisJSON, "[ACCOUNT_ADDR_1]", address.Hex(), 1)
	genesisJSON = strings.Replace(genesisJSON, "[NONCE]", fmt.Sprintf("0x%x", rand.Intn(4294967295)), 1)

	fmt.Printf("genesis.json: %s\n", genesisJSON)
}

func createChainDir(dir string) (error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0600); err != nil {
			fmt.Printf("could not create chain directory: %v\n", err)
			return err
		}
	}
	return nil
}

func startBootNode(filename string) (err error) {
	if err := createBootKey(filename); err != nil {
		return err
	}

	nodeKey, err := crypto.LoadECDSA(filename)
	if err != nil {
		fmt.Printf("could not load boot key file: %v\n", err)
		return err
	}

	natm, err := nat.Parse("none")
	if err != nil {
		fmt.Printf("could not load nat: %v\n", err)
		return err
	}

	var restrictList *netutil.Netlist
	tab, err := discover.ListenUDP(nodeKey, ":30301", natm, "", restrictList)
	if err != nil {
		fmt.Printf("could not start boot node: %v\n", err)
		return err
	}

	fmt.Printf("boot node started: %s\n", strings.Replace(tab.Self().String(), "[::]", "127.0.0.1", 1))
	return nil
}

func createBootKey(filename string) (err error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		nodeKey, err := crypto.GenerateKey()
		if err != nil {
			fmt.Printf("could not generate boot key: %v\n", err)
			return err
		}
		if err = crypto.SaveECDSA(filename, nodeKey); err != nil {
			fmt.Printf("could not create boot key file: %v\n", err)
			return err
		}
	}
	return nil
}

func createAccount() (common.Address, error) {
	var address common.Address
	var password = randString(16)
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	address, err := keystore.StoreKey(defaultChainDir, password, scryptN, scryptP)

	if err != nil {
		fmt.Printf("Failed to create account: %v", err)
		return address, err
	}
	fmt.Printf("Account address: {%x}, account password: %s\n", address, password)
	return address, nil
}

func randString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = charRunes[rand.Intn(len(charRunes))]
    }
    return string(b)
}