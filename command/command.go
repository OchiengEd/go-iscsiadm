package command

// LoginCmd returns command to create a session for a given target and portal
func LoginCmd(useNsenter bool, target, portal string) *Response {
	cmd := []string{"iscsiadm", "-m", "node", "-T", target, "-p", portal, "--login"}
	return prefixCmd(useNsenter, cmd)
}

// LogoutCmd returns command to logout an active target session
func LogoutCmd(useNsenter bool, target, portal string) *Response {
	cmd := []string{"iscsiadm", "-m", "node", " -T", target, "-p", portal, "--logout"}
	return prefixCmd(useNsenter, cmd)
}

// DiscoverCmd sends command to return available targets from a given portal
func DiscoverCmd(useNsenter bool, portal string) *Response {
	cmd := []string{"iscsiadm", "-m", "discoverydb", "-t", "sendtargets", "-p", portal, "--discover"}
	return prefixCmd(useNsenter, cmd)
}

// RemoveCmd returns command to delete target from discovery database
func RemoveCmd(useNsenter bool, target, portal string) *Response {
	cmd := []string{"iscsiadm", "-m", "node", "-T", target, "-p", portal, "-o", "delete"}
	return prefixCmd(useNsenter, cmd)
}

// RescanCmd returns command to delete target from discovery database
func RescanCmd(useNsenter bool, target, portal string) *Response {
	cmd := []string{"iscsiadm", "-m", "node", "-T", target, "-p", portal, "--rescan"}
	return prefixCmd(useNsenter, cmd)
}

// SessionsCmd sends command to list all current sessions
func SessionsCmd(useNsenter bool) *Response {
	cmd := []string{"iscsiadm", "-m", "session"}
	return prefixCmd(useNsenter, cmd)
}

// ListCmd sends command to return list of discovered targets
func ListCmd(useNsenter bool) *Response {
	cmd := []string{"iscsiadm", "-m", "node"}
	return prefixCmd(useNsenter, cmd)
}

func prefixCmd(useNsenter bool, cmd []string) *Response {
	if !useNsenter {
		return &Response{
			command: cmd,
		}
	}

	nsenterCmd := []string{"nsenter",
		"--mount=/proc/1/ns/mnt",
		"--ipc=/proc/1/ns/ipc",
		"--net=/proc/1/ns/net",
		"--",
	}
	nsenterCmd = append(nsenterCmd, cmd...)

	return &Response{
		command: nsenterCmd,
	}
}

type Response struct {
	command []string
}

func (c *Response) Command() string {
	return c.command[0]
}

func (c *Response) Args() []string {
	return c.command[1:]
}
