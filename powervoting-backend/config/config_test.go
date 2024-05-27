package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	InitConfig("../")

	assert.NotEmpty(t, Client)
}
