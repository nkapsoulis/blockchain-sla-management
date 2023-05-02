package ledger

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-private-chaincode/clients/utils"
	"github.com/hyperledger/fabric-private-chaincode/clients/utils/pkg"
	sla_types "github.com/hyperledger/fabric-private-chaincode/clients/utils/types"

	"github.com/tyler-smith/go-bip32"
)

func InitLedger(config *pkg.Config, passphrase string) {
	client := pkg.NewClient(config)

	_, err := client.Invoke("InitLedger")
	if err != nil {
		if err.Error() == "init has already ran" {
			log.Println("Init Ledger has already run. Continuing.")
			return
		}
		log.Fatalln(err)
		return
	}
	users := [10]string{"Tomoko", "Brad", "Jin Soo", "Max", "Adriana", "Michel", "Mario", "Anton", "Marek", "George"}
	for _, u := range users {
		mnemonic, err := utils.CreateMnemonic()
		if err != nil {
			log.Fatalln(err)
		}
		keysSerialized, err := utils.CreateMasterKey(mnemonic, passphrase)
		if err != nil {
			log.Fatalln(err)
		}
		keys, err := bip32.B58Deserialize(keysSerialized)
		if err != nil {
			log.Fatalln(err)
		}
		CreateUser(config, u, keys.PublicKey().B58Serialize())
		fmt.Println(u, mnemonic)
	}

}

func GetUser(config *pkg.Config, name string) sla_types.User {
	client := pkg.NewClient(config)
	res, err := client.Query("ReadUser", name)
	if err != nil {
		log.Fatalln(err)
	}

	var user sla_types.User
	json.Unmarshal([]byte(res), &user)
	return user
}

func CreateUser(config *pkg.Config, name, publicKey string) {
	client := pkg.NewClient(config)
	_, err := client.Invoke("CreateUser", name, publicKey, "500")
	if err != nil {
		log.Fatalln(err)
	}
}

func GetSLA(config *pkg.Config, id string) (sla_types.SLA, error) {
	client := pkg.NewClient(config)
	res, err := client.Query("ReadSLA", id)
	if err != nil {
		return sla_types.SLA{}, err
	}

	var sla sla_types.SLA
	err = json.Unmarshal([]byte(res), &sla)
	if err != nil {
		return sla_types.SLA{}, err
	}
	return sla, nil
}

func CreateSLA(config *pkg.Config, sla sla_types.SLA) error {
	client := pkg.NewClient(config)

	slaJson, err := json.Marshal(sla)
	if err != nil {
		return err
	}
	_, err = client.Invoke("CreateOrUpdateContract", string(slaJson))
	if err != nil {
		return err
	}
	return nil
}

func GetSLAApproval(config *pkg.Config, id string) (sla_types.Approval, error) {
	client := pkg.NewClient(config)
	res, err := client.Query("GetApprovals", id)
	if err != nil {
		return sla_types.Approval{}, err
	}

	var approval sla_types.Approval
	err = json.Unmarshal([]byte(res), &approval)
	if err != nil {
		return sla_types.Approval{}, err
	}
	return approval, nil
}

func Approve(config *pkg.Config, id, username string, signature []byte) error {
	client := pkg.NewClient(config)
	_, err := client.Invoke("Approve", id, username, hex.EncodeToString(signature))
	if err != nil {
		return err
	}

	return nil
}

func CheckForViolation(config *pkg.Config, metric sla_types.Metric) error {
	client := pkg.NewClient(config)
	fmt.Printf("Received metric: %v\n", metric)

	violationJSON, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	_, err = client.Invoke("SLAViolated", string(violationJSON))
	if err != nil {
		return err
	}
	return nil
}
