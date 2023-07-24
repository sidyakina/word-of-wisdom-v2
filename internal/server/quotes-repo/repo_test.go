package quotesrepo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	repo, err := New()

	assert.Nilf(t, err, "err must be nil, but: %v", err)
	assert.Len(t, repo.quotes, 33)
}
