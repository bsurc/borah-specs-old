package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/go-yaml/yaml"
)

func grantText(w io.Writer, c Cluster) error {
	t := template.Must(template.New("grant").ParseFiles("templates/grant.txt"))
	type txt struct {
		ComputeCount int // computed
		ComputeCores int // computed
		ComputeMem   int
		GPUCount     int // computed
		GPUCores     int // computed
		GPUMem       int
		HighMemCount int
		HighMemMem   int
	}
	x := txt{
		ComputeMem:   192,
		GPUMem:       384,
		HighMemCount: 1,
		HighMemMem:   768,
	}

	for _, node := range c.Nodes {
		// Only Boise State and Condo nodes
		if node.Owner == "Idaho Power Company" {
			continue
		}
		switch node.Type {
		case "Compute":
			x.ComputeCount += node.Count
			x.ComputeCores += node.CPUCores * node.Count
		case "GPU":
			x.GPUCount += node.Count
			x.GPUCores += node.GPUCores * node.Count
		}
	}
	return t.ExecuteTemplate(w, "grant", x)
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flagFormat := flag.String("f", "json", "output format (json, yaml)")
	flagPretty := flag.Bool("pretty", false, "pretty-print json")
	flagTemplate := flag.String("t", "", "execute a template")
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
	switch *flagTemplate {
	case "":
		break
	case "templates/grant.txt", "grant.txt", "grant":
		err := grantText(os.Stdout, c)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
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
