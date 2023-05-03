package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
)

const metricFolder = "./metrics/"
const slaFolder = "./slas/"

func main() {
	slaData, err := ioutil.ReadFile(slaFolder + "Incident Resolution Time SLO.json")
	if err != nil {
		panic(err)
	}
	sla, err := iso19086.ReadSLA(slaData)
	log.Println(sla)

	if err != nil {
		panic(err)
	}

	dirs, err := os.ReadDir(metricFolder)
	if err != nil {
		panic(err)
	}
	for _, d := range dirs {
		data, err := ioutil.ReadFile(metricFolder + d.Name())
		if err != nil {
			panic(err)
		}
		fmt.Println(d.Name())
		m, err := iso19086.ReadMetric(data)
		if err != nil {
			panic(err)
		}
		log.Println(m)
		violated, err := iso19086.Parse(slaData, data)
		if err != nil {
			panic(err)
		}
		fmt.Println("Violated: ", violated)
	}
}
