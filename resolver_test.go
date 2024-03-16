package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAreResponsesIdentical(t *testing.T) {
	assert.True(t, areResponsesIdentical([]DnsResponse{
		{Resolver: "a", Value: "1.1.1.1"},
		{Resolver: "b", Value: "1.1.1.1"},
		{Resolver: "c", Value: "1.1.1.1"},
	}))

	assert.False(t, areResponsesIdentical([]DnsResponse{
		{Resolver: "a", Value: "1.1.1.1"},
		{Resolver: "b", Value: "1.1.1.2"},
		{Resolver: "c", Value: "1.1.1.2"},
	}))

	assert.False(t, areResponsesIdentical([]DnsResponse{
		{Resolver: "a", Value: "1.1.1.2"},
		{Resolver: "b", Value: "1.1.1.1"},
		{Resolver: "c", Value: "1.1.1.1"},
	}))
}
