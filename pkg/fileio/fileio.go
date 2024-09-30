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
const rootdir = "/mpk"

var (
	cfgdir string
)

func init() {
	usercfg, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	cfgdir = usercfg + rootdir
	if err := preparedir(cfgdir); err != nil {
		panic(err)
	}

}

func LoadGroup() ([]string, error) {
	gps, err := files(cfgdir)
	if err != nil {
		return nil, err
	}
	return gps, nil
}

func Load[T any](group string) (map[string]T, error) {
	fns, err := files(filepath.Join(cfgdir, group))
	if err != nil {
		return nil, err
	}
	fs := map[string]T{}
	for _, fn := range fns {
		bz, err := os.ReadFile(prepFilePath(group, fn))
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

func CreateGroup(group string) error {
	return preparedir(filepath.Join(cfgdir, group))
}

func Open(group, name string) (*os.File, error) {
	f, err := os.Open(prepFilePath(group, name))
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

func Save[T any](group, name string, data T) error {
	bz, err := json.Marshal(data)
	if err != nil {
		return err
	}

	f, err := preparefile(prepFilePath(group, name))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(bz); err != nil {
		return err
	}
	return nil
}

// preparefile returns file handler, otherwise it returns err if file exist or dir not exist
func preparefile(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return f, nil
}

// preparedir checks dir, returns err if exist
// creates the dir otherwise
func preparedir(dir string) error {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0700); err != nil {
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

	return farr[1:], nil
}

func prepFilePath(group, name string) string {
	return filepath.Join(cfgdir, group, name+subfix)
}
