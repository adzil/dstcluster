// Copyright (c) 2019 Fadhli Dzil Ikram. All rights reserved.
// Use of source code is governed by a MIT license, see LICENSE for more info.

// Command dstcluster provides command-line interface for running and managing
// multiple shards in a Don't Starve Together cluster configuration.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	serverBinary = "./dontstarve_dedicated_server_nullrenderer"
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
		// TODO: Check if the following flag also works on Windows or not
		"-monitor_parent_process", strconv.Itoa(os.Getpid()),
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

func tryChdir() bool {
	if info, err := os.Stat(serverBinary); err == nil && !info.IsDir() {
		return true
	}
	execPath, _ := os.Executable()
	if execPath == "" {
		return false
	}
	if err := os.Chdir(filepath.Dir(execPath)); err != nil {
		return false
	}
	info, err := os.Stat(serverBinary)
	return err == nil && !info.IsDir()
}

func main() {
	if !tryChdir() {
		fmt.Printf("command must be executed and/or stored under the game's \"bin/\" directory\n")
		os.Exit(1)
	}
	opt := getOptions()
	if opt.PersistentStorageRoot == "" {
		fmt.Printf("cannot resolve the current system persistent storage root\n")
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
	var builder strings.Builder
	for i, shard := range shards {
		if i > 0 {
			builder.WriteString(", ")
			if i == len(shards)-1 {
				builder.WriteString(" and ")
			}
		}
		builder.WriteByte('"')
		builder.WriteString(shard)
		builder.WriteByte('"')
	}
	fmt.Printf("starting cluster \"%s\" with %d shard(s): %s\n", opt.Cluster, len(shards), builder.String())
	baseArgs := buildBaseArgs(opt)
	_ = baseArgs // TODO: Do something with it.
	fmt.Printf("%#v\n", opt)
}
