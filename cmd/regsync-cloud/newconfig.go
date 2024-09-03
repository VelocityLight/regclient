package main

import (
	"fmt"
	"os"

	"github.com/regclient/regclient/cmd/regsync-cloud/auth"
	"github.com/regclient/regclient/config"
)

type NewConfig struct {
	RepositorySets   map[string]RepositorySet `yaml:"repositorySets" json:"repositorySets,omitempty"`
	TargetRegistries []Registry               `yaml:"targetRegistries" json:"targetRegistries,omitempty"`
	Creds            []config.Host            `yaml:"creds" json:"creds"`
	Defaults         ConfigDefaults           `yaml:"defaults" json:"defaults"`
}

type Registry struct {
	Name    string        `yaml:"name" json:"name,omitempty"`
	Type    RegisteryType `yaml:"type" json:"type,omitempty"`
	Region  string        `yaml:"region" json:"region,omitempty"`
	Project string        `yaml:"project" json:"project,omitempty"`

	// TargetRegistries only
	IncludeRepositorySets []string     `yaml:"includeRepositorySets" json:"includeRepositorySets,omitempty"`
	Repositories          []Repository `yaml:"repositories" json:"repositories,omitempty"`
}

type Repository struct {
	Name   string     `yaml:"name" json:"name"`
	Source string     `yaml:"source" json:"source"`
	Tags   ConfigTags `yaml:"tags" json:"tags"`
}

type RepositorySet []Repository

type RegisteryType string

const (
	ECR = "ECR"
	GCR = "GCR"
	GAR = "GAR"
	ACR = "ACR"
)

func (c NewConfig) auth(r Registry, cfg *Config) (auth.Authenticator, error) {
	var crAuth auth.Authenticator
	var err error

	h := config.Host{
		Name: r.Name,
	}

	switch r.Type {
	case ECR:
		crAuth, err = auth.NewEcrAuth(os.Getenv("AWS_ROLE_ARN"), r.Name)
		if err != nil {
			log.Errorf("auth ecr failed: %s", err.Error())
			return nil, err
		}
	case GAR:
		crAuth, err = auth.NewGarAuth(r.Name)
		if err != nil {
			log.Errorf("auth ecr failed: %s", err.Error())
			return nil, err
		}
	case GCR:
		crAuth, _ = auth.NewGcrAuth(r.Name, r.Project)
		h.RepoAuth = true
	case ACR:
		crAuth, _ = auth.NewAcrAuth(r.Name)
		h.RepoAuth = true
	default:
		crAuth, _ = auth.NewCommonCRAuth(r.Name)
	}

	//u, p := crAuth.GetLoginPassword()
	//h.Token = crAuth.GetToken()
	//h.User = u
	//h.Pass = p

	for _, cred := range cfg.Creds {
		if h.Name == cred.Name {
			return crAuth, nil
		}
	}

	cfg.Creds = append(cfg.Creds, h)

	return crAuth, nil
}

func (c NewConfig) ToConfig() *Config {
	res := ConfigNew()
	res.Defaults = c.Defaults

	res.Creds = append(res.Creds, c.Creds...)

	for _, r := range c.TargetRegistries {
		crAuth, err := c.auth(r, res)
		if err != nil {
			continue
		}

		var repos []Repository
		for _, s := range r.IncludeRepositorySets {
			r, ok := c.RepositorySets[s]
			if ok {
				repos = append(repos, r...)
			}
		}

		repos = append(repos, r.Repositories...)

		for _, repo := range repos {
			// describe and create
			targetRepo, err := crAuth.CheckRepo(repo.Name)
			if err != nil {
				log.Errorf("check repo %s failed: %s", repo, err.Error())
				continue
			}

			s := ConfigSync{}
			if repo.Tags.Deny == nil && len(repo.Tags.Allow) == 1 {
				s.Source = fmt.Sprintf("%s:%s", repo.Source, repo.Tags.Allow[0])
				s.Target = fmt.Sprintf("%s:%s", targetRepo, repo.Tags.Allow[0])
				s.Type = "image"
			} else {
				s = ConfigSync{
					Source: repo.Source,
					Target: targetRepo,
					Tags:   repo.Tags,
					Type:   "repository",
				}
			}

			res.Sync = append(res.Sync, s)
		}
	}
	return res
}
