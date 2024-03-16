package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidRecordType(t *testing.T) {
	assert.True(t, IsValidRecordType("A"))
	assert.True(t, IsValidRecordType("CNAME"))
	assert.True(t, IsValidRecordType("MX"))
	assert.True(t, IsValidRecordType("TXT"))

	assert.False(t, IsValidRecordType("B"))
	assert.False(t, IsValidRecordType("CNAMEE"))
}

func TestParseRecord(t *testing.T) {
	r, err := parseRecord("A test.com")
	assert.Nil(t, err)
	assert.Equal(t, "A", r.Type)
	assert.Equal(t, "test.com", r.Addr)

	r, err = parseRecord("test.com")
	assert.NotNil(t, err)
}
