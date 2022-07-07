package options

import "github.com/skisocks/foto/pkg/config"

type BaseOptions struct {
	cfg *config.Config
}

func (b *BaseOptions) Validate() (err error) {
	if b.cfg == nil {
		b.cfg, err = config.LoadConfigOrCreateIt()
		if err != nil {
			return err
		}
	}
	return nil
}
