//
// Copyright 2019 Chef Software, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	DefaultChefDirectory            = ".chef"
	DefaultChefWorkstationDirectory = ".chef-workstation"
	DefaultFileName                 = "config.toml"
)

// finds the configuration file (default .chef-workstation/config.toml) inside the current
// directory and recursively, plus inside the $HOME directory
func FindChefWorkstationConfigFile() (string, error) {
	return FindConfigFile(DefaultFileName)
}

// finds the provided configuration file name inside the
// current directory, if the file is not there the, it
// looks up inside the users $HOME directory
func FindConfigFile(name string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "unable to detect current directory")
	}

	//debug("searching config in: %s\n", cwd)
	configFile, exists := configFileExistsInsideDefaultDirectories(cwd, name)
	if exists {
		return configFile, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "unable to detect home directory")
	}
	configFile, exists = configFileExistsInsideDefaultDirectories(home, name)
	if exists {
		return configFile, nil
	}

	// @afiune tell the user the paths we tried to find it?
	return "", errors.Errorf("file '%s' not found", name)
}

// verify if a config file exists inside any of the default directories
// => current directories: ['.chef/', '.chef-workstation/']
func configFileExistsInsideDefaultDirectories(dir, name string) (string, bool) {
	// search inside .chef/
	cfgFileFromChefDir := filepath.Join(dir, DefaultChefDirectory, name)
	if _, err := os.Stat(cfgFileFromChefDir); err == nil {
		return cfgFileFromChefDir, true
	}

	// search inside .chef-workstation/
	cfgFileFromChefWSDir := filepath.Join(dir, DefaultChefWorkstationDirectory, name)
	if _, err := os.Stat(cfgFileFromChefWSDir); err == nil {
		return cfgFileFromChefWSDir, true
	}

	return "", false
}
