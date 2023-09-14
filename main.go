package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/flags"
	"github.com/sunshineplan/utils/httpsvr"
)

var (
	self string

	server = httpsvr.New()
	svc    = service.New()
	meta   metadata.Server

	joinPath = filepath.Join
	dir      = filepath.Dir
)

func init() {
	var err error
	self, err = os.Executable()
	if err != nil {
		svc.Fatalln("Failed to get self path:", err)
	}
	svc.Name = "Requirements"
	svc.Desc = "Instance to serve Requirements"
	svc.Exec = run
	svc.TestExec = test
	svc.Options = service.Options{
		Dependencies:       []string{"After=network.target"},
		Environment:        map[string]string{"GIN_MODE": "release"},
		RemoveBeforeUpdate: []string{"dist/assets"},
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(flag.CommandLine.Output(), svc.Usage())
	}
}

var logPath = flag.String("log", "", "Log file path")

func main() {
	flag.StringVar(&meta.Addr, "server", "", "Metadata Server Address")
	flag.StringVar(&meta.Header, "header", "", "Verify Header Header Name")
	flag.StringVar(&meta.Value, "value", "", "Verify Header Value")
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	flags.SetConfigFile(joinPath(dir(self), "config.ini"))
	flags.Parse()

	if err := svc.ParseAndRun(flag.Args()); err != nil {
		svc.Fatal(err)
	}
}
