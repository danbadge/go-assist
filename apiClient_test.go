package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestExamples(t *testing.T) {
	t.Run("this succeeds", func(t *testing.T){
		assert.Equal(t, 4, 4)
	})

	t.Run("this fails", func(t *testing.T){
		assert.Equal(t, 1, 2)
	})
}