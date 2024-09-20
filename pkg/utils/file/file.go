package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	defaultFileFlag   int         = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	defaultFilePerm   fs.FileMode = 0o644
	defaultFolderPerm fs.FileMode = 0o755
)

func OpenFileWrite(path string, opts ...Option) (*os.File, error) {
	opt := &options{}
	if err := opt.apply(opts...); err != nil {
		return nil, err
	}

	// create folder if not exist
	folder := filepath.Dir(path)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		folderPerm := defaultFolderPerm
		if opt.folderPerm != nil {
			folderPerm = *opt.folderPerm
		}

		if err := os.MkdirAll(folder, folderPerm); err != nil {
			return nil, fmt.Errorf("failed to create folder %s: %w", folder, err)
		}
	}

	filePerm := defaultFilePerm
	if opt.filePerm != nil {
		filePerm = *opt.filePerm
	}

	fileFlag := defaultFileFlag
	if opt.fileFlag != nil {
		fileFlag = *opt.fileFlag
	}

	// open file
	f, err := os.OpenFile(path, fileFlag, filePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	return f, nil
}

func SaveFile(path string, value []byte) error {
	return nil
}
