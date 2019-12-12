package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUser(t *testing.T) {
	tests := map[string]struct {
		id       string
		expected bool
	}{
		"valid ID":   {id: "WR8LAPMKK", expected: true},
		"invalid ID": {id: "THEBADMAN", expected: false},
	}

	for name, test := range tests {
		test := test

		t.Run(name, func(t *testing.T) {
			var handler slashCommandHandler
			result := handler.validateUser(test.id)
			assert.Equal(t, test.expected, result)
		})
	}
}
