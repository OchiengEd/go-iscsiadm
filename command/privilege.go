package command

import (
	"bufio"
	"os"
	"slices"
	"strings"
)

// accessHostNamespaces checks if the container has privileges
// to access the mount, net, ipc namespaces on the host
func accessHostNamespaces(namespaces ...string) bool {
	for _, namespace := range namespaces {
		if _, err := os.Stat(namespace); err != nil {
			return false
		}
	}
	return true
}

// Checks the container can access /etc/iscsi and /var/lib/iscsi
func iscsiBidirectionalMounts() (bool, error) {
	iscsiMounts := []string{"/etc/iscsi", "/var/lib/iscsi"}

	file, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "-")[0]
		line = strings.Trim(line, " ")

		// Structure of mountinfo:
		// [0] mount ID | [4] mount point | [6] browser-stop (separator '-')
		fields := strings.Fields(line)
		mountPoint := fields[3]
		hostPath := fields[4]
		filesystem := fields[6]

		// Ensure the filesystem is showing as "shared:"
		if !strings.Contains(filesystem, "shared:") {
			return false, nil
		}

		// Ensure mountPoint and hostPath are the same
		if !slices.Contains(iscsiMounts, mountPoint) &&
			mountPoint != hostPath {
			return false, nil
		}

	}

	return true, nil
}
