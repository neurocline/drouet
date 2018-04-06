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

// +build !darwin

package commands

import (
	"github.com/neurocline/drouet/pkg/core"

	"github.com/neurocline/cobra"
)

func buildHugoCheckUlimitCmd(hugo *core.Hugo) *hugoCheckUlimitCmd {
	return nil
}

type hugoCheckUlimitCmd struct {
	*core.Hugo
	cmd *cobra.Command
}

func tweakLimit() {
	// nothing to do
}