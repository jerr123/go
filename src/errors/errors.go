// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errors implements functions to manipulate errors.

// 实现了操作错误的方法

package errors

import (
	"internal/errinternal"
	"runtime"
)

// New returns an error that formats as the given text.
//
// The returned error contains a Frame set to the caller's location and
// implements Formatter to show this information when printed with details.

// New 将 text 格式化成 error 返回
//
// 返回的 error 包含一个用来定位调用者的集合，并实现了 Formatter，以在打印详细信息的时候显示此信息
func New(text string) error {
	// Inline call to errors.Callers to improve performance.
	// Inline 调用 errors.Callers 以提高性能
	var s Frame
	runtime.Callers(2, s.frames[:])
	return &errorString{text, nil, s}
}
// ps 这里找了半天没找着 error 类型在哪里定义的，最后发现 error 是 go 中的一种内置类型
// 在 Package builtin 里面文件的末尾
// 先看 init
func init() {
	errinternal.NewError = func(text string, err error) error {
		var s Frame // 2019/03/05/02:16 //TODO: 找到 Frame 类型的定义
		runtime.Callers(3, s.frames[:]) // runtime.Callers 定义在 extern.go
		return &errorString{text, err, s}
	}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s     string
	err   error
	frame Frame
}

func (e *errorString) Error() string {
	if e.err != nil {
		return e.s + ": " + e.err.Error()
	}
	return e.s
}

func (e *errorString) FormatError(p Printer) (next error) {
	p.Print(e.s)
	e.frame.Format(p)
	return e.err
}
