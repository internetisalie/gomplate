package main

import (
	"log/slog"
	"os"

	"code.internetisalie.net/gomplate/pkg/inputs"
)

func main() {
	err := inputs.Constants.Generate()
	if err != nil {
		slog.Error("Failed to generate constants", err)
		os.Exit(1)
	}
}
