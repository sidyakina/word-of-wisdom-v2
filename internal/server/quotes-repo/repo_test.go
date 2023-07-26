package quotesrepo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	repo, err := New()

	assert.Nilf(t, err, "err: %v", err)
	assert.Len(t, repo.quotes, 33)
}
