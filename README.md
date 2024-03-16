# DNS Compare

Program to compare DNS records from multiple nameservers. It's especially useful for DNS migrations. You can compare the
records from the old nameserver with the new one to make sure everything is correct. Right now there is support for A,
CNAME, MX and TXT records.

## Usage

```bash
$ dnscompare --help
DNS Compare is a tool to compare DNS records from different nameservers. It reads a list of DNS records and a list of DNS resolvers from a configuration file, resolves the records using the resolvers and prints the results.

Usage:
  dnscompare [config-file.yml] [flags]

Flags:
  -h, --help        help for dnscompare
      --identical   print identical results as well
      --json        print results in JSON format
```

## Configuration example

```yaml
# Example configuration file for DNS resolver
# Comparing Cloudflare and Google DNS servers
# For following records: A, CNAME, MX, TXT
nameservers:
  - "1.1.1.1:53"
  - "8.8.8.8:53"
records:
  - "A google.com"
  - "CNAME google.com"
  - "MX gmail.com"
  - "TXT google.com"
```

## Example output

```bash
$ dnscompare example_config.yml --identical
DIFFERENT                                                                                                                                                                   
A -> google.com
142.251.36.110 2a00:1450:4014:80b::200e <==> 1.1.1.1:53
142.251.37.110 2a00:1450:4014:80e::200e <==> 8.8.8.8:53
---
CNAME -> google.com
google.com. <==> 1.1.1.1:53
google.com. <==> 8.8.8.8:53
---
MX -> gmail.com
10 alt1.gmail-smtp-in.l.google.com. 20 alt2.gmail-smtp-in.l.google.com. 30 alt3.gmail-smtp-in.l.google.com. 40 alt4.gmail-smtp-in.l.google.com. 5 gmail-smtp-in.l.google.com. <==> 1.1.1.1:53
10 alt1.gmail-smtp-in.l.google.com. 20 alt2.gmail-smtp-in.l.google.com. 30 alt3.gmail-smtp-in.l.google.com. 40 alt4.gmail-smtp-in.l.google.com. 5 gmail-smtp-in.l.google.com. <==> 8.8.8.8:53
---
```
