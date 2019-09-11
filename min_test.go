package min_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arturovm/min"
)

func TestNew(t *testing.T) {
	m := min.New(nil)
	require.NotNil(t, m)
}
