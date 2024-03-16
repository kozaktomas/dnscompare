package main

import (
	"fmt"
	"github.com/fatih/color"
)

type TextResultPrinter struct {
	Identical bool
}

func (printer TextResultPrinter) Print(results []DnsResult) {
	diffCount := 0
	for _, result := range results {

		if !printer.Identical && result.Identical {
			continue
		}

		printer.printResult(result)
		diffCount++
	}

	if !printer.Identical && diffCount == 0 {
		greenColor := color.New(color.FgGreen)
		_, _ = greenColor.Println("All records are identical")
	}
}

func (printer TextResultPrinter) printResult(result DnsResult) {
	headingColor := color.New(color.FgGreen)
	if !result.Identical {
		headingColor = color.New(color.FgRed)
		_, _ = headingColor.Println("DIFFERENT")
	}
	_, _ = headingColor.Println(fmt.Sprintf("%s -> %s", result.Type, result.Host))

	recordColor := color.New(color.FgYellow)
	nameserverColor := color.New(color.FgCyan)
	for _, response := range result.Responses {
		_, _ = recordColor.Print(response.Value)
		if response.Value == "" {
			_, _ = recordColor.Print("EMPTY")
		}
		fmt.Print(" <==> ")
		_, _ = nameserverColor.Println(response.Resolver)
	}
	fmt.Println("---")
}
