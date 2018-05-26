package linkedlist

import (
	"io"

	"github.com/clipperhouse/typewriter"
)

func init() {
	err := typewriter.Register(NewLinkedListWriter())
	if err != nil {
		panic(err)
	}
}

type LinkedListWriter struct{}

func NewLinkedListWriter() *LinkedListWriter {
	return &LinkedListWriter{}
}

func (sw *LinkedListWriter) Name() string {
	return "linkedlist"
}

func (sw *LinkedListWriter) Imports(t typewriter.Type) (result []typewriter.ImportSpec) {
	// none
	return result
}

func (sw *LinkedListWriter) Write(w io.Writer, t typewriter.Type) error {
	tag, found := t.FindTag(sw)

	if !found {
		// nothing to be done
		return nil
	}

	license := `
// List is a modification of http://golang.org/pkg/container/list/
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
`
	if _, err := w.Write([]byte(license)); err != nil {
		return err
	}

	tmpl, err := templates.ByTag(t, tag)

	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, t); err != nil {
		return err
	}

	return nil
}
