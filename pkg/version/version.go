package version

var (
	ServiceName string
	GitBranch   string
	GitVersion  string
	BuildTime   string
)

type Info struct {
	ServiceName  string
	GitBranch    string
	GitVersion   string
	BuildTime    string
	RunningState bool
}

func GetVersionInfo() *Info {
	return &Info{
		ServiceName: ServiceName,
		GitBranch:   GitBranch,
		GitVersion:  GitVersion,
		BuildTime:   BuildTime,
	}
}
