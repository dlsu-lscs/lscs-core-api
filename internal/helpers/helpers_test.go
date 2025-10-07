package helpers

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullStringToString(t *testing.T) {
	t.Run("valid string", func(t *testing.T) {
		s := sql.NullString{String: "hello", Valid: true}
		assert.Equal(t, "hello", NullStringToString(s))
	})

	t.Run("invalid string", func(t *testing.T) {
		s := sql.NullString{String: "", Valid: false}
		assert.Equal(t, "", NullStringToString(s))
	})
}
