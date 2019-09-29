// Copyright (c) 2019 Fadhli Dzil Ikram. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

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

type options struct {
	ServerPath            string
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

func defaultRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	switch runtime.GOOS {
	case "windows":
		return home + "\\Documents\\Klei"
	case "darwin":
		return home + "/Documents/Klei"
	}
	return home + "/.klei"
}

func defaultServerPath() string {
	switch runtime.GOOS {
	case "windows":
		return ".\\dontstarve_dedicated_server_nullrenderer.exe"
	}
	return "./dontstarve_dedicated_server_nullrenderer"
}

func getPlural(n int) string {
	if n > 1 {
		return "s"
	}
	return ""
}

func getConcat(items []string) string {
	var builder strings.Builder
	for i, shard := range items {
		if i > 0 {
			builder.WriteString(", ")
			if i == len(items)-1 {
				builder.WriteString(" and ")
			}
		}
		builder.WriteByte('"')
		builder.WriteString(shard)
		builder.WriteByte('"')
	}
	return builder.String()
}

func getOptions() options {
	var opt options
	flag.StringVar(&opt.ServerPath, "server_path", defaultServerPath(), "Change the dedicated game server binary path.")
	flag.StringVar(&opt.PersistentStorageRoot, "persistent_storage_root", defaultRoot(),
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

func resolveServerPath(serverPath string) string {
	if info, err := os.Stat(serverPath); err == nil && !info.IsDir() {
		return serverPath
	}
	if filepath.IsAbs(serverPath) {
		return ""
	}
	execPath, _ := os.Executable()
	if execPath == "" {
		return ""
	}
	execPath = filepath.Join(filepath.Dir(execPath), serverPath)
	if info, err := os.Stat(execPath); err == nil && !info.IsDir() {
		return execPath
	}
	return ""
}

func buildBaseArgs(opt options) []string {
	baseArgs := []string{
		"-persistent_storage_root", opt.PersistentStorageRoot,
		"-conf_dir", opt.ConfDir,
		"-cluster", opt.Cluster,
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

func fatalf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	os.Exit(1)
}

func main() {
	opt := getOptions()
	serverPath := resolveServerPath(opt.ServerPath)
	if serverPath == "" {
		fatalf("cannot find the game binary in \"%s\"\n", opt.ServerPath)
	}
	serverDir := filepath.Dir(serverPath)
	if err := os.Chdir(serverDir); err != nil {
		fatalf("cannot change working directory to \"%s\": %s\n", serverDir, err.Error())
	}
	if opt.PersistentStorageRoot == "" {
		fatalf("cannot resolve the current system persistent storage root\n")
	}
	clusterDir := filepath.Join(opt.PersistentStorageRoot, opt.ConfDir, opt.Cluster)
	dir, err := os.Open(clusterDir)
	if err != nil {
		if os.IsNotExist(err) {
			fatalf("path \"%s\" is not exist\n", clusterDir)
		}
		fatalf("cannot open \"%s\": %s\n", clusterDir, err.Error())
	}
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		fatalf("path \"%s\" is not a directory\n", clusterDir)
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
		fatalf("configuration \"cluster.ini\" does not exist in cluster \"%s\"\n", opt.Cluster)
	}
	if !clusterToken {
		fatalf("token \"cluster_token.txt\" does not exist in cluster \"%s\"\n", opt.Cluster)
	}
	if len(shards) == 0 {
		fatalf("cluster \"%s\" does not contain any shard configuration\n", opt.Cluster)
	}

	fmt.Printf("starting cluster \"%s\" with %d shard%s: %s\n", opt.Cluster, len(shards), getPlural(len(shards)),
		getConcat(shards))
	baseArgs := buildBaseArgs(opt)
	_ = baseArgs // TODO: Do something with it.
}
