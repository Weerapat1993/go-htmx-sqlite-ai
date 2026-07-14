//go:build e2e

package e2e_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbout(t *testing.T) {
	t.Parallel()
	_, page := newPage(t)
	_, err := page.Goto(getFullPath("/about"))
	require.NoError(t, err)

	require.NoError(t, expect.Locator(page.GetByText("Why this stack")).ToBeVisible())
}

func TestAbout_NavLinkFromHome(t *testing.T) {
	t.Parallel()
	_, page := newPage(t)
	_, err := page.Goto(getFullPath(""))
	require.NoError(t, err)

	err = page.Locator(`a[href="/about"]`).First().Click()
	require.NoError(t, err)

	require.NoError(t, expect.Locator(page.GetByText("Why this stack")).ToBeVisible())
}
