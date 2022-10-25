package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("../")

	require.NoError(t, err)
	require.NotEmpty(t, config)
}
