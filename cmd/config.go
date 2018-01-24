package cmd

import (
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type ConfigCmd struct {
	ui       boshui.UI
	director boshdir.Director
}

func NewConfigCmd(ui boshui.UI, director boshdir.Director) ConfigCmd {
	return ConfigCmd{ui: ui, director: director}
}

func (c ConfigCmd) Run(opts ConfigOpts) error {
	if opts == (ConfigOpts{}) {
		return bosherr.Error("Either <ID> or parameters --type and --name must be provided")
	}

	if opts.Args.ID != "" && (opts.Type != "" || opts.Name != "") {
		return bosherr.Error("Can only specify one of ID or parameters '--type' and '--name'")
	}

	if (opts.Args.ID == "" && opts.Type != "" && opts.Name == "") || (opts.Args.ID == "" && opts.Name != "" && opts.Type == "")  {
		return bosherr.Error("Need to specify both parameters '--type' and '--name'")
	}

	config, err := c.director.LatestConfig(opts.Type, opts.Name)
	if err != nil {
		return err
	}

	c.ui.PrintBlock([]byte(config.Content))
	return nil
}
