package main

type ResultPrinter interface {
	Print(results []DnsResult)
}
