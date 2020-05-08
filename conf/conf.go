package conf

var (
	// LenStackBuf LenStackBuf
	LenStackBuf = 4096

	// log
	// LogLevel string
	// LogPath  string
	// LogFlag  int

	// ConsolePort ConsolePort
	ConsolePort int
	// ConsolePrompt ConsolePrompt
	ConsolePrompt string = "Leaf# "
	// ProfilePath ProfilePath
	ProfilePath string

	// cluster

	// ListenAddr ListenAddr
	ListenAddr string
	// ConnAddrs ConnAddrs
	ConnAddrs []string
	// PendingWriteNum PendingWriteNum
	PendingWriteNum int
)
