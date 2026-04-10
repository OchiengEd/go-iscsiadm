package command

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
)

var _ Runner = (*DefaultRunner)(nil)

type Runner interface {
	Run(ctx context.Context, command string, args ...string) ([]byte, error)
}

type DefaultRunner struct{}

func NewRunner() *DefaultRunner {
	return &DefaultRunner{}
}

func (r *DefaultRunner) Run(ctx context.Context, command string, args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitCode := exitError.ExitCode(); exitCode {
			case ExitCodeSessionExists:
				// ignore iSCSI session exists errors
				return stdout.Bytes(), ErrSessionExists
			case ExitCodeObjectsNotFound:
				return stdout.Bytes(), ErrResourceNotFound
			case ExitCodeLoginFailure:
				return nil, ErrLoginFailed
			case ExitCodeLogoutFailure:
				return nil, ErrLogoutFailed
			case ExitCodeAccessDenied:
				return nil, fmt.Errorf(
					"insufficient OS permissions; %w", ErrPermissionDenied,
				)
			default:
				log.Println("generic error", "stderr", stderr.String(), "exitcode", fmt.Sprintf("%d", exitCode))
				return nil, err
			}
		}

	}
	return stdout.Bytes(), nil
}
