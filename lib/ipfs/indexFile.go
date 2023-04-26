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
	return sh.FilesWrite(ctx, indexPath, strings.NewReader(data+" "))
}

func overwriteIndexFile(ctx context.Context, sh *shell.Shell, path, data string) error {
	indexPath := path + "/index"
	return sh.FilesWrite(ctx, indexPath, strings.NewReader(data), shell.FilesWrite.Offset(0))
}

func parseIndexFile(ctx context.Context, sh *shell.Shell, path string) ([]string, error) {
	indexPath := path + "/index"
	content, err := sh.FilesRead(ctx, indexPath)
	if err != nil {
		return []string{}, err
	}

	data, err := io.ReadAll(content)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(data), " "), nil
}
