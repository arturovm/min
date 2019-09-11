package min_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arturovm/min"
)

func TestNewSubGroup(t *testing.T) {
	m := min.New(nil)
	g := m.NewGroup("/test")

	require.NotNil(t, g)
	require.Equal(t, m.Group, g.Parent())
}

func TestRootPath(t *testing.T) {
	m := min.New(nil)
	require.Equal(t, "/", m.Path)
}

func TestSubGroupPath(t *testing.T) {
	g := min.New(nil).NewGroup("/sub")

	require.NotNil(t, g)
	require.Equal(t, "/sub", g.Path)
}

func TestFullPath(t *testing.T) {
	g := min.New(nil).NewGroup("/sub").NewGroup("/group")

	require.NotNil(t, g)
	require.Equal(t, "/sub/group", g.FullPath())
}
