package file

import (
	"fmt"
	"os"

	"github.com/worldline-go/struct2"
	"gopkg.in/yaml.v3"
)

var mapDecoder = struct2.Decoder{
	TagName: "cfg",
}

func Load(fileName string, v interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	var vMap interface{}
	if err := yaml.NewDecoder(file).Decode(&vMap); err != nil {
		return fmt.Errorf("failed to decode yaml file: %w", err)
	}

	return mapDecoder.Decode(vMap, v)
}
