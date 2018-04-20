// Copyright 2017 The Hugo Authors. All rights reserved.
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

package hugolib

// Site contains all the information relevant for constructing a static
// site.  The basic flow of information is as follows:
//
// 1. A list of Files is parsed and then converted into Pages.
//
// 2. Pages contain sections (based on the file they were generated from),
//    aliases and slugs (included in a pages frontmatter) which are the
//    various targets that will get generated.  There will be canonical
//    listing.  The canonical path can be overruled based on a pattern.
//
// 3. Taxonomies are created via configuration and will present some aspect of
//    the final page and typically a perm url.
//
// 4. All Pages are passed through a template based on their desired
//    layout based on numerous different elements.
//
// 5. The entire collection of files is written to disk.
type Site struct {
	owner *HugoSites

	*PageCollections
}
