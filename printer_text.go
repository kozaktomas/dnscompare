package main

import (
	"fmt"
	"github.com/fatih/color"
)

type TextResultPrinter struct {
	NoColors  bool
	Identical bool
}

func (printer TextResultPrinter) Print(results []DnsResult) {
	for _, result := range results {

		if !printer.Identical && result.Identical {
			continue
		}

		if printer.NoColors {
			printer.PrintNoColors(result)
		} else {
			printer.PrintColors(result)
		}
	}
}

func (printer TextResultPrinter) PrintColors(result DnsResult) {
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

func (printer TextResultPrinter) PrintNoColors(result DnsResult) {
	if !result.Identical {
		fmt.Println("DIFFERENT")
	}
	fmt.Println(fmt.Sprintf("%s -> %s", result.Type, result.Host))

	for _, response := range result.Responses {
		fmt.Print(response.Value)
		if response.Value == "" {
			fmt.Print("EMPTY")
		}
		fmt.Print(" <==> ")
		fmt.Println(response.Resolver)
	}
	fmt.Println("---")
}
