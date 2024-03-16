package main

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

// DnsResolvers is a list of DNS resolvers including IP addresses and ports
// Example: []string{"8.8.8.8:53", "1.1.1.1:53",}
type DnsResolvers []string

// DnsResponse is a response from a DNS resolver
// The Value field contains string representation of the response
// It encapsulates the response for any DNS record type
// For multiple records, the values are joined with a space and sorted by alphabet
type DnsResponse struct {
	Resolver string `json:"resolver,omitempty"`
	Value    string `json:"value,omitempty"`
}

// DnsResult is a result of a DNS query for multiple nameservers
type DnsResult struct {
	Type      string        `json:"type,omitempty"`
	Host      string        `json:"host,omitempty"`
	Responses []DnsResponse `json:"responses,omitempty"`
	Identical bool          `json:"identical,omitempty"` // true when the responses from all resolvers are identical
}

// LookupA returns the list of A records for the given host
func (rs DnsResolvers) LookupA(host string) (DnsResult, error) {
	responses := make([]DnsResponse, len(rs))

	for resolverIndex, resolver := range rs {
		r := createResolver(resolver)
		ips, err := r.LookupIPAddr(context.Background(), host)
		if err != nil {
			responses[resolverIndex] = DnsResponse{
				Resolver: resolver,
				Value:    "",
			}
		} else {
			// convert ips to the list of strings
			values := make([]string, len(ips))
			for i, ip := range ips {
				values[i] = ip.String()
			}
			sort.Strings(values)
			responses[resolverIndex] = DnsResponse{
				Resolver: resolver,
				Value:    strings.Join(values, " "),
			}
		}
	}

	return DnsResult{
		Type:      "A",
		Host:      host,
		Responses: responses,
		Identical: areResponsesIdentical(responses),
	}, nil
}

// LookupCname returns the CNAME record for the given host
func (rs DnsResolvers) LookupCname(host string) (DnsResult, error) {
	responses := make([]DnsResponse, len(rs))

	for resolverIndex, resolver := range rs {
		r := createResolver(resolver)
		cname, err := r.LookupCNAME(context.Background(), host)
		if err != nil {
			responses[resolverIndex] = DnsResponse{
				Resolver: resolver,
				Value:    "",
			}
		} else {
			responses[resolverIndex] = DnsResponse{
				Resolver: resolver,
				Value:    cname,
			}
		}
	}

	return DnsResult{
		Type:      "CNAME",
		Host:      host,
		Responses: responses,
		Identical: areResponsesIdentical(responses),
	}, nil
}

// LookupMx returns the list of MX records for the given host
func (rs DnsResolvers) LookupMx(host string) (DnsResult, error) {
	responses := make([]DnsResponse, len(rs))

	for resolverIndex, resolver := range rs {
		r := createResolver(resolver)
		mxs, err := r.LookupMX(context.Background(), host)
		if err != nil {
			responses[resolverIndex] = DnsResponse{
				Resolver: resolver,
				Value:    "",
			}
		} else {
			// convert ips to the list of strings
			values := make([]string, len(mxs))
			for i, mx := range mxs {
				values[i] = fmt.Sprintf("%d %s", mx.Pref, mx.Host)
			}
			sort.Strings(values)
			responses[resolverIndex] = DnsResponse{
				Resolver: resolver,
				Value:    strings.Join(values, " "),
			}
		}
	}

	return DnsResult{
		Type:      "MX",
		Host:      host,
		Responses: responses,
		Identical: areResponsesIdentical(responses),
	}, nil
}

// LookupTxt returns the TXT records for the given host
func (rs DnsResolvers) LookupTxt(host string) (DnsResult, error) {
	responses := make([]DnsResponse, len(rs))

	for resolverIndex, resolver := range rs {
		r := createResolver(resolver)
		txts, err := r.LookupTXT(context.Background(), host)
		if err != nil {
			responses[resolverIndex] = DnsResponse{
				Resolver: resolver,
				Value:    "",
			}
		} else {
			sort.Strings(txts)
			responses[resolverIndex] = DnsResponse{
				Resolver: resolver,
				Value:    strings.Join(txts, " "),
			}
		}
	}

	return DnsResult{
		Type:      "TXT",
		Host:      host,
		Responses: responses,
		Identical: areResponsesIdentical(responses),
	}, nil
}

// createResolver creates a resolver with a custom dialer
func createResolver(server string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, network, server)
		},
	}
}

// areResponsesIdentical returns true when the response from all resolvers is identical
func areResponsesIdentical(responses []DnsResponse) bool {
	if len(responses) == 0 {
		return true
	}

	f := responses[0].Value
	for _, response := range responses {
		if f != response.Value {
			return false
		}
	}

	return true
}
