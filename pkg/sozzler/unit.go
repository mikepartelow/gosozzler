package sozzler

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

type unit struct {
	Name string `yaml:"name"`
}

//go:embed units.yaml
var unitsYAML []byte

var knownUnits map[string]struct{}

func init() {
	var units []unit
	if err := yaml.Unmarshal(unitsYAML, &units); err != nil {
		panic(err)
	}

	knownUnits = make(map[string]struct{})

	for _, u := range units {
		knownUnits[u.Name] = struct{}{}
	}
}
