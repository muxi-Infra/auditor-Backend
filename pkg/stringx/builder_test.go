package stringx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuild(t *testing.T) {
	s1 := "test1"
	s2 := "test2"
	re := Build(s1, s2)
	assert.Equal(t, "test1test2", re)
}
