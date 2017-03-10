package looli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	assert.True(t, router.TrailingSlashRedirect)
	assert.NotNil(t, router.allowMethods)
	assert.NotNil(t, router.tree)
	assert.NotNil(t, router.tree.children)
	assert.NotNil(t, router.tree.handlers)
}
