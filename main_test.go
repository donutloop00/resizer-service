package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidDimensionFunction(t *testing.T) {
	assert.Equal(t, checkValidDimensions(30, 50), true, "Function should return false for an invalid dimension")
	assert.Equal(t, checkValidDimensions(0, 50), false, "Function should return false for an invalid dimension")
	assert.Equal(t, checkValidDimensions(-30, 50), false, "Function should return false for a negative dimension")
	assert.Equal(t, checkValidDimensions(30, -50), false, "Function should return false for a negative dimension")
}