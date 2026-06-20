package scanner

type ScanMode string

const (
	ScanFull       ScanMode = "full"
	ScanComponents ScanMode = "components"
	ScanModules    ScanMode = "modules"
	ScanSecurity   ScanMode = "security"
)
