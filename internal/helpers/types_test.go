package helpers

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullableString_MarshalJSON(t *testing.T) {
	t.Run("valid string", func(t *testing.T) {
		ns := NullableString{sql.NullString{String: "hello", Valid: true}}
		b, err := json.Marshal(ns)
		assert.NoError(t, err)
		assert.Equal(t, `"hello"`, string(b))
	})

	t.Run("invalid string", func(t *testing.T) {
		ns := NullableString{sql.NullString{String: "", Valid: false}}
		b, err := json.Marshal(ns)
		assert.NoError(t, err)
		assert.Equal(t, `""`, string(b))
	})
}
