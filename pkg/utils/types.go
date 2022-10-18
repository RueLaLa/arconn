package utils

type Args struct {
	Profile, Target string
}

type Target struct {
	ResolvedName, SessionInfo, Type string
	Resolved                        bool
}
