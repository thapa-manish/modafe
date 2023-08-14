package main

import (
	"encoding/json"
	"fmt"
	"io"
	"modafe/pkg/decoder"
	"modafe/pkg/encoder"
	"modafe/pkg/types"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		panic(fmt.Errorf("please pass the schema json file as parameter"))
	}
	file, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	schema, err := encoder.NewEncoder().EncodeString(string(data))
	if err != nil {
		panic(err)
	}

	var output = make(types.JSON)
	decoder.NewDecoder().Decode(schema, output)
	outData, err := json.MarshalIndent([]types.JSON{output}, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(outData))
}
