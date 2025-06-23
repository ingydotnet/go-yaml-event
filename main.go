// Package main provides a YAML node inspection tool that reads YAML from stdin
// and outputs a detailed analysis of its node structure, including comments
// and content organization.
package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// main reads YAML from stdin, parses it, and outputs the node structure
func main() {
	parser, err := yaml.NewParser(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating parser: %v\n", err)
		os.Exit(1)
	}
	defer parser.Close()

	for {
		event, err := parser.Next()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parser error: %v\n", err)
			os.Exit(1)
		}
		if event == nil {
			break
		}

		// Print event information in YAML format
		fmt.Printf("- Event: %v\n", event.Type)
		fmt.Printf("  Start: {Line: %d, Column: %d}\n", 
			event.StartMark.Line+1, event.StartMark.Column)
		fmt.Printf("  End: {Line: %d, Column: %d}\n", 
			event.EndMark.Line+1, event.EndMark.Column)

		// Print any comments associated with the event
		if event.HeadComment != nil && len(event.HeadComment) > 0 {
			fmt.Printf("  HeadComment: %q\n", string(event.HeadComment))
		}
		if event.LineComment != nil && len(event.LineComment) > 0 {
			fmt.Printf("  LineComment: %q\n", string(event.LineComment))
		}
		if event.FootComment != nil && len(event.FootComment) > 0 {
			fmt.Printf("  FootComment: %q\n", string(event.FootComment))
		}
		if event.TailComment != nil && len(event.TailComment) > 0 {
			fmt.Printf("  TailComment: %q\n", string(event.TailComment))
		}

		switch event.Type {
		case yaml.EventScalar:
			fmt.Printf("  Value: %q\n", event.Value)
			if style := event.StyleString(); style != "" && style != "plain" {
				fmt.Printf("  Style: %s\n", style)
			}
			if event.Tag != "" {
				fmt.Printf("  Tag: %s\n", event.Tag)
			}
			if event.Anchor != "" {
				fmt.Printf("  Anchor: %s\n", event.Anchor)
			}
			fmt.Printf("  Implicit: %v\n", event.Implicit)
		case yaml.EventAlias:
			fmt.Printf("  Anchor: %s\n", event.Anchor)
		case yaml.EventSequenceStart, yaml.EventMappingStart:
			if style := event.StyleString(); style != "" && style != "block" {
				fmt.Printf("  Style: %s\n", style)
			}
			if event.Tag != "" {
				fmt.Printf("  Tag: %s\n", event.Tag)
			}
			if event.Anchor != "" {
				fmt.Printf("  Anchor: %s\n", event.Anchor)
			}
			fmt.Printf("  Implicit: %v\n", event.Implicit)
		case yaml.EventDocumentStart, yaml.EventDocumentEnd:
			fmt.Printf("  Implicit: %v\n", event.Implicit)
		}
		fmt.Println()
	}
}
