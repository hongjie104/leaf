package leaf

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hongjie104/leaf/cluster"
	"github.com/hongjie104/leaf/console"
	"github.com/hongjie104/leaf/log"
	"github.com/hongjie104/leaf/module"
)

// Run Run
func Run(mods ...module.Module) {
	log.Logger = log.New()
	defer log.Logger.Sync()

	log.Logger.Infof("Leaf %s starting up", version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-c
	log.Infof("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
	log.Info("Leaf closed success")
}
