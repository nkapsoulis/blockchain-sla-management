package ipfs

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/hyperledger/fabric-private-chaincode/lib"
	shell "github.com/ipfs/go-ipfs-api"
)

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

func AddViolation(ctx context.Context, sh *shell.Shell, v lib.Violation) error {
	path := "/sla/" + v.SLAID + "/" + v.ID

	violationBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = sh.FilesWrite(ctx, path, strings.NewReader(string(violationBytes)))
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

func GetUnreadViolations(ctx context.Context, sh *shell.Shell, id string) ([]lib.Violation, error) {
	path := "/sla/" + id
	violationIDs, err := parseIndexFile(ctx, sh, path)
	if err != nil {
		return nil, err
	}

	var violations []lib.Violation
	for _, vid := range violationIDs {
		data, err := ReadFile(ctx, sh, path+"/"+vid)
		if err != nil {
			return nil, err
		}
		var v lib.Violation
		err = json.Unmarshal([]byte(data), &v)
		if err != nil {
			return nil, err
		}
		violations = append(violations, v)
	}

	err = overwriteIndexFile(ctx, sh, path, "")
	if err != nil {
		return nil, err
	}
	return violations, nil
}

func ReadAllViolations(ctx context.Context, sh *shell.Shell, id string) ([]lib.Violation, error) {
	files, err := sh.FilesLs(ctx, "/sla/"+id)
	if err != nil {
		return nil, err
	}

	var violations []lib.Violation
	for _, f := range files {
		if f.Name == "index" {
			continue
		}
		path := "/sla/" + id + "/" + f.Name
		d, err := ReadFile(ctx, sh, path)
		if err != nil {
			return nil, err
		}
		var v lib.Violation
		err = json.Unmarshal([]byte(d), &v)
		if err != nil {
			return nil, err
		}
		violations = append(violations, v)
	}
	return violations, nil

}
