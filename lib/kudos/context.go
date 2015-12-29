package kudos

import (
	"path/filepath"

	"github.com/joshlf/kudos/lib/config"
	"github.com/joshlf/kudos/lib/log"
)

type Context struct {
	GlobalConfig *GlobalConfig
	CourseCode   string // if the empty string, considered to be unset
	Course       *Course
	*log.Logger
}

func (c *Context) CourseRoot() string     { return CourseCodeToPath(c.CourseCode, c.GlobalConfig) }
func (c *Context) CourseKudosDir() string { return filepath.Join(c.CourseRoot(), config.KudosDirName) }
func (c *Context) CourseConfigFile() string {
	return filepath.Join(c.CourseKudosDir(), config.CourseConfigFileName)
}

func (c *Context) CourseHandinDir() string {
	return filepath.Join(c.CourseKudosDir(), config.HandinDirName)
}

func (c *Context) CourseAssignmentDir() string {
	return filepath.Join(c.CourseKudosDir(), config.AssignmentDirName)
}

func (c *Context) AssignmentHandinDir(code string) string {
	return filepath.Join(c.CourseHandinDir(), code)
}

func (c *Context) HandinHandinDir(assignment, handin string) string {
	return filepath.Join(c.CourseHandinDir(), assignment, handin)
}

func (c *Context) UserAssignmentHandinFile(code, uid string) string {
	return filepath.Join(c.AssignmentHandinDir(code), uid, config.HandinFileName)
}

func (c *Context) UserHandinHandinFile(assignment, handin, uid string) string {
	return filepath.Join(c.HandinHandinDir(assignment, handin), uid, config.HandinFileName)
}
