package kudos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joshlf/kudos/lib/config"
	"github.com/joshlf/kudos/lib/db"
	"github.com/joshlf/kudos/lib/kudos/internal"
)

// InitCourse initializes the course specified by ctx.
// It assumes that the course root already exists.
// It logs messages about what files and directories
// it is creating at the verbose level.
func InitCourse(ctx *Context) (err error) {
	conf := internal.MustAsset(filepath.Join("example", config.CourseConfigFileName))
	assign := internal.MustAsset(filepath.Join("example", config.AssignmentDirName, "assignment.sample"))

	fi, err := os.Stat(ctx.CourseRoot())
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("course root does not exist")
		}
		return fmt.Errorf("could not stat course root: %v", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("course root exists but is not directory")
	}

	logAndMkdir := func(path string, perms os.FileMode) error {
		ctx.Verbose.Printf("creating %v\n", path)
		err := os.Mkdir(path, perms|os.ModeDir)
		if err != nil {
			return err
		}
		// in case permissions are masked out by umask
		return os.Chmod(path, perms|os.ModeDir)
	}
	logAndWriteNewFile := func(path string, perms os.FileMode, contents []byte) error {
		ctx.Verbose.Printf("creating %v\n", path)
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perms)
		if err != nil {
			return err
		}
		// in case permissions are masked out by umask
		err = os.Chmod(path, perms)
		if err != nil {
			return err
		}
		_, err = f.Write(conf)
		if err != nil {
			return err
		}
		return f.Sync()
	}

	err = logAndMkdir(ctx.CourseKudosDir(), config.KudosDirPerms)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("course already initialized (%v already exists)", ctx.CourseKudosDir())
		}
		return err
	}
	err = logAndMkdir(ctx.CourseAssignmentDir(), config.AssignmentDirPerms)
	if err != nil {
		return
	}
	err = logAndMkdir(ctx.CourseHandinDir(), config.HandinDirPerms)
	if err != nil {
		return
	}
	err = logAndMkdir(ctx.CourseSavedHandinsDir(), config.SavedHandinsDirPerms)
	if err != nil {
		return
	}
	err = logAndMkdir(ctx.CourseHooksDir(), config.HooksDirPerms)
	if err != nil {
		return
	}
	err = logAndMkdir(ctx.CourseDBDir(), config.DBDirPerms)
	if err != nil {
		return
	}
	err = logAndMkdir(ctx.CoursePubDBDir(), config.PubDBDirPerms)
	if err != nil {
		return
	}

	ctx.Verbose.Println("initializing database")
	err = db.Init(NewDB(), ctx.CourseDBDir())
	if err != nil {
		return
	}

	ctx.Verbose.Println("initializing public database")
	err = db.Init(NewPubDB(), ctx.CoursePubDBDir())
	if err != nil {
		return
	}

	err = logAndWriteNewFile(ctx.CourseConfigFile(), config.CourseConfigFilePerms, conf)
	if err != nil {
		return
	}
	path := filepath.Join(ctx.CourseAssignmentDir(), "assignment.sample")
	return logAndWriteNewFile(path, 0664, assign)
}
