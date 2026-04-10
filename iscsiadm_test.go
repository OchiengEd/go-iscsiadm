package iscsiadm_test

import (
	"context"
	"strings"
	"testing"

	"github.com/OchiengEd/go-iscsiadm"
	"github.com/OchiengEd/go-iscsiadm/command"
)

var _ command.Runner = (*MockRunner)(nil)

type MockRunner struct{}

func (r *MockRunner) Run(ctx context.Context, command string, args ...string) ([]byte, error) {
	return []byte(strings.Join(args, " ")), nil
}

func TestLogin(t *testing.T) {
	tt := []struct {
		name    string
		args    iscsiadm.LoginRequest
		want    iscsiadm.Device
		wantErr bool
	}{
		{
			name: "with target IQN and portal",
			args: iscsiadm.LoginRequest{
				Portal:    "192.168.100.5:3260",
				TargetIQN: "iqn.1993-08.org.debian.iscsi:01:107dc7e4254a",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "missing target IQN",
			args: iscsiadm.LoginRequest{
				Portal: "192.168.100.5:3260",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "missing portal",
			args: iscsiadm.LoginRequest{
				TargetIQN: "iqn.1993-08.org.debian.iscsi:01:107dc7e4254a",
			},
			want:    "",
			wantErr: true,
		},
	}

	ctrl := iscsiadm.SystemController{
		Runner: &MockRunner{},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ctrl.Login(context.Background(), &tc.args)
			if got != tc.want {
			}
			if (err != nil) != tc.wantErr {
				t.Errorf("wanted presense of error to be %t but, got %t\n.", tc.wantErr, err)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	tt := []struct {
		name    string
		args    iscsiadm.LogoutRequest
		want    bool
		wantErr bool
	}{
		{
			name: "with target IQN and portal",
			args: iscsiadm.LogoutRequest{
				Portal:    "192.168.100.5:3260",
				TargetIQN: "iqn.1993-08.org.debian.iscsi:01:107dc7e4254a",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "missing target IQN",
			args: iscsiadm.LogoutRequest{
				Portal: "192.168.100.5:3260",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "missing portal",
			args: iscsiadm.LogoutRequest{
				TargetIQN: "iqn.1993-08.org.debian.iscsi:01:107dc7e4254a",
			},
			want:    false,
			wantErr: true,
		},
	}

	ctrl := iscsiadm.SystemController{
		Runner: &MockRunner{},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ctrl.Logout(context.Background(), &tc.args)
			if got != tc.want {
			}
			if (err != nil) != tc.wantErr {
				t.Errorf("wanted presense of error to be %t but, got %t\n.", tc.wantErr, err)
			}
		})
	}
}
