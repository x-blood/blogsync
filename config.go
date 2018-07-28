package main

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	Default *blogConfig
	Blogs   map[string]*blogConfig
}

type blogConfig struct {
	RemoteRoot string `yaml:"-"`
	LocalRoot  string `yaml:"local_root"`
	Username   string
	Password   string
	OmitDomain *bool `yaml:"omit_domain"`
}

func loadConfig(r io.Reader) (*config, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var blogs map[string]*blogConfig
	err = yaml.Unmarshal(bytes, &blogs)
	if err != nil {
		return nil, err
	}

	c := &config{
		Default: blogs["default"],
		Blogs:   blogs,
	}

	delete(blogs, "default")
	for key, b := range blogs {
		b.RemoteRoot = key
	}
	return c, nil
}

func (c *config) Get(remoteRoot string) *blogConfig {
	bc, ok := c.Blogs[remoteRoot]
	if !ok {
		return nil
	}

	if bc.LocalRoot == "" {
		bc.LocalRoot = c.Default.LocalRoot
	}

	return bc
}

func mergeBlogConfig(b1, b2 *blogConfig) *blogConfig {
	if b1 == nil {
		if b2 != nil {
			return b2
		}
		b1 = &blogConfig{}
	}
	if b2 == nil {
		return b1
	}
	if b1.LocalRoot == "" && b2.LocalRoot != "" {
		b1.LocalRoot = b2.LocalRoot
	}
	if b1.Username == "" && b2.Username != "" {
		b1.Username = b2.Username
	}
	if b1.Password == "" && b2.Password != "" {
		b1.Password = b2.Password
	}
	if b1.OmitDomain == nil && b2.OmitDomain != nil {
		b1.OmitDomain = b2.OmitDomain
	}
	return b1
}

func mergeConfig(c1, c2 *config) *config {
	if c1 == nil {
		c1 = &config{
			Blogs: make(map[string]*blogConfig),
		}
	}
	if c2 == nil {
		return c1
	}

	c1.Default = mergeBlogConfig(c1.Default, c2.Default)
	for k, bc := range c2.Blogs {
		c1.Blogs[k] = mergeBlogConfig(c1.Blogs[k], bc)
	}
	return c1
}
