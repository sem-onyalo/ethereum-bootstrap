package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/p2p/netutil"
)


func main() {
	createChainDir("ethbs\\private-chain")
	startBootNode("ethbs\\boot.key")
}

func createChainDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0600);
	}
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