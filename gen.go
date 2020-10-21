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
	"sort"

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
	Owner    string `json:"owner" yaml:"owner"`
}

type Cluster struct {
	Name         string  `json:"name" yaml:"name"`
	Interconnect string  `json":interconnect" yaml:"interconnect"`
	Storage      string  `json:"storage" yaml:"storage"`
	Nodes        []Nodes `json:"nodes" yaml:"nodes"`
}

func main() {
	flagFormat := flag.String("f", "json", "output format (json, yaml)")
	flagPretty := flag.Bool("pretty", false, "pretty-print json")
	flag.Parse()
	var (
		files []string
		err   error
	)
	files, err = filepath.Glob("./*.yml")
	if err != nil {
		log.Fatal(err)
	}
	sort.Slice(files, func(i, j int) bool {
		if filepath.Base(files[i]) == "borah.yml" {
			return true
		} else if filepath.Base(files[j]) == "borah.yml" {
			return false
		}
		return files[i] < files[j]
	})

	if files[0] != "borah.yml" {
		panic("failed to sort files with borah first")
	}

	readers := make([]io.Reader, len(files))
	for i, f := range files {
		fin, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		readers[i] = fin
		defer fin.Close()
	}
	r := io.MultiReader(readers...)
	var c Cluster
	err = yaml.NewDecoder(r).Decode(&c)
	buf := &bytes.Buffer{}
	for _, c := range c.Nodes {
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
