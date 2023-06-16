package utils

type Args struct {
	Profile, Target, Command, PortForward, RemoteHost, Vault string
}

type Target struct {
	ResolvedName, SessionInfo, Type, RemoteHost string
	Resolved                                    bool
	PortForwarding                              []string
}
