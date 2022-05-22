package sync

import "io/ioutil"

type FilePathSync struct {
}

func (fs *FilePathSync) Fetch(input string, strategy SYNC_STRATEGY) (string, error) {

	rawFile, err := ioutil.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(rawFile), nil
}
