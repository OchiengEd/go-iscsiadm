package iscsiadm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/OchiengEd/go-iscsiadm/command"
)

type SystemController struct {
	command.Runner
	useNsenter bool
}

type Target struct {
	name   string
	portal string
	lun    int
}

type ControllerOptions func(*SystemController)

func WithCommandRunner(f command.Runner) ControllerOptions {
	return func(sc *SystemController) {
		sc.Runner = f
	}
}

func WithNsenter(b bool) ControllerOptions {
	return func(sc *SystemController) {
		sc.useNsenter = b
	}
}

func New(opts ...ControllerOptions) *SystemController {
	ctrl := new(SystemController)

	for _, opt := range opts {
		opt(ctrl)
	}

	if ctrl.Runner == nil {
		ctrl.Runner = command.NewRunner()
	}

	return ctrl
}

type Device string

type LoginRequest struct {
	TargetIQN string
	Portal    string
}

// Login to an iSCSI target that has been discovered from a portal
func (c *SystemController) Login(ctx context.Context, req *LoginRequest) (Device, error) {
	if req.Portal == "" {
		return "", errors.New("iSCSI portal required")
	}
	if req.TargetIQN == "" {
		return "", errors.New("iSCSI target IQN required")
	}

	// Check if the targetIQN exists in discoveryDB
	resp := command.ListCmd(c.useNsenter)
	listOut, err := c.Run(ctx, resp.Command(), resp.Args()...)
	if err != nil {
		return "", fmt.Errorf("failed listing discovered targets; %w", err)
	}
	if len(listOut) == 0 {
		return "", errors.New("no discovered targets found")
	}

	if !isTargetDiscovered(listOut, req.TargetIQN, req.Portal) {
		return "", fmt.Errorf("target %s, not in discovered nodes", req.TargetIQN)
	}

	resp = command.LoginCmd(c.useNsenter, req.TargetIQN, req.Portal)
	loginOut, err := c.Run(ctx, resp.Command(), resp.Args()...)
	if err != nil {
		return "", fmt.Errorf("login to target failed; %w", err)
	}
	if _, ok := parseLoginData(loginOut); !ok {
		return "", errors.New("unexpected loging failure")
	}

	// Once there is an active session, we can attempt to get the device path
	devicePrefix := fmt.Sprintf(
		"/dev/disk/by-path/ip-%s-iscsi-%s-lun", req.Portal, req.TargetIQN,
	)
	var devicePath Device
	errCh := make(chan error)

	go func() {
		for {
			if devicePath != "" {
				errCh <- nil
			}
			if err := filepath.WalkDir("/dev/disk/by-path",
				func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if strings.Contains(path, devicePrefix) && !d.IsDir() {
						devicePath = realDevicePath(path)
						return filepath.SkipAll
					}
					return nil
				},
			); err != nil {
				errCh <- err
			}
			time.Sleep(20 * time.Millisecond)
		}
	}()

	select {
	case <-ctx.Done():
		return devicePath, fmt.Errorf("timeout getting device path; %w", ctx.Err())
	case err := <-errCh:
		if err != nil {
			return devicePath, err
		}
		return devicePath, nil
	}
}

type LogoutRequest struct {
	Portal    string
	TargetIQN string
}

// Logout an active iSCSI session to a target
func (c *SystemController) Logout(ctx context.Context, req *LogoutRequest) (bool, error) {
	if req.Portal == "" {
		return false, errors.New("iSCSI portal required")
	}
	if req.TargetIQN == "" {
		return false, errors.New("iSCSI target IQN required")
	}

	resp := command.LogoutCmd(c.useNsenter, req.TargetIQN, req.Portal)
	logutOut, err := c.Run(ctx, resp.Command(), resp.Args()...)
	if err != nil {
		return false, fmt.Errorf("iSCSI target logout failed; %w", err)
	}

	if _, ok := parseLogoutData(logutOut); !ok {
		return false, fmt.Errorf("target logout failed")
	}

	return true, nil
}

// Rescan checks for changes on a currently loggied in iSCSI target
func (c *SystemController) Rescan(ctx context.Context, req *LogoutRequest) (bool, error) {
	if req.Portal == "" {
		return false, errors.New("iSCSI portal required")
	}
	if req.TargetIQN == "" {
		return false, errors.New("iSCSI target IQN required")
	}

	resp := command.RescanCmd(c.useNsenter, req.TargetIQN, req.Portal)
	rescanOut, err := c.Run(ctx, resp.Command(), resp.Args()...)
	if err != nil {
		return false, fmt.Errorf("iSCSI target rescan failed; %w", err)
	}

	// TODO: parse rescan target output
	_ = rescanOut

	return true, nil
}

type Session struct {
	target  string
	portal  string
	session string
}

// Sessions returns a list of active sessions on the current system
func (c *SystemController) Sessions(ctx context.Context) ([]Session, error) {
	resp := command.SessionsCmd(c.useNsenter)
	sessionsOut, err := c.Run(ctx, resp.Command(), resp.Args()...)
	if err != nil {
		return []Session{}, fmt.Errorf("getting iSCSI sessions failed; %w", err)
	}
	sessions := parseSessionData(sessionsOut)

	return sessions, nil
}

type DiscoverRequest struct {
	Portal string
}

// Discover performs an iSCSI discovery with sendtargets
func (c *SystemController) Discover(ctx context.Context, req *DiscoverRequest) ([]Target, error) {
	if req.Portal == "" {
		return []Target{}, errors.New("iSCSI portal required")
	}

	resp := command.DiscoverCmd(c.useNsenter, req.Portal)
	stdout, err := c.Run(ctx, resp.Command(), resp.Args()...)
	if err != nil {
		return []Target{}, fmt.Errorf("discovering iSCSI targets failed; %w", err)
	}

	if len(stdout) == 0 {
		return []Target{}, nil
	}
	targets := parseDiscoverInfo(stdout)

	return targets, nil
}

type RemoveRequest struct {
	Portal    string
	TargetIQN string
}

// Remove deletes a specific target from the discovery DB
func (c *SystemController) Remove(ctx context.Context, req *RemoveRequest) (bool, error) {
	if req.Portal == "" {
		return false, errors.New("iSCSI portal required")
	}
	if req.TargetIQN == "" {
		return false, errors.New("iSCSI target IQN required")
	}

	// Delete from discoveryDB
	resp := command.RemoveCmd(c.useNsenter, req.TargetIQN, req.Portal)
	_, err := c.Run(ctx, resp.Command(), resp.Args()...)
	if err != nil {
		return false, fmt.Errorf("removal of iSCSI node failed; %w", err)
	}

	return true, nil
}

func isTargetDiscovered(data []byte, target, portal string) bool {
	for _, line := range bytes.Split(data, []byte("\n")) {
		fields := strings.Fields(string(line))
		if len(fields) != 2 {
			return false
		}
		if strings.Contains(fields[0], portal) &&
			strings.Contains(fields[1], target) {
			return true
		}
	}
	return false
}

func realDevicePath(path string) Device {
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return Device(path)
	}
	return Device(realPath)
}
