package upload

import (
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/skisocks/foto/pkg/cmd/options"
	"github.com/skisocks/foto/pkg/helpers/command"
	"github.com/skisocks/foto/pkg/helpers/files"
	"github.com/skisocks/foto/pkg/helpers/flickrUtils"
	"github.com/skisocks/foto/pkg/rootcmd"
	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v2"
	"gopkg.in/masci/flickr.v2/photosets"
)

// PhotoOptions the command line options
type PhotoOptions struct {
	Options
	options.BaseOptions

	Key       string
	Secret    string
	AuthToken string
	AlbumName string
	client    *flickr.FlickrClient
}

var (
	avmDesc = command.LongDesc(`
		Uploads photos from a directory to Flickr'
`)

	avmExample = command.Examples(`
		%s upload photos [source]
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
		Aliases: []string{"ph"},
	}
	cmd.Flags().StringVarP(&photoOptions.Key, "key", "k", "", "your Flickr consumer key")
	cmd.Flags().StringVarP(&photoOptions.Secret, "secret", "s", "", "your Flickr consumer secret")
	cmd.Flags().StringVarP(&photoOptions.AlbumName, "album-name", "a", "", "the name of the album you'd like to create")
	return cmd
}

// Run implements this command
func (o *PhotoOptions) Run(args []string) error {
	err := o.validate()
	if err != nil {
		return err
	}

	if o.AuthToken == "" {
		o.client, err = flickrUtils.Init(o.Key, o.Secret)
		if err != nil {
			return err
		}
	}

	isDir, err := files.IsDir(args[0])
	if err != nil {
		return err
	}
	if !isDir {
		id, err := flickrUtils.UploadPhoto(o.client, args[0])
		if err != nil {
			return err
		}
		if o.AlbumName != "" {
			_, err := photosets.Create(o.client, o.AlbumName, "", id)
			if err != nil {
				return err
			}
		}
		return nil
	}

	paths, err := files.GetPaths(args[0])
	if err != nil {
		return err
	}

	var photoSetID string
	bar := progressbar.Default(int64(len(paths)), "Uploading photos")
	for _, f := range paths {
		bar.Add(1)
		photoID, err := flickrUtils.UploadPhoto(o.client, f)
		if err != nil {
			return err
		}

		if o.AlbumName != "" {
			if photoSetID == "" {
				photoSetID, err = flickrUtils.CreatePhotoSet(o.client, o.AlbumName, "", photoID)
				if err != nil {
					return err
				}
			} else {
				_, err = photosets.AddPhoto(o.client, photoSetID, photoID)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (o *PhotoOptions) validate() error {
	err := o.BaseOptions.Validate()
	if err != nil {
		return err
	}

	if o.Key == "" {
		return errors.New("no API key set")
	}
	if o.Secret == "" {
		return errors.New("no secret set")
	}
	return nil
}
