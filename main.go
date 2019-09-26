// Copyright (c) 2019 Fadhli Dzil Ikram. All rights reserved.
// This source code is governed by a MIT license. See LICENSE for more info.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type options struct {
	PersistentStorageRoot string
	ConfDir               string
	Cluster               string
	Offline               bool
	DisableDataCollection bool
	BindIP                string
	Players               int
	BackupLogs            bool
	Tick                  int
}

func buildBaseArgs(opt options) []string {
	baseArgs := []string{
		"-persistent_storage_root", opt.PersistentStorageRoot,
		"-conf_dir", opt.ConfDir,
		"-cluster", opt.Cluster,
	}
	if opt.Offline {
		baseArgs = append(baseArgs, "-offline")
	}
	if opt.DisableDataCollection {
		baseArgs = append(baseArgs, "-disabledatacollection")
	}
	if opt.BindIP != "" {
		baseArgs = append(baseArgs, "-bind_ip", opt.BindIP)
	}
	if opt.Players > 0 {
		baseArgs = append(baseArgs, "-players", strconv.Itoa(opt.Players))
	}
	if opt.BackupLogs {
		baseArgs = append(baseArgs, "-backup_logs")
	}
	if opt.Tick > 0 {
		baseArgs = append(baseArgs, "-tick", strconv.Itoa(opt.Tick))
	}
	return baseArgs
}

func getDefaultRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	switch runtime.GOOS {
	case "windows":
		return home + "\\Documents\\Klei"
	case "darwin":
		return home + "/Documents/Klei"
	case "linux":
		return home + "/.klei"
	}
	return ""
}

func getOptions() options {
	var opt options
	flag.StringVar(&opt.PersistentStorageRoot, "persistent_storage_root", getDefaultRoot(),
		"Change the directory that your configuration directory resides in.")
	flag.StringVar(&opt.ConfDir, "conf_dir", "DoNotStarveTogether", "Change the name of your configuration directory.")
	flag.StringVar(&opt.Cluster, "cluster", "Cluster_1",
		"Set the name of the cluster directory that this server will use.")
	flag.BoolVar(&opt.Offline, "offline", false, "Start the server in offline mode.")
	flag.BoolVar(&opt.DisableDataCollection, "disabledatacollection", false, "Disable data collection for the server.")
	flag.StringVar(&opt.BindIP, "bind_ip", "",
		"Change the address that the server binds to when listening for player connections.")
	flag.IntVar(&opt.Players, "players", 0, "Set the maximum number of players that will be allowed to join the game.")
	flag.BoolVar(&opt.BackupLogs, "backup_logs", false,
		"Create a backup of the previous log files each time the server is run.")
	flag.IntVar(&opt.Tick, "tick", 0,
		"This is the number of times per-second that the server sends updates to clients.")
	flag.Parse()
	return opt
}

func main() {
	opt := getOptions()
	if opt.PersistentStorageRoot == "" {
		fmt.Printf("cannot resolve the system persistent storage root\n")
		os.Exit(1)
	}
	clusterDir := filepath.Join(opt.PersistentStorageRoot, opt.ConfDir, opt.Cluster)
	dir, err := os.Open(clusterDir)
	if err != nil {
		fmt.Printf("path \"%s\" is not exist\n", clusterDir)
		os.Exit(1)
	}
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		fmt.Printf("path \"%s\" is not a directory\n", clusterDir)
		os.Exit(1)
	}
	var shards []string
	var clusterToken, clusterConfig bool
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			switch fileInfo.Name() {
			case "cluster.ini":
				clusterConfig = true
			case "cluster_token.txt":
				clusterToken = true
			}
			continue
		}
		shard := fileInfo.Name()
		if info, err := os.Stat(filepath.Join(clusterDir, shard, "server.ini")); err != nil || info.IsDir() {
			continue
		}
		shards = append(shards, shard)
	}
	if !clusterConfig {
		fmt.Printf("configuration \"cluster.ini\" does not exist in cluster \"%s\"\n", opt.Cluster)
		os.Exit(1)
	}
	if !clusterToken {
		fmt.Printf("token \"cluster_token.txt\" does not exist in cluster \"%s\"\n", opt.Cluster)
		os.Exit(1)
	}
	if len(shards) == 0 {
		fmt.Printf("cluster \"%s\" does not contain any shard configuration\n", opt.Cluster)
		os.Exit(1)
	}
	baseArgs := buildBaseArgs(opt)
	_ = baseArgs // TODO: Do something with it.
	fmt.Printf("%#v\n", opt)
}
