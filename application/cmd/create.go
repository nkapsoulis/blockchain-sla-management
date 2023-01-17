package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hyperledger/fabric-private-chaincode/samples/application/simple-cli-go/pkg"
	"github.com/hyperledger/fabric-private-chaincode/structs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create sla_json",
	Short: "create SLA on the chaincode",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile, err := os.Open(args[0])
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		// read our opened jsonFile as a byte array.
		byteValue, err := ioutil.ReadAll(jsonFile)
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

		client := pkg.NewClient(config)
		res := client.Invoke("CreateOrUpdateContract", string(byteValue))
		fmt.Println("> " + res)
	},
}
