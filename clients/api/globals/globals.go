package globals

import (
	"math/rand"

	"github.com/hyperledger/fabric-private-chaincode/clients/utils/pkg"
)

var Config *pkg.Config
var Secret = []byte(randSeq(32))

const Userkey = "user"
const AppName = "fabric-private-chaincode"
const Passphrase = "password"
const SessionName = "fabricAuth"

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
