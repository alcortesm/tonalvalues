package staircase_test

import (
	"fmt"
	"testing"

	"github.com/alcortesm/tonalvalues/staircase"
)

func TestTransform(t *testing.T) {
	t.Parallel()

	subtests := []struct {
		min, max, n uint
		input       int
		want        int
	}{
		{min: 0, max: 1, n: 2, input: 0, want: 0},
		{min: 0, max: 1, n: 2, input: 1, want: 1},

		{min: 0, max: 1, n: 3, input: 0, want: 0},
		{min: 0, max: 1, n: 3, input: 1, want: 1},

		{min: 0, max: 1, n: 100, input: 0, want: 0},
		{min: 0, max: 1, n: 100, input: 1, want: 1},

		{min: 0, max: 2, n: 2, input: 0, want: 0},
		{min: 0, max: 2, n: 2, input: 1, want: 2},
		{min: 0, max: 2, n: 2, input: 2, want: 2},

		{min: 0, max: 2, n: 3, input: 0, want: 0},
		{min: 0, max: 2, n: 3, input: 1, want: 1},
		{min: 0, max: 2, n: 3, input: 2, want: 2},

		{min: 0, max: 2, n: 4, input: 0, want: 0},
		{min: 0, max: 2, n: 4, input: 1, want: 1},
		{min: 0, max: 2, n: 4, input: 2, want: 2},

		{min: 0, max: 3, n: 2, input: 0, want: 0},
		{min: 0, max: 3, n: 2, input: 1, want: 0},
		{min: 0, max: 3, n: 2, input: 2, want: 3},
		{min: 0, max: 3, n: 2, input: 3, want: 3},

		{min: 0, max: 3, n: 3, input: 0, want: 0},
		{min: 0, max: 3, n: 3, input: 1, want: 1},
		{min: 0, max: 3, n: 3, input: 2, want: 3},
		{min: 0, max: 3, n: 3, input: 3, want: 3},

		{min: 0, max: 3, n: 4, input: 0, want: 0},
		{min: 0, max: 3, n: 4, input: 1, want: 1},
		{min: 0, max: 3, n: 4, input: 2, want: 2},
		{min: 0, max: 3, n: 4, input: 3, want: 3},

		{min: 10, max: 13, n: 4, input: 10, want: 10},
		{min: 10, max: 13, n: 4, input: 11, want: 11},
		{min: 10, max: 13, n: 4, input: 12, want: 12},
		{min: 10, max: 13, n: 4, input: 13, want: 13},
	}

	for _, test := range subtests {
		name := fmt.Sprintf("min=%d max=%d n=%d input=%d", test.min, test.max, test.n, test.input)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sc, err := staircase.New(test.min, test.max, test.n)
			if err != nil {
				t.Fatalf("creating staricase: %v", err)
			}

			got := sc.Transform(test.input)
			if got != test.want {
				t.Fatalf("staircase %v, input %d, step %d, got %d, want %d",
					sc, test.input, sc.Step(test.input), got, test.want)
			}
		})
	}
}
