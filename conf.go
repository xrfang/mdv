package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type (
	config struct {
		CSS  string `json:"css"`
		Quit int    `json:"quit"`
		dir  string
	}
)

func (c *config) Load() error {
	if c.dir == "" {
		return nil
	}
	f, err := os.Open(filepath.Join(c.dir, "config.json"))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(c)
}

func (c config) Save() (err error) {
	if c.dir == "" {
		return nil
	}
	defer func() {
		if e := recover(); e != nil {
			err = trace("%v", e)
		}
	}()
	f, err := os.Create(filepath.Join(c.dir, "config.json"))
	assert(err)
	defer func() { assert(f.Close()) }()
	return json.NewEncoder(f).Encode(c)
}

var cf config

func init() {
	cf.CSS = "default.css"
	cf.Quit = 9
	dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	dir = filepath.Join(dir, "mdv")
	if err := os.MkdirAll(dir, 0700); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	cf.dir = dir
	if err := cf.Load(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		err = cf.Save()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
