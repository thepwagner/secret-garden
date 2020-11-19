package action

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
	"github.com/thepwagner/secret-garden/token"
	"gopkg.in/yaml.v3"
)

type ConfigProcessor struct {
	m   token.Minter
	s   token.Storer
	org string
}

func NewConfigProcessor(org string, minter token.Minter, storer token.Storer) *ConfigProcessor {
	return &ConfigProcessor{org: org, m: minter, s: storer}
}

func (c *ConfigProcessor) ProcessConfig(ctx context.Context, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !strings.HasSuffix(path, "secrets.yaml") {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		var raw map[string]interface{}
		if err := yaml.NewDecoder(f).Decode(&raw); err != nil {
			return err
		}

		org := filepath.Dir(path) == dir
		if org {
			return c.processOrgSecrets(ctx, raw)
		}
		repo := filepath.Base(filepath.Dir(path))
		return c.processRepoSecrets(ctx, repo, raw)
	})
}

func (c *ConfigProcessor) processOrgSecrets(ctx context.Context, secrets map[string]interface{}) error {
	for secretName, secretValue := range secrets {
		secretCfg := secretValue.(map[string]interface{})
		log := logrus.WithField("name", secretName)
		log.Debug("processing org secret")

		tok, err := c.mintToken(ctx, secretCfg)
		if err != nil {
			return err
		}

		consumers := stringList(secretCfg, "consumers")
		log.WithField("consumers", consumers).Debug("storing secret")
		if err := c.s.StoreOrg(ctx, c.org, secretName, tok, consumers); err != nil {
			return err
		}
	}
	return nil
}

func (c *ConfigProcessor) processRepoSecrets(ctx context.Context, repo string, secrets map[string]interface{}) error {
	for secretName, secretValue := range secrets {
		secretCfg := secretValue.(map[string]interface{})
		log := logrus.WithField("name", secretName)
		log.Debug("processing org secret")

		tok, err := c.mintToken(ctx, secretCfg)
		if err != nil {
			return err
		}

		log.Debug("storing secret")
		if err := c.s.StoreRepo(ctx, c.org, repo, secretName, tok); err != nil {
			return err
		}
	}
	return nil
}

func (c *ConfigProcessor) mintToken(ctx context.Context, secretCfg map[string]interface{}) (string, error) {
	perms, err := secretPermissions(secretCfg)
	if err != nil {
		return "", err
	}
	targets := stringList(secretCfg, "targets")

	tok, err := c.m.Mint(ctx, targets, perms)
	if err != nil {
		return "", fmt.Errorf("minting token: %w", err)
	}
	return tok, nil
}

func stringList(secretCfg map[string]interface{}, key string) (ret []string) {
	if t, ok := secretCfg[key].([]interface{}); ok {
		for _, target := range t {
			if s, ok := target.(string); ok {
				ret = append(ret, s)
			}
		}
	}
	return
}

func secretPermissions(secretCfg map[string]interface{}) (*github.InstallationPermissions, error) {
	var perms github.InstallationPermissions
	if p, ok := secretCfg["permissions"].(map[string]interface{}); ok {
		for k, v := range p {
			s, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("unknown permission value, expected 'read' or 'write': %v", k)
			}

			switch k {
			case "contents":
				perms.Contents = &s
			case "issues":
				perms.Issues = &s
			default:
				return nil, fmt.Errorf("unknown permission: %v", k)
			}
		}

		if v, ok := p["contents"].(string); ok {
			perms.Contents = &v
		}
		if v, ok := p["issues"].(string); ok {
			perms.Issues = &v
		}
	}
	return &perms, nil
}
