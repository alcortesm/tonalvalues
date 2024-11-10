package staircase

import (
	"fmt"
)

type Staircase struct {
	min, max, n int

	stepWidth  float64
	stepHeight float64
}

func New(min, max, n uint) (*Staircase, error) {
	if max < min {
		return nil, fmt.Errorf("min (%d) is bigger than max (%d)", min, max)
	}

	return &Staircase{
		min:        int(min),
		max:        int(max),
		n:          int(n),
		stepWidth:  float64(max-min) / float64(n),
		stepHeight: float64(max-min) / float64(n-1),
	}, nil
}

func (s *Staircase) Transform(i int) int {
	switch {
	case i <= s.min:
		return s.min
	case i >= s.max:
		return s.max
	default: // i is between min and max
		return s.min + int(float64(s.Step(i))*s.stepHeight)
	}
}

// the Step the i value belongs to.
// steps go from 0, 1, 2 ... n-1.
//
// only valid for i between s.min and s.max.
func (s *Staircase) Step(i int) int {
	result := int(float64(i-s.min) / s.stepWidth)
	return result
}

func (s *Staircase) String() string {
	return fmt.Sprintf("[min=%d, max=%d, n=%d, stepWidth=%f, stepHeight=%f]",
		s.min, s.max, s.n, s.stepWidth, s.stepHeight)
}
