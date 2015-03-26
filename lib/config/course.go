package config

type CourseConfig struct {
	Name             string `toml:"course-name"`
	TaGroup          string `toml:"ta-group"`
	StudentGroup     string `toml:"student-group"`
	ShortDescription string `toml:"short-description"`
	LongDescription  string `toml:"long-description"`
}

func DefaultCourseConfig() CourseConfig {
	return CourseConfig{
		Name:             "cs0",
		TaGroup:          "cs0tas",
		StudentGroup:     "cs0students",
		ShortDescription: "Test CS course",
		LongDescription: `This is a test file providing an example for
how course configs will look`,
	}
}
