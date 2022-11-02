package utils

type Args struct {
	Profile, Target, Command, PortForward string
}

type Target struct {
	ResolvedName, SessionInfo, Type string
	Resolved                        bool
	PortForwarding                  []string
}
