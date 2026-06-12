package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		length      int
		expectError bool
	}{
		{
			name:        "success - length 7",
			length:      7,
			expectError: false,
		},
		{
			name:        "success - length 10",
			length:      10,
			expectError: false,
		},
		{
			name:        "zero length",
			length:      0,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := GenerateShortCode(tc.length)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tc.length)
			}
		})
	}
}
