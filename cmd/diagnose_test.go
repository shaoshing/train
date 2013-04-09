package main

import (
	"github.com/shaoshing/gotest"
	"testing"
)

func TestDiagnose(t *testing.T) {
	assert.Test = t
	assert.True(diagnose())
}
