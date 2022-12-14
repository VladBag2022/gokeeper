package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomBytes(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{
			name: "positive - zero-length",
			n:    0,
		},
		{
			name: "positive - 10 bytes",
			n:    10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomBytes(tt.n)
			require.NoError(t, err)
			assert.Equal(t, tt.n, len(got))
		})
	}
}
