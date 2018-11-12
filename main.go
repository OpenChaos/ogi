package main

import (
	"fmt"
	"runtime"

	ogiconsumer "github.com/OpenChaos/ogi/consumer"
	logger "github.com/OpenChaos/ogi/logger"
)

func main() {
	logger.SetupLogger()
	fmt.Printf(`

		   oooo          ggg      iiii      \  cores available: %d
		 ooo  ooo      gg  gg       ii      |
		000    000       gggg       ii      |
		 ooo  ooo           gg      ii      |
		   oooo         gggg        ii      /

`,
		runtime.NumCPU(),
	)
	ogiconsumer.Consume()
}
