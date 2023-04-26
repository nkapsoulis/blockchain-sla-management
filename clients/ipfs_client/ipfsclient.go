package main

import (
	"context"
	"log"
	"time"

	"github.com/hyperledger/fabric-private-chaincode/clients/utils"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils/ledger"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils/pkg"
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
			log.Println(err)
			continue
		}

		for _, d := range dirs {
			vio, err := ipfs.GetUnreadViolations(ctx, sh, "/sla/"+d.Name)
			if err != nil {
				log.Println(err)
				break
			}
			for _, v := range vio {
				err = ledger.SLAViolation(Config, v)
				if err != nil {
					log.Println(err)
					break
				}
			}
		}
	}

}
