package command

import "testing"

func Test_accessHostNamespaces(t *testing.T) {
	tt := []struct {
		name string
		args []string
		want bool
	}{
		{
			name: "access mount namespace",
			args: []string{"/proc/1/ns/mnt"},
			want: false,
		},
		{
			name: "access ipc namespace",
			args: []string{"/proc/1/ns/ipc"},
			want: false,
		},
		{
			name: "access network namespace",
			args: []string{"/proc/1/ns/net"},
			want: false,
		},
		{
			name: "access multiple namespace",
			args: []string{"/proc/1/ns/ipc", "/proc/1/ns/net", "/proc/1/ns/mnt"},
			want: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := accessHostNamespaces(tc.args...)
			if tc.want != got {
				t.Logf("Got %t but wanted %t for %v namespace(s)\n", got, tc.want, tc.args)
			}
		})
	}
}
