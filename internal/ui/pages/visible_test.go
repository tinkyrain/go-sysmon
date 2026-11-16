package pages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisible(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		elements []int
		count    int
		expected []int
	}{
		{
			name:     "everything fits and comes back whole",
			elements: []int{1, 2, 3},
			count:    3,
			expected: []int{1, 2, 3},
		},
		{
			name:     "more items than room keeps only the first ones",
			elements: []int{1, 2, 3},
			count:    1,
			expected: []int{1},
		},
		{
			name:     "no room returns nothing",
			elements: []int{1, 2, 3},
			count:    0,
			expected: nil,
		},
		{
			name:     "a negative count returns nothing",
			elements: []int{1, 2, 3},
			count:    -1,
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := visible(tc.elements, tc.count)
			assert.Equal(t, tc.expected, result)
		})
	}
}
