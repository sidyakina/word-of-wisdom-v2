package challenger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_calculateLeadingZeros(t *testing.T) {
	tests := []struct {
		name      string
		hash      string
		wantZeros int32
	}{
		{
			name:      "only zeros",
			hash:      "00000",
			wantZeros: 20,
		},
		{
			name:      "1 = 3 zeros, stop after",
			hash:      "010",
			wantZeros: 7,
		},
		{
			name:      "2 = 2 zeros, stop after",
			hash:      "020",
			wantZeros: 6,
		},
		{
			name:      "3 = 2 zeros, stop after",
			hash:      "030",
			wantZeros: 6,
		},
		{
			name:      "4 = 1 zero, stop after",
			hash:      "040",
			wantZeros: 5,
		},
		{
			name:      "5 = 1 zero, stop after",
			hash:      "050",
			wantZeros: 5,
		},
		{
			name:      "6 = 1 zero, stop after",
			hash:      "060",
			wantZeros: 5,
		},
		{
			name:      "7 = 1 zero, stop after",
			hash:      "070",
			wantZeros: 5,
		},
		{
			name:      "8 hasn't zeros, stop after",
			hash:      "080",
			wantZeros: 4,
		},
		{
			name:      "9 hasn't zeros, stop after",
			hash:      "090",
			wantZeros: 4,
		},
		{
			name:      "A hasn't zeros, stop after",
			hash:      "0A0",
			wantZeros: 4,
		},
		{
			name:      "B hasn't zeros, stop after",
			hash:      "0B0",
			wantZeros: 4,
		},
		{
			name:      "C hasn't zeros, stop after",
			hash:      "0C0",
			wantZeros: 4,
		},
		{
			name:      "D hasn't zeros, stop after",
			hash:      "0D0",
			wantZeros: 4,
		},
		{
			name:      "E hasn't zeros, stop after",
			hash:      "0E0",
			wantZeros: 4,
		},
		{
			name:      "F hasn't zeros, stop after",
			hash:      "0F0",
			wantZeros: 4,
		},
		{
			name:      "0 leading zeros",
			hash:      "A",
			wantZeros: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotZeros := calculateLeadingZeros(tt.hash)

			assert.Equal(t, tt.wantZeros, gotZeros)
		})
	}
}
