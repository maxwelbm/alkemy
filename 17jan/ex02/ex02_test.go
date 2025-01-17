package ex02

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactorial(t *testing.T) {
	tests := []struct {
		name string
		args int
		want int
	}{
		{
			name: "fatorial de 0 deve ser 1",
			args: 0,
			want: 1,
		},
		{
			name: "fatorial de 1 deve ser 1",
			args: 1,
			want: 1,
		},
		{
			name: "fatorial de 2 deve ser 2",
			args: 2,
			want: 2,
		},
		{
			name: "fatorial de 3 deve ser 6",
			args: 3,
			want: 6,
		},
		{
			name: "fatorial de 4 deve ser 24",
			args: 4,
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Factorial(tt.args)
			assert.Equal(t, tt.want, got, "Factorial(%d) = %v, want %v", tt.args, got, tt.want)
		})
	}
}
