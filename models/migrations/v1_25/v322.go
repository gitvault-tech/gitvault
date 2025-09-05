// Copyright 2025 GitVault Technologies. All rights reserved.
// SPDX-License-Identifier: MIT

package v1_25

import (
	"xorm.io/xorm"
)

func CreatePhantomKitKeysTable(x *xorm.Engine) error {
	// Create phantomkit_keys table
	type PhantomKitKey struct {
		ID           int64  `xorm:"pk autoincr"`
		UserID       int64  `xorm:"NOT NULL"`
		Name         string `xorm:"NOT NULL"`
		Description  string `xorm:"TEXT"`
		Key          string `xorm:"UNIQUE NOT NULL"`
		CreatedUnix  int64  `xorm:"created"`
		UpdatedUnix  int64  `xorm:"updated"`
		LastUsedUnix int64  `xorm:"last_used"`
	}

	return x.Sync2(new(PhantomKitKey))
}
