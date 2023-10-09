package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"go.xrfang.cn/yal"
)

type (
	config struct {
		MainCSS string `json:"main_css"`
		CodeCSS string `json:"code_css"`
		Recurse int    `json:"recurse"`
		Port    int    `json:"port"`
		Rev     int    `json:"rev"`
		dir     string
	}
)

func (c *config) Load() (err error) {
	defer yal.Catch(func(e error) error {
		err = e
		if os.IsNotExist(err) {
			return nil
		}
		return err
	})
	cfg := filepath.Join(c.dir, "config.json")
	f, err := os.Open(cfg)
	yal.Assert(err)
	defer f.Close()
	yal.Assert(json.NewDecoder(f).Decode(c))
	rev, _ := strconv.Atoi(_G_REVS)
	if rev > c.Rev {
		os.Remove(cfg)
		os.Remove(filepath.Join(c.dir, "default.css"))
		os.Remove(filepath.Join(c.dir, "highlight.css"))
		return errors.New("version upgraded, reset config")
	}
	return
}

func (c config) Save() (err error) {
	defer yal.Catch(&err)
	f, err := os.Create(filepath.Join(c.dir, "config.json"))
	yal.Assert(err)
	defer func() { yal.Assert(f.Close()) }()
	c.Rev, _ = strconv.Atoi(_G_REVS)
	return json.NewEncoder(f).Encode(c)
}

var cf config

func init() {
	cf.MainCSS = "default.css"
	cf.CodeCSS = "highlight.css"
	cf.Recurse = 2
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
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, err)
		}
		err = cf.Save()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
