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
	"path/filepath"
	"strings"
)

// PhotoOptions the command line options
type PhotoOptions struct {
	Options
	options.BaseOptions

	Exclude       []string
	Copy          bool
	DirFormat     string
	AddDateToName bool
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
	cmd.Flags().StringVarP(&photoOptions.DirFormat, "dir-format", "d", "", "dir formatting using time.Format")
	cmd.Flags().BoolVarP(&photoOptions.AddDateToName, "add-date-to-name", "n", false, "Adds the date that the photo was taken to the name")
	return cmd
}

// Run implements this command
func (o *PhotoOptions) Run(args []string) error {
	err := o.BaseOptions.Validate()
	if err != nil {
		return err
	}
	source, destination, err := validateArgs(args)
	if err != nil {
		return err
	}

	paths, err := files.GetPaths(source)
	if err != nil {
		return err
	}

	// Remove any excluded file types
	paths = o.filterPaths(paths)

	act := "Moving"
	if o.Copy {
		act = "Copying"
	}

	bar := progressbar.Default(int64(len(paths)), fmt.Sprintf("%s images", act))
	errs := make(map[string]error)
	for _, path := range paths {
		bar.Add(1)
		dst := destination

		img, err := image.Read(path)
		if err != nil {
			errs[path] = err
			continue
		}

		if o.Exclude != nil && o.isExtensionExcluded(files.GetExtension(path)) {
			continue
		}

		if o.DirFormat != "" {
			takenTime, err := img.GetDateTime()
			if err != nil {
				errs[path] = err
				continue
			}
			dst = filepath.Join(dst, takenTime.Format(o.DirFormat))
		}

		imageName := img.GetName()
		if o.AddDateToName {
			takenTime, err := img.GetDateTime()
			if err != nil {
				errs[path] = err
				continue
			}
			imageName = fmt.Sprintf("%s_%s", takenTime.Format("20060102"), imageName)
		}

		dst = filepath.Join(dst, imageName)
		if o.Copy {
			err = img.Copy(dst)
			if err != nil {
				errs[path] = err
				continue
			}
		} else {
			err = img.Move(dst)
			if err != nil {
				errs[path] = err
				continue
			}
		}
	}

	if len(errs) > 0 {
		fmt.Printf("%d errors occured when importing images\n", len(errs))
		for k, v := range errs {
			fmt.Printf("%s:\n%v\n", k, v)
		}
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

	isEmpty, err := files.IsEmpty(args[0])
	if err != nil {
		return "", "", err
	}
	if isEmpty {
		return "", "", fmt.Errorf("path %s contains no files", args[0])
	}

	return args[0], args[1], nil
}

func (o *PhotoOptions) filterPaths(paths []string) []string {
	var filteredPaths []string
	for _, path := range paths {
		if o.Exclude != nil && o.isExtensionExcluded(files.GetExtension(path)) {
			continue
		}
		filteredPaths = append(filteredPaths, path)
	}
	return filteredPaths
}

func (o *PhotoOptions) isExtensionExcluded(imgExt string) bool {
	for _, ext := range o.Exclude {
		if imgExt == strings.ToLower(ext) {
			return true
		}
	}
	return false
}
