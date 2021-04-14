package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

func grantText(c Cluster) string {
	return ""
}

// A set of homogeneous nodes
type Nodes struct {
	Type     string `json:"type" yaml:"type"`
	Make     string `json:"make" yaml:"make"`
	Model    string `json:"model" yaml:"model"`
	CPU      string `json:"cpu" yaml:"cpu"`
	CPUs     int    `json:"cpus" yaml:"cpus"`
	CPUCores int    `json:"cpucores" yaml:"cpucores"`
	GPU      string `json:"gpu,omitempty" yaml:"gpu,omitempty"`
	GPUs     int    `json:"gpus,omitempty" yaml:"gpus,omitempty"`
	GPUCores int    `json:"gpucores,omitempty" yaml:"gpucores,omitempty"`
	RAM      int    `json:"ram" yaml:"ram"` // GB
	Count    int    `json:"count" yaml:"count"`
	Owner    string `json:"owner" yaml:"owner"`
	Added    string `json:"added" yaml:"added"`
}

type Cluster struct {
	Name         string  `json:"name" yaml:"name"`
	Interconnect string  `json":interconnect" yaml:"interconnect"`
	Storage      string  `json:"storage" yaml:"storage"`
	Desc         string  `json:"desc" yaml:"desc"`
	Nodes        []Nodes `json:"nodes" yaml:"nodes"`
}

func main() {
	flagFormat := flag.String("f", "json", "output format (json, yaml)")
	flagPretty := flag.Bool("pretty", false, "pretty-print json")
	flag.Parse()
	r, err := os.Open("./borah.yml")
	if err != nil {
		log.Fatal(err)
	}
	var c Cluster
	err = yaml.NewDecoder(r).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
	var b []byte
	switch *flagFormat {
	case "json":
		if *flagPretty {
			b, err = json.MarshalIndent(c, "", "  ")
		} else {
			b, err = json.Marshal(c)
		}
	case "yaml":
		b, err = yaml.Marshal(c)
	default:
		log.Fatalf("invalid format: %s", *flagFormat)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
