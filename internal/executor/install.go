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

const exit = "[exit]"

func InstallRegistry() error {
	registries, err := tools.GetRegistries()
	if err != nil {
		return err
	}
	registries = append([]string{exit}, registries...)

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

		if selected == exit {
			break
		}

		reg, err := tools.GetRegistryInfo(selected)
		if err != nil {
			return err
		}

		installationConfirm := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure to install %s?\nrepository: %s\ndescription: %s\n", selected, reg.Repository, reg.Description),
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

	return nil
}

type Candidate struct {
	Repository string
	Version    string
}

func InstallMultipleRegistries() error {
	registries, err := tools.GetRegistries()
	if err != nil {
		return err
	}

	multiSelectRegistry := &survey.MultiSelect{
		Message: "Select registries to install",
		Options: registries,
		Default: nil,
	}

	selected := make([]string, 0)
	if err := survey.AskOne(multiSelectRegistry, &selected); err != nil {
		return err
	}

	candidates := make([]Candidate, 0, len(selected))

	for _, registry := range selected {
		reg, err := tools.GetRegistryInfo(registry)
		if err != nil {
			return err
		}

		sureConfirm := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure to install %s?\nrepository: %s\ndescription: %s\n", registry, reg.Repository, reg.Description),
			Default: true,
		}

		sure := true
		if err := survey.AskOne(sureConfirm, &sure); err != nil {
			return err
		}

		if !sure {
			continue
		}

		versionInput := &survey.Input{
			Message: fmt.Sprintf("Input version to install %s", registry),
			Default: "latest",
		}

		version := ""
		if err := survey.AskOne(versionInput, &version); err != nil {
			return err
		}

		candidates = append(candidates, Candidate{
			Repository: reg.Repository,
			Version:    version,
		})
	}

	for _, candidate := range candidates {
		if err := tools.InstallRegistry(candidate.Repository, candidate.Version); err != nil {
			return err
		}
	}

	return nil
}
