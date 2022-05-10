package exec

import (
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

const SchemeExec = "exec"

type UpdaterHandler struct{}

func (h *UpdaterHandler) Match(url *url.URL) bool {
	return url.Scheme == SchemeExec
}

func (u *UpdaterHandler) Update(url *url.URL) (io.WriteCloser, error) {
	path := filepath.Join(url.Host, url.Path)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cmd := exec.Command(absPath)

	if url.Query().Get("env") == "yes" {
		cmd.Env = os.Environ()
	}

	if url.Query().Get("stdout") == "yes" {
		cmd.Stdout = os.Stdout
	}

	if url.Query().Get("stderr") == "yes" {
		cmd.Stderr = os.Stderr
	}

	writer, err := cmd.StdinPipe()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := cmd.Start(); err != nil {
		panic(errors.WithStack(err))
	}

	return writer, nil
}

func NewUpdaterHandler() *UpdaterHandler {
	return &UpdaterHandler{}
}
