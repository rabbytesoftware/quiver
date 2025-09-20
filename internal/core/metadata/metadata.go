package metadata

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed metadata.yaml
	metadataByte []byte
	metadata *Metadata
)

type Maintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
	URL   string `yaml:"url"`
}

type Metadata struct {
	Version     string       `yaml:"version"`
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Author      string       `yaml:"author"`
	URL         string       `yaml:"url"`
	License     string       `yaml:"license"`
	Copyright   string       `yaml:"copyright"`
	Maintainers []Maintainer `yaml:"maintainers"`
}

func Get() (*Metadata) {
	if metadata != nil {
		return metadata
	}
	
	err := yaml.Unmarshal(metadataByte, metadata)
	if err != nil {
		// Fallback to default metadata
		return &Metadata{
			Version:     "25.9.0",
			Name:        "Quiver",
			Description: "The future of wizards and package managers.",
			Author:      "Rabbyte Software",
			URL:         "https://quiver.ar",
			License:     "GPL-3.0",
			Copyright:   "Copyright 2025 Rabbyte Software & char2cs.net",
			Maintainers: []Maintainer{
				{
					Name:  "Mateo Urrutia",
					Email: "me@char2cs.net",
					URL:   "https://char2cs.net",
				},
			},
		}
	}

	return metadata
}

func GetVersion() (string) {
	return Get().Version
}

func GetName() (string) {
	return Get().Name
}

func GetDescription() (string) {
	return Get().Description
}

func GetAuthor() (string) {
	return Get().Author
}

func GetURL() (string) {
	return Get().URL
}

func GetLicense() (string) {
	return Get().License
}

func GetCopyright() (string) {
	return Get().Copyright
}

func GetMaintainers() ([]Maintainer) {
	return Get().Maintainers
}
