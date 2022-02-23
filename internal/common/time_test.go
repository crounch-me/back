package common_test

import (
	"testing"
	"time"

	"github.com/crounch-me/back/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestIsInRangeNominal(t *testing.T) {
	now := time.Now()
	in_range := common.IsInRange(now, now.Add(-time.Hour), now.Add(time.Hour))

	assert.True(t, in_range)
}
