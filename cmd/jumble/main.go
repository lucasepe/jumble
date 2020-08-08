package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/jumble"
	"github.com/lucasepe/jumble/config"
)

const (
	banner = `
   _           _   _
  |_|_ _ _____| |_| |___    Crafted with passion by Luca Sepe 
  | | | |     | . | | -_|   
 _| |___|_|_|_|___|_|___|   https://github.com/lucasepe/jumble         
|___| v{{VERSION}}`
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	flagTileSize int
	flagOutput   string
)

func main() {
	configureFlags()

	if flag.CommandLine.Arg(0) == "" {
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	uri := flag.Args()[0]

	cfg, err := config.DecodeURI(uri)
	handleErr(err)

	if flagTileSize <= 16 {
		flagTileSize = 16
	}

	if flagTileSize > 96 {
		flagTileSize = 96
	}

	grid, err := jumble.NewGrid(cfg.Rows, cfg.Cols, flagTileSize,
		jumble.GridBackground(cfg.Background),
		jumble.GridMargin(cfg.Margin),
	)
	handleErr(err)

	if cfg.Grid {
		grid.DrawGrid()
	}
	if cfg.Border {
		grid.DrawBorder()
	}
	if cfg.Hints {
		grid.DrawCoords()
	}

	for _, tile := range cfg.Tiles {
		handleErr(tile.Plot(grid))
	}

	if len(flagOutput) <= 1 {
		handleErr(grid.EncodePNG(os.Stdout))
	} else {
		handleErr(grid.SavePNG(flagOutput))
	}
}

func configureFlags() {
	name := appName()

	flag.CommandLine.Usage = func() {
		printBanner()
		fmt.Printf("Create diagrams stitching and connecting images.\n\n")

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [options] <hcl file or url>\n\n", name)

		fmt.Print("EXAMPLE:\n\n")
		fmt.Printf("  %s -s 64 -o test.png test.hcl\n", name)
		fmt.Println()

		fmt.Print("OPTIONS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
		fmt.Print("  -help\n\tprints this message\n")
		fmt.Println()
	}

	flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)

	flag.CommandLine.IntVar(&flagTileSize, "s", 72, "cell size in pixel; min:16 max:96")
	flag.CommandLine.StringVar(&flagOutput, "o", "", "write to file instead of stdout")

	flag.CommandLine.Parse(os.Args[1:])
}

// handleErr check for an error and eventually exit
func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func printBanner() {
	fmt.Print(strings.Trim(strings.Replace(banner, "{{VERSION}}", version, 1), "\n"), "\n\n")
}

func appName() string {
	return filepath.Base(os.Args[0])
}
