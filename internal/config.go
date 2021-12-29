package internal

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/manifoldco/promptui"
)

type Config struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Dockerfile string `json:"dockerfile"`
	// Image      string `json:"image"`
	// Environment []string `json:"environment"`
	// Command     string   `json:"command"`
	// Restart     string   `json:"restart"`
	// WorkingDir  string   `json:"working_dir"`
	Ports []string `json:"ports"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	err = json.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c *Config) Prompt() error {
	err := c.PromptProjectName()
	if err != nil {
		return err
	}

	err = c.PromptRepository()
	if err != nil {
		return err
	}

	err = c.PromptBranch()
	if err != nil {
		return err
	}

	err = c.PromptImage()
	if err != nil {
		return err
	}

	err = c.PromptPort()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) PromptProjectName() error {
	validate := func(input string) error {
		if input == "" {
			return errors.New("project name cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Project Name",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	c.Name = result
	return nil
}

func (c *Config) PromptRepository() error {
	validate := func(input string) error {
		if input == "" {
			return errors.New("repository cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Repository",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	c.Repository = result
	return nil
}

func (c *Config) PromptBranch() error {
	validate := func(input string) error {
		if input == "" {
			return errors.New("branch cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Branch",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	c.Branch = result
	return nil
}

func (c *Config) PromptImage() error {
	validate := func(input string) error {
		if input == "" {
			return errors.New("path cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Dockerfile path",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	c.Dockerfile = result
	return nil
}

func (c *Config) PromptPort() error {
	validate := func(input string) error {
		if input == "" {
			return errors.New("port cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Port",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	c.Ports = append(c.Ports, result)
	return nil
}
