package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

// A set of hetergeneous nodes
type Nodes struct {
	Type     string `json:"type" yaml:"type"`
	Make     string `json:"make" yaml:"make"`
	Model    string `json:"model" yaml:"model"`
	CPU      string `json:"cpu" yaml:"cpu"`
	CPUs     int    `json:"cpus" yaml:"cpus"`
	CPUCores int    `json:"cpucores" yaml:"cpucores"`
	GPU      string `json:"gpu,omitempty" yaml:"gpu"`
	GPUs     int    `json:"gpus,omitempty" yaml:"gpus"`
	GPUCores int    `json:"gpucores,omitempty" yaml:"gpucores"`
	RAM      int    `json:"ram" yaml:"ram"` // GB
	Count    int    `json:"count" yaml:"count"`
}

type Cluster struct {
	Name         string  `json:"name" yaml:"name"`
	Interconnect string  `json":interconnect" yaml":interconnect"`
	Nodes        []Nodes `json:"nodes" yaml:"nodes"`
}

func grantSummary(w io.Writer, c Cluster) error {
	return nil
}

func main() {
	flagFormat := flag.String("f", "json", "output format (json, yaml)")
	flagPretty := flag.Bool("pretty", false, "pretty-print json")
	flag.Parse()
	var (
		files []string
		err   error
	)
	if flag.NArg() < 1 {
		files, err = filepath.Glob("./*.yml")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files = flag.Args()
	}
	var all []Cluster
	for _, yml := range files {
		var c Cluster
		fin, err := os.Open(yml)
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.NewDecoder(fin).Decode(&c)
		fin.Close()
		if err != nil {
			log.Fatal(err)
		}
		all = append(all, c)
	}
	buf := &bytes.Buffer{}
	for _, c := range all {
		var b []byte
		if *flagFormat == "json" {
			if *flagPretty {
				b, err = json.MarshalIndent(c, "", "  ")
			} else {
				b, err = json.Marshal(c)
			}
		} else if *flagFormat == "yaml" {
			b, err = yaml.Marshal(c)
		}
		if err != nil {
			log.Fatal(err)
		}
		buf.Write(b)
	}
	fmt.Println(buf.String())
}
