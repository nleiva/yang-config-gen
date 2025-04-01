package model

import (
	"encoding/json"
	"fmt"
	"io"
)

func ReadData(r io.Reader, object any) error {
	d := json.NewDecoder(r)

	err := d.Decode(object)
	if err != nil {
		return fmt.Errorf("can't decode object: %w", err)
	}
	return nil
}
