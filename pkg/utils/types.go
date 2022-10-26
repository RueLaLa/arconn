package utils

type Args struct {
	Profile, Target, Command string
}

type Target struct {
	ResolvedName, SessionInfo, Type string
	Resolved                        bool
}
