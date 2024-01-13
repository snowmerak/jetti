package executor

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/snowmerak/jetti/v2/lib/tools"
)

func InstallRegistriesRenew() error {
	if err := tools.ResetTempDir(); err != nil {
		return err
	}

	if err := tools.CloneIfNotExists(); err != nil {
		return err
	}

	return nil
}

func InstallRegistry() error {
	registries, err := tools.GetRegistries()
	if err != nil {
		return err
	}

loop:
	for {
		selectRegistry := &survey.Select{
			Message:     "Select registry to install",
			Options:     registries,
			Description: nil,
		}

		selected := ""
		if err := survey.AskOne(selectRegistry, &selected); err != nil {
			return err
		}

		reg, err := tools.GetRegistryInfo(selected)
		if err != nil {
			return err
		}

		installationConfirm := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure to install %s?\nrepository: %s\ndescription: %s", selected, reg.Repository, reg.Description),
			Default: false,
		}

		ok := false
		if err := survey.AskOne(installationConfirm, &ok); err != nil {
			return err
		}

		if !ok {
			continue loop
		}

		versionInput := &survey.Input{
			Message: "Input version to install",
			Default: "latest",
		}

		version := ""
		if err := survey.AskOne(versionInput, &version); err != nil {
			return err
		}

		if err := tools.InstallRegistry(reg.Repository, version); err != nil {
			return err
		}
	}
}
