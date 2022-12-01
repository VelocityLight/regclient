package main

import (
	"github.com/regclient/regclient/cmd/regsync-cloud/auth"
	"github.com/regclient/regclient/config"
	"os"
)

type NewConfig struct {
	CommonRepositories []Repository   `yaml:"commonRepositories" json:"commonRepositories,omitempty"`
	TargetRegistries   []Registry     `yaml:"targetRegistries" json:"targetRegistries,omitempty"`
	SourceRegistries   []Registry     `yaml:"sourceRegistries" json:"sourceRegistries,omitempty"`
	Defaults           ConfigDefaults `yaml:"defaults" json:"defaults"`
}

type Registry struct {
	Name    string        `yaml:"name" json:"name,omitempty"`
	Type    RegisteryType `yaml:"type" json:"type,omitempty"`
	Region  string        `yaml:"region" json:"region,omitempty"`
	Project string        `yaml:"project" json:"project,omitempty"`

	// TargetRegistries only
	Repositories []Repository `yaml:"repositories" json:"repositories,omitempty"`
}

type Repository struct {
	Name   string     `yaml:"name" json:"name"`
	Source string     `yaml:"source" json:"source"`
	Tags   ConfigTags `yaml:"tags" json:"tags"`
}

type RegisteryType string

const (
	ECR = "ECR"
	GCR = "GCR"
)

func (c NewConfig) ToConfig() *Config {
	res := ConfigNew()
	res.Defaults = c.Defaults
	for _, r := range c.SourceRegistries {
		res.Creds = append(res.Creds, config.Host{Name: r.Name})
	}
	for _, r := range c.TargetRegistries {
		crAuth, err := auth.NewEcrAuth(os.Getenv("AWS_ROLE_ARN"), r.Name)

		res.Creds = append(res.Creds,
			config.Host{Name: r.Name, Token: crAuth.GetToken()},
		)
		repos := append(c.CommonRepositories, r.Repositories...)

		if err != nil {
			log.Errorf("auth ecr failed: %s", err.Error())
			continue
		}
		for _, repo := range repos {
			// describe and create
			targetRepo, err := crAuth.CheckRepo(repo.Name)
			if err != nil {
				log.Errorf("check repo failed: %s", err.Error())
				continue
			}
			s := ConfigSync{
				Source: repo.Source,
				Target: targetRepo,
				Tags:   repo.Tags,
				Type:   "repository",
			}
			res.Sync = append(res.Sync, s)
		}
	}
	return res
}
