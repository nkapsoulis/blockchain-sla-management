package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-private-chaincode/lib/ipfs"
	iso19086parser "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"

	shell "github.com/ipfs/go-ipfs-api"
)

var names []string

// var providers []iso19086parser.Entity
// var clients []iso19086parser.Entity

const IPFSHost = "localhost:5001"
const nViolations = 10

func main() {
	rand.Seed(time.Now().UnixNano())
	data, err := os.ReadFile("./first-names.json")
	if err != nil {
		log.Fatalf("failed to read names file: %v", err)
	}
	err = json.Unmarshal(data, &names)
	if err != nil {
		log.Fatalf("failed to unmarshal files: %v", err)
	}

	// providers, _ = createUsers(0, len(names))
	// clients, _ = createUsers(0, len(names))

	metrics := createMetrics(nViolations, nViolations)

	ctx := context.Background()
	sh := shell.NewShell(IPFSHost)

	err = ipfs.CreateRootFolder(ctx, sh)
	if err != nil {
		panic(err)
	}

	for _, metric := range metrics {
		err = ipfs.CreateSLAFolder(ctx, sh, metric.SLAID)
		if err != nil {
			panic(err)
		}

		err = ipfs.AddMetric(ctx, sh, metric)
		if err != nil {
			panic(err)
		}
		log.Println("Wrote metric")

		time.Sleep(5 * time.Second)
	}
}

func CreateFile(metric iso19086parser.Metrics) {
	f, err := os.Create("./metrics/" + metric.ID + ".json")
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}
	f.Write(data)
	f.Close()
}

// func CreateAssets(nAssets int) []lib.SLA {
// 	states := []string{"started", "ongoing"} // , "stopped", "deleted"}
// 	types := []string{"agreement"}

// 	assets := make([]lib.SLA, nAssets)
// 	for i := 0; i < nAssets; i++ {
// 		nProvider := rand.Intn(len(providers))
// 		nClient := rand.Intn(len(clients))

// 		id := fmt.Sprintf("a%d", i)
// 		name := fmt.Sprintf("Agreement %d", i)
// 		importance := []lib.Importance{
// 			{Name: "Warning", Constraint: "> 30"},
// 			{Name: "Mild", Constraint: "> 30"},
// 			{Name: "Serious", Constraint: "> 30"},
// 			{Name: "Sever", Constraint: "> 70"},
// 			{Name: "Catastrophic", Constraint: "> 70"},
// 		}
// 		asset := iso19086parser.SLA{
// 			ReferenceID: id, Name: name, State: states[rand.Intn(len(states))],
// 			Assessment: lib.Assessment{FirstExecution: time.Now().Add(-1000 * time.Hour).Format(time.RFC3339),
// 				LastExecution: time.Now().Format(time.RFC3339)},
// 			Details: lib.Detail{
// 				ID:       id,
// 				Type:     types[rand.Intn(len(types))],
// 				Name:     name,
// 				Provider: providers[nProvider],
// 				Client:   clients[nClient],
// 				Creation: time.Now().Format(time.RFC3339),
// 				Guarantees: []lib.Guarantee{{Name: "TestGuarantee", Constraint: "[test_value] < 0.7", Importance: []lib.Importance{}},
// 					{Name: "TestGuarantee2", Constraint: "[test_value] < 0.2", Importance: importance}},
// 				Service: "8",
// 			},
// 		}
// 		assets[i] = asset
// 	}
// 	return assets
// }

func createMetrics(nViolations, nAssets int) []iso19086parser.Metrics {
	if nAssets == 0 {
		nAssets = 5
	}

	metrics := make([]iso19086parser.Metrics, nViolations)
	for i := 0; i < nViolations; i++ {
		metric := iso19086parser.Metrics{
			ID:    fmt.Sprintf("v%d", i),
			SLAID: "a1",
			Sample: iso19086parser.SampleData{
				IncidentReportTime:     strconv.Itoa(rand.Int()),
				IncidentResolutionTime: strconv.Itoa(rand.Int()),
				IncidentResponseTime:   strconv.Itoa(rand.Int()),
			},
		}
		metrics[i] = metric
	}

	return metrics
}

func createUsers(startID, nUsers int) ([]iso19086parser.Entity, int) {
	var users []iso19086parser.Entity
	var id int
	for id = startID; id <= startID+nUsers-1; id++ {
		users = append(users, iso19086parser.Entity{ID: strconv.Itoa(id), Name: names[id]})
	}
	return users, id
}
