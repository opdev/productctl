package file

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/opdev/productctl/internal/logger"
)

var ErrBackupError = errors.New("error creating backup")

// LazyOverwriter is an io.Writer that will wait to open and truncate a file
// until the exact point of writing data. If specified, a backup will be created
// with the original contents before overwriting.
//
// Otherwise, the os.OpenFile call will truncate a given file's contents when
// called. For this reason, all other open operations to Filename should use
// read-only operations.
type LazyOverwriter struct {
	Filename string
	DoBackup bool
	// OptionalLogger is a logger that will emit information about file
	// writing operations. If not set, logs will be discarded.
	OptionalLogger *slog.Logger

	// BackupFilenameGenFn is a configurable string generator, taking an input
	// string and returning a modified representation.
	BackupFilenameGenFn BackupNameGenerator
}

func (w *LazyOverwriter) logger() *slog.Logger {
	if w.OptionalLogger == nil {
		return logger.DiscardingLogger()
	}

	return w.OptionalLogger
}

// nameGenFn returns the in-instance BackupFilenameGenFn, or a default if
// none has been provided.
func (w *LazyOverwriter) nameGenFn() BackupNameGenerator {
	if w.BackupFilenameGenFn == nil {
		return prependWithSecondsSinceEpoch
	}

	return w.BackupFilenameGenFn
}

// CreateBackup will create a backup file for the given w.Filename using
// w.BackupFilenameGenFn to produce the output file. Filenames can be full paths. This
// function will ensure that only basenames are passed to BackupFilenameGenFn.
func (w *LazyOverwriter) CreateBackup() error {
	newNameFn := w.nameGenFn()

	originalFile, err := os.Open(w.Filename)
	if err != nil {
		return err
	}

	newBaseName := newNameFn(filepath.Base(w.Filename))
	newName := filepath.Join(filepath.Dir(w.Filename), newBaseName)
	w.logger().Debug("writing backup file before overwriting original", "backup", newName, "original", w.Filename)
	newFile, err := os.OpenFile(newName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, originalFile)
	return err
}

func (w *LazyOverwriter) Write(p []byte) (int, error) {
	if w.DoBackup {
		w.logger().Debug("backup before overwriting was requested")
		err := w.CreateBackup()
		if err != nil {
			return 0, err
		}
		w.logger().Debug("backup completed successfully")
	}

	w.logger().Debug("overwriting original file")
	f, err := os.OpenFile(w.Filename, os.O_RDWR|os.O_TRUNC, 0o400)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return fmt.Fprint(f, string(p))
}

// prependWithSecondsSinceEpoch returns string s with the seconds since epoch at
// call time, using a dot separator.
func prependWithSecondsSinceEpoch(s string) string {
	generatedVal := strconv.FormatInt(time.Now().Unix(), 10)
	return strings.Join([]string{generatedVal, s}, ".")
}

// Ensure this meets the function signature over time.
var _ BackupNameGenerator = prependWithSecondsSinceEpoch

// BackupNameGenerator describes the function to accept an originalBaseName and
// produce a newBaseName.
type BackupNameGenerator = func(originalBaseName string) (newBaseName string)
