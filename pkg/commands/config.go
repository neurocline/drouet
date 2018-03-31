// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"reflect"
	"sort"

	"github.com/neurocline/drouet/pkg/core"
	//"github.com/neurocline/viper"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// Build "hugo config" command.
func buildHugoConfigCmd(h *hugoCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Print the site configuration",
		Long:  `Print the site configuration, both default and custom settings.`,
		RunE:  h.config,
	}

	cmd.Flags().StringP("source", "s", "", "filesystem path to read files relative from")
	cmd.Flags().StringP("theme", "t", "", "theme to use (located in /themes/THEMENAME/)")
	cmd.Flags().String("themesDir", "", "filesystem path to themes directory")

	return cmd
}

// ----------------------------------------------------------------------------------------------

func (h *hugoCmd) config(cmd *cobra.Command, args []string) error {
	//fmt.Println("hugo config - print config goes here")

	// Load config
	var err error
	h.Config, err = core.InitializeConfig(h.Hugo, cmd)
	if err != nil {
		return err
	}

	// Show all config organized by origin
	allSettings := h.Config.AllSettings()
	allOrigins := h.Config.AllSettingsLevels()

	// Keys set by flags
	var keysOverride []string
	var keysFlags []string
	var keysConfig []string
	var keysDefault []string

	for k := range allOrigins {
		switch allOrigins[k] {
		case "override":
			keysOverride = append(keysFlags, k)
		case "flag":
			keysFlags = append(keysFlags, k)
		case "config":
			keysConfig = append(keysConfig, k)
		default:
			keysDefault = append(keysDefault, k)
		}
	}
	fn := func(tag string, keys []string) {
		if len(keys) > 0 {
			jww.FEEDBACK.Printf("config from %s:\n", tag)
			sort.Strings(keys)
			for _, v := range keys {
				jww.FEEDBACK.Printf("  %s\n", v)
			}
		}
	}
	fn("override", keysOverride)
	fn("flags", keysFlags)
	fn("config", keysConfig)
	fn("default", keysDefault)

	for i, v := range allOrigins {
		jww.FEEDBACK.Printf("%s = %s\n", i, v)
	}

	var separator string
	if allSettings["metadataformat"] == "toml" {
		separator = " = "
	} else {
		separator = ": "
	}

	// Put keys in sorted order (people like it better that way)
	var keys []string
	for k := range allSettings {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tag := ""
		switch allOrigins[k] {
		case "flag":
			tag = "*"
		case "config":
			tag = "!"
		}
		kv := reflect.ValueOf(allSettings[k])
		if kv.Kind() == reflect.String {
			jww.FEEDBACK.Printf("%s%s%s\"%+v\"\n", tag, k, separator, allSettings[k])
		} else {
			jww.FEEDBACK.Printf("%s%s%s%+v\n", tag, k, separator, allSettings[k])
		}
	}

	return nil
}
