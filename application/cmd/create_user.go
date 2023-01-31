package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hyperledger/fabric-private-chaincode/application/pkg"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.AddCommand(createUserCmd)
}

var createUserCmd = &cobra.Command{
	Use:   "user username",
	Short: "create user on the chaincode",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			fmt.Printf("Generation error : %s", err)
			os.Exit(1)
		}

		b, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: b,
		}

		err = ioutil.WriteFile(name, pem.EncodeToMemory(block), 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		client := pkg.NewClient(config)
		res := client.Invoke("CreateUser", name, string(publicKey), "500")
		fmt.Println("> " + res)
	},
}
