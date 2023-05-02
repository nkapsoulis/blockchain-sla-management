package main

import (
	"fmt"
	"io/ioutil"
	"os"

	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
)

const testFolder = "./slas/"

func main() {
	dirs, err := os.ReadDir(testFolder)
	if err != nil {
		panic(err)
	}
	for _, d := range dirs {
		data, err := ioutil.ReadFile(testFolder + d.Name())
		if err != nil {
			panic(err)
		}
		fmt.Println(d.Name())
		sla, err := iso19086.ReadMetric(data)
		if err != nil {
			panic(err)
		}

		iso19086.ParseMetrics(*sla)
	}
}
