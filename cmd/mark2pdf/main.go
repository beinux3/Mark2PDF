package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/beinux3/Mark2PDF"
)

const version = "0.1.0"

func main() {
	// Define command line flags
	inputFile := flag.String("input", "", "Input Markdown file (required)")
	outputFile := flag.String("output", "", "Output PDF file (required)")
	showVersion := flag.Bool("version", false, "Show version information")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("mark2pdf version %s\n", version)
		fmt.Println("A pure Go Markdown to PDF converter")
		return
	}

	// Show help
	if *help || *inputFile == "" || *outputFile == "" {
		printHelp()
		if *inputFile == "" || *outputFile == "" {
			os.Exit(1)
		}
		return
	}

	// Check if input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Input file '%s' does not exist\n", *inputFile)
		os.Exit(1)
	}

	// Convert the file
	fmt.Printf("Converting '%s' to '%s'...\n", *inputFile, *outputFile)

	err := mark2pdf.ConvertFile(*inputFile, *outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during conversion: %v\n", err)
		os.Exit(1)
	}

	// Get output file size
	fileInfo, err := os.Stat(*outputFile)
	if err == nil {
		fmt.Printf("Success! PDF created: %s (%.2f KB)\n", *outputFile, float64(fileInfo.Size())/1024.0)
	} else {
		fmt.Printf("Success! PDF created: %s\n", *outputFile)
	}
}

func printHelp() {
	fmt.Println("Mark2PDF - Convert Markdown to PDF")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  mark2pdf -input <input.md> -output <output.pdf>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -input string")
	fmt.Println("        Input Markdown file (required)")
	fmt.Println("  -output string")
	fmt.Println("        Output PDF file (required)")
	fmt.Println("  -version")
	fmt.Println("        Show version information")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  mark2pdf -input README.md -output README.pdf")
	fmt.Println("  mark2pdf -input document.md -output document.pdf")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/beinux3/Mark2PDF")
}
