package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/nleiva/yang-config-gen/compiler/junos"
	"github.com/nleiva/yang-config-gen/model"
)

func run(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("can't open file "+file, err)
	}

	defer f.Close()

	target := new(model.Target)

	err = ReadData(f, target)
	if err != nil {
		return "", fmt.Errorf("can't unmarshal data from "+file, err)
	}

	j := junos.NewCompiler()
	err = j.CompileConfig(*target)
	if err != nil {
		return "", fmt.Errorf("can't compile config from "+file, err)
	}

	output, err := j.EmitConfig()
	if err != nil {
		return "", fmt.Errorf("can't create json config from "+file, err)
	}
	return output, nil

}

func ReadData(r io.Reader, object any) error {
	d := json.NewDecoder(r)

	err := d.Decode(object)
	if err != nil {
		return fmt.Errorf("can't decode object: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("add a JSON input file to the command")
		os.Exit(1)
	}

	if result, err := run(os.Args[1]); err != nil {
		fmt.Printf("error generating the config: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("%v\n", result)
	}

}
