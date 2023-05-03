package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/hyperledger/fabric-private-chaincode/clients/utils"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils/ledger"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils/pkg"
	t "github.com/hyperledger/fabric-private-chaincode/clients/utils/types"
	"github.com/hyperledger/fabric-private-chaincode/lib/ipfs"

	shell "github.com/ipfs/go-ipfs-api"
)

const IPFSHost = "localhost:5001"

var Config *pkg.Config

func main() {
	Config = utils.InitConfig()

	ctx := context.Background()
	sh := shell.NewShell(IPFSHost)

	err := ipfs.CreateRootFolder(ctx, sh)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(5 * time.Second)
		log.Println("Checking for new violations")
		dirs, err := sh.FilesLs(ctx, "/sla")
		if err != nil {
			log.Printf("failed while checking sla folders: %v\n", err)
			continue
		}

		for _, d := range dirs {
			metricsBytes, err := ipfs.GetUnreadMetrics(ctx, sh, d.Name)
			if err != nil {
				log.Printf("failed while reading sla %v folder: %v\n", d.Name, err)
				continue
			}

			if metricsBytes == nil {
				log.Printf("No metrics to read in sla %v folder\n", d.Name)
				continue
			}

			var metrics []t.Metric
			err = json.Unmarshal(metricsBytes, &metrics)
			if err != nil {
				log.Printf("failed to unmarshal metrics json: %v\n", err)
				continue
			}

			failed := make([]string, 0)
			for _, m := range metrics {
				err = ledger.CheckForViolation(Config, m)
				if err != nil {
					log.Printf("failed while submitting violation %v: %v\n", m.ID, err)
					failed = append(failed, m.ID)
					continue
				}
				log.Printf("Submitted violation %v", m.ID)
			}

			// Don't delete the violations that have not been submitted
			remainingMetrics := ""
			if len(failed) != 0 {
				remainingMetrics += strings.Join(failed, " ")
			}
			err = ipfs.OverwriteUnreadViolations(ctx, sh, d.Name, remainingMetrics)
			if err != nil {
				log.Printf("failed while clearing violations for SLA %v: %v\n", d.Name, err)
				continue
			}
		}
	}

}
