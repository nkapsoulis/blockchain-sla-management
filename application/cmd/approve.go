package cmd

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hyperledger/fabric-private-chaincode/samples/application/simple-cli-go/pkg"
	"github.com/hyperledger/fabric-private-chaincode/structs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(approveCmd)
}

var approveCmd = &cobra.Command{
	Use:   "approve sla_json username priv_key",
	Short: "approve SLA on the chaincode",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		// read our opened jsonFile as a byte array.
		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			panic(err)
		}

		// Verify that the data indeed fits the struct
		// before we send it to the chaincode
		var sla structs.SLA
		err = json.Unmarshal(byteValue, &sla)
		if err != nil {
			panic(err)
		}

		privKeyFile, err := os.Open(args[2])
		if err != nil {
			panic(err)
		}
		defer privKeyFile.Close()

		privKeyBytes, err := io.ReadAll(privKeyFile)
		if err != nil {
			panic(err)
		}
		privKey := ed25519.PrivateKey(privKeyBytes)
		signature := ed25519.Sign(privKey, byteValue)

		client := pkg.NewClient(config)
		res := client.Invoke("Approve", sla.ID, args[1], string(signature))
		fmt.Println("> " + res)
	},
}
