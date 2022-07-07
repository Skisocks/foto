package _import

import (
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/skisocks/foto/pkg/cmd/options"
	"github.com/skisocks/foto/pkg/helpers/command"
	"github.com/skisocks/foto/pkg/helpers/files"
	"github.com/skisocks/foto/pkg/helpers/image"
	"github.com/skisocks/foto/pkg/rootcmd"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// PhotoOptions the command line options
type PhotoOptions struct {
	Options
	options.BaseOptions

	Exclude   []string
	Copy      bool
	sourceDir string
}

var (
	avmDesc = command.LongDesc(`
		Imports photos from a external device to a local directory'
`)

	avmExample = command.Examples(`
		%s import photos [source] [destination]
	`)
)

// NewCmdPhoto generates a
func NewCmdPhoto() *cobra.Command {
	photoOptions := &PhotoOptions{
		Options: Options{},
	}

	cmd := &cobra.Command{
		Use:     "photos",
		Short:   "ph",
		Long:    avmDesc,
		Example: fmt.Sprintf(avmExample, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			photoOptions.Cmd = cmd
			photoOptions.Args = args
			err := photoOptions.Run(args)
			command.CheckErr(err)
		},
		Aliases: []string{"av"},
	}
	cmd.Flags().BoolVarP(&photoOptions.Copy, "copy", "c", true, "this copies the photos instead of moving them (default = true)")
	cmd.Flags().StringArrayVarP(&photoOptions.Exclude, "exclude-types", "e", nil, "a list of file extensions to ignore (format as .jpg)")
	return cmd
}

// Run implements this command
func (o *PhotoOptions) Run(args []string) error {
	err := o.BaseOptions.Validate()
	if err != nil {
		return err
	}
	src, dst, err := validateArgs(args)
	if err != nil {
		return err
	}

	paths, err := files.GetPaths(src)
	if err != nil {
		return err
	}

	act := "Moving"
	if o.Copy {
		act = "Copying"
	}

	bar := progressbar.Default(int64(len(paths)), fmt.Sprintf("%s images", act))
	for _, path := range paths {
		img, err := image.Read(path)
		if err != nil {
			return err
		}

		if o.Exclude != nil && o.isExtensionExcluded(img.GetExt()) {
			continue
		}

		if o.Copy {
			err = img.Copy(dst)
			if err != nil {
				return err
			}
		} else {
			err = img.Move(dst)
			if err != nil {
				return err
			}
		}

		bar.Add(1)
	}
	return nil
}

func validateArgs(args []string) (string, string, error) {
	if len(args) != 2 {
		return "", "", errors.New("please enter the source and destination directories")
	}
	for _, arg := range args {
		stat, err := os.Stat(arg)
		if err != nil {
			return "", "", err
		}
		if !stat.IsDir() {
			return "", "", fmt.Errorf("path %s is not a directory", arg)
		}
	}

	isEmpty, err := files.IsEmpty(args[1])
	if err != nil {
		return "", "", err
	}
	if isEmpty {
		return "", "", fmt.Errorf("path %s contains no files", args[0])
	}

	return args[0], args[1], nil
}

func (o *PhotoOptions) isExtensionExcluded(imgExt string) bool {
	for _, ext := range o.Exclude {
		if imgExt == strings.ToLower(ext) {
			return true
		}
	}
	return false
}
