package iscsiadm

import (
	"fmt"
	"log"
	"regexp"
)

func parseDiscoverInfo(data []byte) []Target {
	re := regexp.MustCompile(`(?m)^(?P<portal>[\d:.]+),\d+\s+(?P<target>iqn.[0-9A-Za-z:.-]+)$`)
	var targets []Target
	for _, match := range re.FindAllSubmatch(data, -1) {
		targets = append(targets, Target{
			name:   fmt.Sprintf("%s", match[re.SubexpIndex("target")]),
			portal: fmt.Sprintf("%s", match[re.SubexpIndex("portal")]),
		})
	}

	return targets
}

func parseLoginData(data []byte) (Target, bool) {
	re := regexp.MustCompile(`(?m)^Login[\s\w\[:]+,\s+target:\s(?P<target>iqn.[\w:.-]+),\s+portal:\s+(?P<portal>[\d.]{7,19}),(?P<port>\d+)\]\s+(?P<result>\w+).`)
	match := re.FindSubmatch(data)
	if len(match) == 0 {
		return Target{}, false
	}
	return Target{
		name: fmt.Sprintf("%s", match[re.SubexpIndex("target")]),
		portal: fmt.Sprintf("%s:%s",
			match[re.SubexpIndex("portal")],
			match[re.SubexpIndex("port")],
		),
	}, true
}

func parseLogoutData(data []byte) (Target, bool) {
	re := regexp.MustCompile(`(?m)^Logout\s[\[\s\w]+sid:\s+(?P<session>\d+),\s+target:\s(?P<target>iqn.[\w:.-]+),\s+portal:\s+(?P<portal>[\d.]{7,19}),(?P<port>\d+)\]\s+(?<ok>\w+).`)
	match := re.FindSubmatch(data)
	if len(match) == 0 {
		return Target{}, false
	}
	log.Printf("%s", match[re.SubexpIndex("target")])
	return Target{
		name: fmt.Sprintf("%s", match[re.SubexpIndex("target")]),
		portal: fmt.Sprintf("%s:%s",
			match[re.SubexpIndex("portal")],
			match[re.SubexpIndex("port")],
		),
	}, true
}

func parseSessionData(data []byte) []Session {
	_ = `tcp: [1] 192.168.20.90:3260,1 iqn.2024-04.com.example.io:csi-pool1-pvc-6a26fc18-5ff5-4737-9462-9cff3ff91864 (non-flash)`
	if len(data) == 0 {
		return []Session{}
	}
	re := regexp.MustCompile(`(?m)^[a-z]{3}:\s+\[(?P<session>[0-9.]+)\]\s+(?P<portal>[0-9:.]+),[0-9]+\s(?P<target>iqn.[0-9]{4}-[0-9]{2}\.[A-Za-z0-9:.-]+)`)
	var sessions []Session
	matches := re.FindAllSubmatch(data, -1)
	for _, match := range matches {
		sessions = append(sessions,
			Session{
				target:  fmt.Sprintf("%s", match[re.SubexpIndex("target")]),
				portal:  fmt.Sprintf("%s", match[re.SubexpIndex("portal")]),
				session: fmt.Sprintf("%s", match[re.SubexpIndex("session")]),
			},
		)
	}

	return sessions
}
