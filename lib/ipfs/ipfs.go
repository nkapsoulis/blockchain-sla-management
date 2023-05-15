package ipfs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	iso19086 "github.com/hyperledger/fabric-private-chaincode/lib/iso-19086"
	shell "github.com/ipfs/go-ipfs-api"
)

type Metric struct {
	iso19086.Metrics
}

func CreateRootFolder(ctx context.Context, sh *shell.Shell) error {
	err := sh.FilesMkdir(ctx, "/sla", shell.FilesMkdir.Parents(true))

	if err != nil {
		return err
	}
	return nil
}

// Creates a folder for a new SLA add adds an index file. If SLA/folder exists, does nothing.
func CreateSLAFolder(ctx context.Context, sh *shell.Shell, id string) error {
	err := sh.FilesMkdir(ctx, "/sla/"+id, shell.FilesMkdir.Parents(true))

	if err != nil {
		return err
	}
	return createIndexFile(ctx, sh, "/sla/"+id)
}

func AddMetric(ctx context.Context, sh *shell.Shell, v iso19086.Metrics) error {
	path := "/sla/" + v.SLAID
	metricPath := path + "/" + v.ID

	metricBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = sh.FilesWrite(ctx, metricPath, strings.NewReader(string(metricBytes)), shell.FilesWrite.Create(true))
	if err != nil {
		return err
	}

	err = writeToIndexFile(ctx, sh, path, v.ID)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(ctx context.Context, sh *shell.Shell, path string) (string, error) {
	content, err := sh.FilesRead(ctx, path)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(content)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func GetUnreadMetrics(ctx context.Context, sh *shell.Shell, id string) ([]byte, error) {
	path := "/sla/" + id
	metricIDs, err := parseIndexFile(ctx, sh, path)
	log.Println("IDS:", metricIDs)
	if err != nil {
		return nil, fmt.Errorf("cannot read index file: %v", err)
	}

	var metrics []Metric

	if metricIDs == nil {
		return nil, nil
	}

	for _, vid := range metricIDs {
		data, err := ReadFile(ctx, sh, path+"/"+vid)
		if err != nil {
			return nil, fmt.Errorf("cannot read violation %v file: %v", vid, err)
		}
		var v Metric
		err = json.Unmarshal([]byte(data), &v)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal %v file: %v", v.ID, err)
		}
		metrics = append(metrics, v)
	}
	metricsBytes, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}
	return metricsBytes, nil
}

func OverwriteUnreadViolations(ctx context.Context, sh *shell.Shell, id, data string) error {
	path := "/sla/" + id
	err := overwriteIndexFile(ctx, sh, path, data)
	if err != nil {
		return err
	}
	return nil
}

func ReadAllMetrics(ctx context.Context, sh *shell.Shell, id string) ([]Metric, error) {
	files, err := sh.FilesLs(ctx, "/sla/"+id)
	if err != nil {
		return nil, err
	}

	var metrics []Metric
	for _, f := range files {
		if f.Name == "index" {
			continue
		}
		path := "/sla/" + id + "/" + f.Name
		d, err := ReadFile(ctx, sh, path)
		if err != nil {
			return nil, err
		}
		var v Metric
		err = json.Unmarshal([]byte(d), &v)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, v)
	}
	return metrics, nil
}
