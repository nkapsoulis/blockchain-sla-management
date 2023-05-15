package ipfs

import (
	"context"
	"io"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func createIndexFile(ctx context.Context, sh *shell.Shell, path string) error {
	indexPath := path + "/index"
	return sh.FilesWrite(ctx, indexPath, strings.NewReader(""), shell.FilesWrite.Create(true))
}

func writeToIndexFile(ctx context.Context, sh *shell.Shell, path, data string) error {
	indexPath := path + "/index"
	prevData, err := readIndexFile(ctx, sh, path)
	if err != nil {
		return err
	}
	return sh.FilesWrite(ctx, indexPath, strings.NewReader(string(prevData)+data+" "))
}

func overwriteIndexFile(ctx context.Context, sh *shell.Shell, path, data string) error {
	indexPath := path + "/index"
	err := sh.FilesRm(ctx, indexPath, false)
	if err != nil {
		return err
	}
	err = createIndexFile(ctx, sh, path)
	if err != nil {
		return err
	}
	return writeToIndexFile(ctx, sh, path, data)
}

func readIndexFile(ctx context.Context, sh *shell.Shell, path string) ([]byte, error) {
	indexPath := path + "/index"
	content, err := sh.FilesRead(ctx, indexPath)
	if err != nil {
		return []byte{}, err
	}

	data, err := io.ReadAll(content)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func parseIndexFile(ctx context.Context, sh *shell.Shell, path string) ([]string, error) {
	data, err := readIndexFile(ctx, sh, path)
	if err != nil {
		return []string{}, err
	}

	if string(data) == "" || string(data) == " " {
		return nil, nil
	}

	return strings.Split(strings.TrimSpace(string(data)), " "), nil
}
