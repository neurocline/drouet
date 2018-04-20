// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package source

import (
	"io"
	"os"
	//"path/filepath"
	//"strings"
	"sync"

	"github.com/neurocline/drouet/pkg/core"
	//"github.com/neurocline/drouet/pkg/hugofs"
)

// fileInfo implements the File interface.
var (
	_ File         = (*FileInfo)(nil)
	_ ReadableFile = (*FileInfo)(nil)
)

type File interface {

	// Filename gets the full path and filename to the file.
	Filename() string

	// Path gets the relative path including file name and extension.
	// The directory is relative to the content root.
	Path() string

	// Dir gets the name of the directory that contains this file.
	// The directory is relative to the content root.
	Dir() string

	// Extension gets the file extension, i.e "myblogpost.md" will return "md".
	Extension() string
	// Ext is an alias for Extension.
	Ext() string // Hmm... Deprecate Extension

	// Lang for this page, if `Multilingual` is enabled on your site.
	Lang() string

	// LogicalName is filename and extension of the file.
	LogicalName() string

	// Section is first directory below the content root.
	// For page bundles in root, the Section will be empty.
	Section() string

	// BaseFileName is a filename without extension.
	BaseFileName() string

	// TranslationBaseName is a filename with no extension,
	// not even the optional language extension part.
	TranslationBaseName() string

	// UniqueID is the MD5 hash of the file's path and is for most practical applications,
	// Hugo content files being one of them, considered to be unique.
	UniqueID() string

	FileInfo() os.FileInfo

	String() string

	// Deprecated
	Bytes() []byte
}

// A ReadableFile is a File that is readable.
type ReadableFile interface {
	File
	Open() (io.ReadCloser, error)
}

type FileInfo struct {

	// Absolute filename to the file on disk.
	filename string

	sp *SourceSpec

	fi os.FileInfo

	// Derived from filename
	ext  string // Extension without any "."
	lang string

	name string

	dir                 string
	relDir              string
	relPath             string
	baseName            string
	translationBaseName string
	section             string
	isLeafBundle        bool

	uniqueID string

	lazyInit sync.Once
}

func (fi *FileInfo) Filename() string            { return fi.filename }
func (fi *FileInfo) Path() string                { return fi.relPath }
func (fi *FileInfo) Dir() string                 { return fi.relDir }
func (fi *FileInfo) Extension() string           { return fi.Ext() }
func (fi *FileInfo) Ext() string                 { return fi.ext }
func (fi *FileInfo) Lang() string                { return fi.lang }
func (fi *FileInfo) LogicalName() string         { return fi.name }
func (fi *FileInfo) BaseFileName() string        { return fi.baseName }
func (fi *FileInfo) TranslationBaseName() string { return fi.translationBaseName }

func (fi *FileInfo) Section() string {
	fi.init()
	return fi.section
}

func (fi *FileInfo) UniqueID() string {
	fi.init()
	return fi.uniqueID
}

func (fi *FileInfo) FileInfo() os.FileInfo {
	return fi.fi
}

func (fi *FileInfo) Bytes() []byte {
	// Remove in Hugo 0.38
	core.Deprecated("File", "Bytes", "", true)
	return []byte("")
}

func (fi *FileInfo) String() string { return fi.BaseFileName() }

// We create a lot of these FileInfo objects, but there are parts of it used only
// in some cases that is slightly expensive to construct.
func (fi *FileInfo) init() {
	/*
	fi.lazyInit.Do(func() {
		relDir := strings.Trim(fi.relDir, hugofs.FilePathSeparator)
		parts := strings.Split(relDir, hugofs.FilePathSeparator)
		var section string
		if (!fi.isLeafBundle && len(parts) == 1) || len(parts) > 1 {
			section = parts[0]
		}

		fi.section = section

		fi.uniqueID = helpers.MD5String(filepath.ToSlash(fi.relPath))

	})*/
}

// Open implements ReadableFile.
func (fi *FileInfo) Open() (io.ReadCloser, error) {
	//f, err := fi.sp.PathSpec.Fs.Source.Open(fi.Filename())
	//return f, err
	return nil, nil
}
