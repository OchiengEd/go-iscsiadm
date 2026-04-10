package iscsiadm

import (
	"reflect"
	"testing"
)

func Test_parseDiscoverInfo(t *testing.T) {
	tt := []struct {
		name string
		args string
		want int
	}{
		{
			name: "single discovered target",
			args: "192.168.20.90:3260,1 iqn.2024-04.xyz.mtaani.zeus:csi-hibernate-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864",
			want: 1,
		},
		{
			name: "multiple discovered targets",
			args: `192.168.20.90:3260,1 iqn.2024-04.xyz.mtaani.zeus:csi-hibernate-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864
192.168.20.90:3260,1 iqn.2024-04.xyz.mtaani.zeus:csi-hibernate-pvc-be803236-44b9-4352-a0f6-73762cf8d6c2`,
			want: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := parseDiscoverInfo([]byte(tc.args))
			if len(got) != tc.want {
				t.Errorf("Wanted %d targets but, got %d\n", tc.want, len(got))
			}
		})
	}
}

// tcp: [1] 192.168.1.50:3260,1 iqn.2003-01.com.example:storage.target1 (non-flash)
// tcp: [2] 192.168.1.52:3260,1 iqn.2003-01.com.example:storage.target2 (non-flash)
// ptotocol: tcp
// session id: 1
// portal: 192.168.1.50:3260
// portal group tag: 1
// iqn
func Test_parseSessionData(t *testing.T) {
	tt := []struct {
		name string
		args string
		want int
	}{
		{
			name: "single iscsi session",
			args: `tcp: [1] 192.168.1.50:3260,1 iqn.2003-01.com.example:storage.target1 (non-flash)`,
			want: 1,
		},
		{
			name: "two iscsi sessions",
			args: `tcp: [1] 192.168.1.50:3260,1 iqn.2003-01.com.example:storage.target1 (non-flash)
tcp: [2] 192.168.1.52:3260,1 iqn.2003-01.com.example:storage.target2 (non-flash)`,
			want: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := parseSessionData([]byte(tc.args))
			if len(got) != tc.want {
				t.Errorf("wanted %d iSCSI sessions but got %d\n", tc.want, len(got))
			}
		})
	}
}

func Test_parseLoginData(t *testing.T) {
	tt := []struct {
		name string
		args string
		want Target
	}{
		{
			name: "successful login",
			args: "Login to [iface: default, target: iqn.2024-04.com.example.io:csi-pool1-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864, portal: 192.168.20.90,3260] successful.",
			want: Target{
				name:   "iqn.2024-04.com.example.io:csi-pool1-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864",
				portal: "192.168.20.90:3260",
			},
		},
		// {
		// 	name: "failed iscsi login",
		// 	args: "Login to [iface: default, target: iqn.2024-04.com.example.io:csi-pool1-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864, portal: 192.168.20.90,3260] successful.",
		// },
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := parseLoginData([]byte(tc.args))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("wanted to login to target, %s but got  %s", tc.want.name, got.name)
			}
		})
	}
}

func Test_parseLogoutData(t *testing.T) {
	tt := []struct {
		name string
		args string
		want Target
	}{
		{
			name: "successful logout",
			args: `Logging out of session [sid: 1, target: iqn.2024-04.com.example.io:csi-pool1-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864, portal: 192.168.20.90,3260]
Logout of [sid: 1, target: iqn.2024-04.com.example.io:csi-pool1-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864, portal: 192.168.20.90,3260] successful.`,
			want: Target{
				name:   "iqn.2024-04.com.example.io:csi-pool1-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864",
				portal: "192.168.20.90:3260",
			},
		},
		// {
		// 	name: "failed iscsi logout",
		// 	args: "Login to [iface: default, target: iqn.2024-04.com.example.io:csi-pool1-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864, portal: 192.168.20.90,3260] successful.",
		// },
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := parseLogoutData([]byte(tc.args))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("expected %v but got, %v\n", tc.want, got)
			}
		})
	}
}
