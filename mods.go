package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var clusterRegexp = regexp.MustCompile(`\[\s*"workshop-(\d+)"\s*]\s*=\s*{`)
var serverRegexp = regexp.MustCompile(`[^(--)] *ServerModSetup *\(\s*"(\d+)"\s*\)`)

func getClusterMods(clusterDir, shard string, mods map[string]struct{}) error {
	f, err := os.Open(filepath.Join(clusterDir, shard, "modoverrides.lua"))
	if err != nil {
		return err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	for _, match := range clusterRegexp.FindAllStringSubmatch(string(buf), -1) {
		mods[match[1]] = struct{}{}
	}
	return nil
}

func getServerMods(setupPath string) ([]string, bool, error) {
	f, err := os.Open(setupPath)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, false, err
	}
	matches := serverRegexp.FindAllStringSubmatch(string(buf), -1)
	serverMods := make([]string, len(matches))
	for i := range serverMods {
		serverMods[i] = matches[i][1]
	}

	var noEOF bool
	if n := len(buf); n >= 1 && buf[n-1] != '\n' {
		noEOF = true
	}
	return serverMods, noEOF, nil
}

func appendServerMods(serverPath string, mods map[string]struct{}) error {
	setupPath := filepath.Join(filepath.Dir(serverPath), "../mods/dedicated_server_mods_setup.lua")
	serverMods, noEOF, err := getServerMods(setupPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, mod := range serverMods {
		delete(mods, mod)
	}
	if len(mods) == 0 {
		return nil
	}
	f, err := os.OpenFile(setupPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if noEOF {
		if _, err := fmt.Fprintf(f, "\n"); err != nil {
			return err
		}
	}
	for mod := range mods {
		if _, err := fmt.Fprintf(f, "ServerModSetup(\"%s\")\n", mod); err != nil {
			return err
		}
	}
	return nil
}
