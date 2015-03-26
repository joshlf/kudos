package config

type CourseConfig struct {
	Name             string `toml:"course-name"`
	TaGroup          string `toml:"ta-group"`
	StudentGroup     string `toml:"student-group"`
	ShortDescription string `toml:"short-description"`
	LongDescription  string `toml:"long-description"`
}
