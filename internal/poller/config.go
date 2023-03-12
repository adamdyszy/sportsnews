package poller

type ListConfig interface {
	GetListURL() string
	GetListCount() int
	GetListSchedule() string
	GetTeamId() string
}

type DetailsConfig interface {
	GetDetailsURL() string
	GetDetailsSchedule() string
	GetTeamId() string
}

type Config struct {
	TeamId        string `mapstructure:"teamId"`
	RunOnceAtBoot bool   `mapstructure:"runOnceAtBoot"`
	List          struct {
		URL      string `mapstructure:"url"`
		Count    int    `mapstructure:"count"`
		Schedule string `mapstructure:"schedule"`
	} `mapstructure:"list"`
	Details struct {
		URL      string `mapstructure:"url"`
		Schedule string `mapstructure:"schedule"`
	} `mapstructure:"details"`
}

func (c Config) GetListURL() string {
	return c.List.URL
}

func (c Config) GetListCount() int {
	return c.List.Count
}

func (c Config) GetListSchedule() string {
	return c.List.Schedule
}

func (c Config) GetDetailsURL() string {
	return c.Details.URL
}

func (c Config) GetDetailsSchedule() string {
	return c.Details.Schedule
}

func (c Config) GetTeamId() string {
	return c.TeamId
}
