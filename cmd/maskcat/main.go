// Package main controls the primary logic for the application
//
// The package leans on /internal/cli to perform command line actions
// The application logic is stored within /pkg/*
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/jakewnuk/maskcat/internal/cli"
)

var version = "2.0.0"

func main() {
	flagSet := flag.NewFlagSet("maskcat", flag.ExitOnError)
	doVerbose := flagSet.Bool("v", false, "Show verbose information about masks\nExample: maskcat [MODE] -v")
	doMultiByte := flagSet.Bool("m", false, "Process multibyte text (warning: slows processes)\nExample: maskcat [MODE] -m")
	doNumberOfReplacements := flagSet.Int("n", 1, "Max number of replacements to make per item (default: 1)\nExample: maskcat [MODE] -n 1")
	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "Options for maskcat (version %s):\n", version)
		flagSet.PrintDefaults()
		printUsage()
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(flagSet.Output(), "Options for maskcat (version %s):\n", version)
		flagSet.PrintDefaults()
		printUsage()
		os.Exit(0)
	}

	stdIn := bufio.NewScanner(os.Stdin)

	switch os.Args[1] {
	case "mask":
		flagSet.Parse(os.Args[2:])
		cli.GenerateMasks(stdIn, *doMultiByte, *doVerbose)
	case "match":
		flagSet.Parse(os.Args[2:])
		cli.CheckIfArgExists(2, os.Args)
		cli.MatchMasks(stdIn, os.Args[2], *doMultiByte)
	case "sub":
		flagSet.Parse(os.Args[3:])
		cli.CheckIfArgExists(2, os.Args)
		cli.SubMasks(stdIn, os.Args[2], *doMultiByte, *doNumberOfReplacements)
	case "mutate":
		flagSet.Parse(os.Args[3:])
		cli.CheckIfArgExists(2, os.Args)
		cli.MutateMasks(stdIn, os.Args[2], *doMultiByte, *doNumberOfReplacements)
	case "tokens":
		flagSet.Parse(os.Args[1:])
		cli.CheckIfArgExists(2, os.Args)
		cli.GenerateTokens(stdIn, os.Args[2])
	case "partial":
		flagSet.Parse(os.Args[1:])
		cli.CheckIfArgExists(2, os.Args)
		cli.GeneratePartialMasks(stdIn, os.Args[2])
	case "remove":
		flagSet.Parse(os.Args[1:])
		cli.CheckIfArgExists(2, os.Args)
		cli.GeneratePartialRemoveMasks(stdIn, os.Args[2])
	}
}

// printUsage prints usage information for the application
func printUsage() {
	fmt.Println(fmt.Sprintf("\nModes for maskcat (version %s):", version))
	fmt.Println("\n  mask\t\tCreates masks from text")
	fmt.Println("\t\tExample: stdin | maskcat mask [OPTIONS]")
	fmt.Println("\n  match\t\tMatches text to masks")
	fmt.Println("\t\tExample: stdin | maskcat match [MASK-FILE] [OPTIONS]")
	fmt.Println("\n  sub\t\tReplaces text with text from a file with masks")
	fmt.Println("\t\tExample: stdin | maskcat sub [TOKENS-FILE] [OPTIONS]")
	fmt.Println("\n  mutate\tMutates text by using chunking and token swapping")
	fmt.Println("\t\tExample: stdin | maskcat mutate [CHUNK-SIZE] [OPTIONS]")
	fmt.Println("\n  tokens\tSplits text into chunks by length (values over 99 allow all)")
	fmt.Println("\t\tExample: stdin | maskcat tokens [TOKEN-LEN] [OPTIONS]")
	fmt.Println("\n  partial\tPartially replaces characters with mask characters")
	fmt.Println("\t\tExample: stdin | maskcat partial [MASK-CHARS] [OPTIONS]")
	fmt.Println("\n  remove\tRemoves characters that match given mask characters")
	fmt.Println("\t\tExample: stdin | maskcat remove [MASK-CHARS] [OPTIONS]")
}