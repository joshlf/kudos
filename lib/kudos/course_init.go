package kudos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joshlf/kudos/lib/config"
	"github.com/joshlf/kudos/lib/db"
	"github.com/joshlf/kudos/lib/kudos/internal"
	"github.com/joshlf/kudos/lib/perm"
)

// InitCourse initializes the course specified by ctx.
// It assumes that the course root already exists.
func InitCourse(ctx *Context) (err error) {
	conf := internal.MustAsset(filepath.Join("example", config.CourseConfigFileName))
	assign := internal.MustAsset(filepath.Join("example", config.AssignmentDirName, "assignment.sample"))

	dirMode := os.ModeDir | perm.Parse("rwxrwxr-x")
	fileMode := perm.Parse("rw-rw-r--")

	ctx.Verbose.Printf("creating %v\n", ctx.CourseKudosDir())
	err = os.Mkdir(ctx.CourseKudosDir(), dirMode)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("course already initialized (%v already exists)", ctx.CourseKudosDir())
		}
		return err
	}

	ctx.Verbose.Printf("creating %v\n", ctx.CourseAssignmentDir())
	err = os.Mkdir(ctx.CourseAssignmentDir(), dirMode)
	if err != nil {
		return
	}

	ctx.Verbose.Printf("creating %v\n", ctx.CourseHandinDir())
	err = os.Mkdir(ctx.CourseHandinDir(), dirMode)
	if err != nil {
		return
	}

	ctx.Verbose.Printf("creating %v\n", ctx.CourseDBDir())
	err = os.Mkdir(ctx.CourseDBDir(), dirMode)
	if err != nil {
		return
	}

	ctx.Verbose.Println("initializing database")
	err = db.Init(NewDB(), ctx.CourseDBDir())
	if err != nil {
		return
	}

	ctx.Verbose.Printf("creating %v\n", ctx.CourseConfigFile())
	file, err := os.OpenFile(ctx.CourseConfigFile(), os.O_CREATE|os.O_EXCL|os.O_WRONLY, fileMode)
	if err != nil {
		return
	}
	_, err = file.Write(conf)
	if err != nil {
		return
	}
	err = file.Sync()
	if err != nil {
		return
	}

	path := filepath.Join(ctx.CourseAssignmentDir(), "assignment.sample")
	ctx.Verbose.Printf("creating %v\n", path)
	file, err = os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, fileMode)
	if err != nil {
		return err
	}
	_, err = file.Write(assign)
	if err != nil {
		return
	}
	err = file.Sync()
	return
}
