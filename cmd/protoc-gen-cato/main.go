package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/ncuhome/cato/src"
)

func main() {

	protoInput, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("[-] cato read data from stdin: %#v", err)
	}
	pbRequest := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(protoInput, pbRequest); err != nil {
		log.Fatalf("[-] cato unmarshal pbRequest data: %#v", err)
	}
	pbResponse := new(pluginpb.CodeGeneratorResponse)
	generator := src.NewDBGenerator(pbRequest)
	generator.Generate(pbResponse)
	output, err := proto.Marshal(pbResponse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(output)
}
