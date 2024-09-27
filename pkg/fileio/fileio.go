package fileio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const subfix = ".json"

var (
	cfgdir string
)

func init() {
	usercfg, _ := os.UserConfigDir()
	cfgdir = usercfg + "/multi-signer/"
	if err := preparedir(cfgdir); err != nil {
		panic(err)
	}

}

func Load[T any]() (map[string]T, error) {
	fns, err := files(cfgdir)
	if err != nil {
		return nil, err
	}
	fs := map[string]T{}
	for _, fn := range fns {
		bz, err := os.ReadFile(prepFilePath(fn))
		if err != nil {
			return nil, err
		}
		stru := new(T)
		err = json.Unmarshal(bz, stru)
		if err != nil {
			return nil, err
		}
		fs[fn] = *stru
	}
	return fs, nil
}

func Open(name string) (*os.File, error) {
	f, err := os.Open(prepFilePath(name))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Close(fs ...*os.File) error {
	errs := []error{}
	for _, f := range fs {
		if err := f.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func Save[T any](name string, data T) error {
	bz, err := json.Marshal(data)
	if err != nil {
		return err
	}

	f, err := preparefile(prepFilePath(name))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(bz); err != nil {
		return err
	}
	return nil
}

// preparedir checks the file, returns file and err if exist
// creates the dir and file otherwise
func preparefile(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return f, nil
}

func preparedir(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(cfgdir, 0700); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func files(dir string) ([]string, error) {
	farr := []string{}

	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		farr = append(farr, info.Name())
		return nil
	}); err != nil {
		return nil, err
	}

	return farr, nil
}

func prepFilePath(name string) string {
	return cfgdir + name + subfix
}
