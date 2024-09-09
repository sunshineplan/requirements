package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/flags"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/utils/mail"
)

var (
	self string

	server = httpsvr.New()
	svc    = service.New()
	meta   metadata.Server
	dialer mail.Dialer

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
		Dependencies: []string{"After=network.target"},
		Environment:  map[string]string{"GIN_MODE": "release"},
	}
	svc.RegisterCommand("backup", "backup", func(_ ...string) error {
		if err := initSrv(); err != nil {
			return err
		}
		backup()
		return nil
	}, 0, true)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(flag.CommandLine.Output(), svc.Usage())
	}
}

var (
	name      = flag.String("name", "业务系统", "Done value")
	doneValue = flag.String("done", "已完成", "Done value")
	poll      = flag.Duration("poll", 50*time.Second, "Poll interval")
	logPath   = flag.String("log", "", "Log file path")
)

func main() {
	flag.TextVar(&to, "to", mail.Receipts(nil), "Mail Subscriber")
	flag.StringVar(&meta.Addr, "server", "", "Metadata Server Address")
	flag.StringVar(&meta.Header, "header", "", "Verify Header Header Name")
	flag.StringVar(&meta.Value, "value", "", "Verify Header Value")
	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	flag.StringVar(&svc.Options.PIDFile, "pid", "/var/run/requirements.pid", "PID file path")
	flags.SetConfigFile(joinPath(dir(self), "config.ini"))
	flags.Parse()

	if *logPath != "" {
		svc.SetLogger(*logPath, "", log.LstdFlags)
		gin.DefaultWriter = svc.Logger
		gin.DefaultErrorWriter = svc.Logger
	}

	if err := svc.ParseAndRun(flag.Args()); err != nil {
		svc.Fatal(err)
	}
}
