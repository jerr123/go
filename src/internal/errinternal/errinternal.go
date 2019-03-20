// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errinternal

// NewError creates a new error as created by errors.New, but with one
// additional stack frame depth.

// NewError 创建一个由 errors.New 创建的新的 error，并且它会有额外的堆栈深度。

// 定义了 NewError 函数类型 传入 string error, 返回 error
// ps 具体有什么用?为什么要这么定义?暂时不知道，先不管
var NewError func(msg string, err error) error
