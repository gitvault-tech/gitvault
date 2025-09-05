// Copyright 2025 GitVault Technologies. All rights reserved.
// SPDX-License-Identifier: MIT

package forms



// PhantomKitKeyForm form for creating PhantomKit API keys
type PhantomKitKeyForm struct {
	Name        string `binding:"Required;MaxSize(100)"`
	Description string `binding:"MaxSize(500)"`
}
