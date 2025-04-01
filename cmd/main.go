package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/nleiva/yang-config-gen/compiler/junos"
	"github.com/nleiva/yang-config-gen/model"
)

func check(s string, err error) {
	if err != nil {
		err = fmt.Errorf("%s: %w", s, err)
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("add a JSON input file to the command")
		os.Exit(1)
	}

	file := cmp.Or(os.Args[1], "../model/testdata/interface.json")

	f, err := os.Open(file)
	check("can't open file "+file, err)

	defer f.Close()

	target := new(model.Target)

	err = ReadData(f, target)
	check("can't unmarshal data from "+file, err)

	j := junos.NewCompiler()
	err = j.CompileConfig(*target)
	check("can't compile config from "+file, err)

	output, err := j.EmitConfig()
	check("can't create json config from "+file, err)

	fmt.Printf("%v\n", output)

}

func ReadData(r io.Reader, object any) error {
	d := json.NewDecoder(r)

	err := d.Decode(object)
	if err != nil {
		return fmt.Errorf("can't decode object: %w", err)
	}
	return nil
}
