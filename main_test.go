package main

import (
	"fmt"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSteps(t *testing.T) {
	t.Parallel()

	subtests := []struct {
		n    uint
		want []uint8
	}{
		{
			n: 2,
			want: slices.Concat(
				slices.Repeat([]uint8{0}, 128),
				slices.Repeat([]uint8{255}, 128),
			),
		},
		{
			n: 3,
			want: slices.Concat(
				slices.Repeat([]uint8{0}, 86),
				slices.Repeat([]uint8{128}, 85),
				slices.Repeat([]uint8{255}, 85),
			),
		},
		{
			n: 4,
			want: slices.Concat(
				slices.Repeat([]uint8{0}, 64),
				slices.Repeat([]uint8{85}, 64),
				slices.Repeat([]uint8{170}, 64),
				slices.Repeat([]uint8{255}, 64),
			),
		},
		{
			n: 5,
			want: slices.Concat(
				slices.Repeat([]uint8{0}, 52),
				slices.Repeat([]uint8{64}, 51),
				slices.Repeat([]uint8{128}, 51),
				slices.Repeat([]uint8{192}, 51),
				slices.Repeat([]uint8{255}, 51),
			),
		},
	}

	for _, test := range subtests {
		name := fmt.Sprintf("n %d", test.n)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := steps(test.n)

			want := [256]uint8{}
			copy(want[:], test.want)

			if diff := cmp.Diff(&want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
