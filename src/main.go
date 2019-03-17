package main

import (
	"nipplebot/cmd"
)

func main() {
	// CPU profiling by default
	//defer profile.Start(profile.CPUProfile).Stop()
	// Memory profiling
	//defer profile.Start(profile.MemProfile).Stop()
	nipplebot.Execute()
}
