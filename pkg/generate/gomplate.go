package generate

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Gomplate struct {
	OutputFile string
	Data       any
	Template   string
}

func (g Gomplate) GoModDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cur := cwd
	for {
		if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
			return cwd, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}

		parent := filepath.Dir(cur)
		if parent == cur {
			return "", os.ErrNotExist
		}

		cur = parent
	}
}

func (g Gomplate) Generate() (err error) {
	jsonData, err := json.Marshal(g.Data)
	if err != nil {
		return
	}

	templateFile, err := os.CreateTemp("", "*.gotpl")
	if err != nil {
		return
	}
	defer templateFile.Close()

	_, err = templateFile.Write([]byte(g.Template))
	if err != nil {
		return
	}

	outputFile := g.OutputFile
	if strings.HasPrefix(outputFile, "/") {
		goModDir, err := g.GoModDir()
		if err != nil {
			return err
		}

		outputFile = filepath.Join(goModDir, outputFile[1:])
	}

	cmd := exec.Command(
		"go", "tool", "gomplate",
		"-f", templateFile.Name(),
		"-d", "data=stdin:///in.json",
	)
	cmd.Stdin = bytes.NewBuffer(jsonData)
	cmd.Stderr = os.Stderr
	cmd.Stdout, err = os.Create(outputFile)
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	return nil
}
