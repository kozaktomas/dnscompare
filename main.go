package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type Config struct {
	Nameservers []string `yaml:"nameservers"`
	Records     []string `yaml:"records"`
}

type Record struct {
	Type string
	Addr string
}

const (
	ParamJson      = "json"
	ParamIdentical = "identical"
)

func main() {

	rootCmd := &cobra.Command{
		Use:   "dnscompare [config-file.yml]",
		Short: "DNS Compare is a tool to compare DNS records from different nameservers.",
		Long: "DNS Compare is a tool to compare DNS records from different nameservers. " +
			"It reads a list of DNS records and a list of DNS resolvers from a configuration file, " +
			"resolves the records using the resolvers and prints the results.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires a config file as an argument")
			}
			if _, err := os.Stat(args[0]); err != nil {
				return fmt.Errorf("config file %q does not exist", args[0])
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			configFilename := args[0]

			cfgData, err := os.ReadFile(configFilename)
			if err != nil {
				log.Fatalf("could not read config file %q: %s", configFilename, err)
			}

			var config Config
			err = yaml.Unmarshal(cfgData, &config)
			if err != nil {
				log.Fatalf("could not unmarshal config: %s", err)
			}

			resolvers := DnsResolvers(config.Nameservers)
			results := make([]DnsResult, len(config.Records))
			bar := progressbar.Default(int64(len(config.Records)), "Resolving records")

			for i, recordText := range config.Records {
				var result DnsResult
				record, err := parseRecord(recordText)
				if err != nil {
					log.Fatalf("could not parse record: %s", err)
				}

				if record.Type == "A" {
					result, err = resolvers.LookupA(record.Addr)
				}

				if record.Type == "CNAME" {
					result, err = resolvers.LookupCname(record.Addr)
				}

				if record.Type == "MX" {
					result, err = resolvers.LookupMx(record.Addr)
				}

				if record.Type == "TXT" {
					result, err = resolvers.LookupTxt(record.Addr)
				}

				if err != nil {
					log.Fatalf("could not resolve record: %s", err)
				}

				results[i] = result
				_ = bar.Add(1)
			}

			_ = bar.Close()

			ignoreErrBool := func(x bool, err error) bool { return x }
			if ignoreErrBool(cmd.Flags().GetBool(ParamJson)) {
				printer := JsonResultPrinter{}
				printer.Print(results)
			} else {
				printer := TextResultPrinter{
					Identical: ignoreErrBool(cmd.Flags().GetBool(ParamIdentical)),
				}
				printer.Print(results)
			}
		},
	}

	rootCmd.PersistentFlags().Bool(ParamJson, false, "print results in JSON format")
	rootCmd.PersistentFlags().Bool(ParamIdentical, false, "print identical results as well")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("could not run root command")
	}
}

func parseRecord(line string) (Record, error) {
	chunks := strings.Split(strings.TrimSpace(line), " ")
	if len(chunks) != 2 {
		return Record{}, fmt.Errorf("could not parse line: %s", line)
	}

	rt := strings.ToUpper(chunks[0]) // record type
	if !IsValidRecordType(rt) {
		return Record{}, fmt.Errorf("invalid record type: %s", rt)

	}

	return Record{
		Type: rt,
		Addr: chunks[1],
	}, nil
}

func IsValidRecordType(recordType string) bool {
	switch recordType {
	case
		"A",
		"CNAME",
		"TXT",
		"MX":
		return true
	}
	return false
}
