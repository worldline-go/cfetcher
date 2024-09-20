package file

import (
	"io/fs"
	"strconv"
)

type options struct {
	filePerm   *fs.FileMode
	folderPerm *fs.FileMode
	fileFlag   *int
}

type Option func(opt *options) error

func (o *options) apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}

	return nil
}

// WithFilePerm sets the file permission.
//
// The permission is a string of octal digits, such as "0644".
func WithFilePerm(filePerm string) Option {
	return func(options *options) error {
		if filePerm == "" {
			return nil
		}

		perm, err := strconv.ParseUint(filePerm, 8, 32)
		if err != nil {
			return err
		}

		v := fs.FileMode(perm)
		options.filePerm = &v

		return nil
	}
}

// WithFolderPerm sets the folder permission.
//
// The permission is a string of octal digits, such as "0755".
func WithFolderPerm(folderPerm string) Option {
	return func(options *options) error {
		if folderPerm == "" {
			return nil
		}

		perm, err := strconv.ParseUint(folderPerm, 8, 32)
		if err != nil {
			return err
		}

		v := fs.FileMode(perm)
		options.folderPerm = &v

		return nil
	}
}
