package main

import (
	"errors"
	"fmt"
	"os"
)

type config struct {
	imagePath string
}

func loadConfig() (*config, error) {
	if len(os.Args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments, want 1 got %d", len(os.Args)-1)
	}

	imagePath := os.Args[1]
	if imagePath == "" {
		return nil, errors.New("empty image path")
	}

	return &config{
		imagePath: imagePath,
	}, nil
}
