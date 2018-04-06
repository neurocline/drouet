// Copyright 2016n The Hugo Authors. All rights reserved.
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

package parser

const (
	// TODO(bep) Do we really have to export these?

	// HTMLLead identifies the start of HTML documents.
	HTMLLead = "<"
	// YAMLLead identifies the start of YAML frontmatter.
	YAMLLead = "-"
	// YAMLDelimUnix identifies the end of YAML front matter on Unix.
	YAMLDelimUnix = "---\n"
	// YAMLDelimDOS identifies the end of YAML front matter on Windows.
	YAMLDelimDOS = "---\r\n"
	// YAMLDelim identifies the YAML front matter delimiter.
	YAMLDelim = "---"
	// TOMLLead identifies the start of TOML front matter.
	TOMLLead = "+"
	// TOMLDelimUnix identifies the end of TOML front matter on Unix.
	TOMLDelimUnix = "+++\n"
	// TOMLDelimDOS identifies the end of TOML front matter on Windows.
	TOMLDelimDOS = "+++\r\n"
	// TOMLDelim identifies the TOML front matter delimiter.
	TOMLDelim = "+++"
	// JSONLead identifies the start of JSON frontmatter.
	JSONLead = "{"
	// HTMLCommentStart identifies the start of HTML comment.
	HTMLCommentStart = "<!--"
	// HTMLCommentEnd identifies the end of HTML comment.
	HTMLCommentEnd = "-->"
	// BOM Unicode byte order marker
	BOM = '\ufeff'
)
