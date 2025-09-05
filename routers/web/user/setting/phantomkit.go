// Copyright 2025 GitVault Technologies. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/templates"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/forms"
)

const (
	tplSettingsPhantomKit templates.TplName = "user/settings/phantomkit"
)

// PhantomKitKey represents a PhantomKit API key
type PhantomKitKey struct {
	ID          int64     `xorm:"pk autoincr"`
	UserID      int64     `xorm:"NOT NULL"`
	Name        string    `xorm:"NOT NULL"`
	Description string    `xorm:"TEXT"`
	Key         string    `xorm:"UNIQUE NOT NULL"`
	CreatedUnix int64     `xorm:"created"`
	UpdatedUnix int64     `xorm:"updated"`
	LastUsedUnix int64    `xorm:"last_used"`
	ShowKey     bool      `xorm:"-"` // For template display
	HasRecentActivity bool `xorm:"-"` // For template display
	HasUsed     bool      `xorm:"-"` // For template display
}

// TableName returns the table name for PhantomKitKey
func (p *PhantomKitKey) TableName() string {
	return "phantomkit_keys"
}

// PhantomKitStats represents usage statistics
type PhantomKitStats struct {
	TotalKeys       int64
	TotalUploads    int64
	TotalExecutions int64
}

// PhantomKit render the PhantomKit API keys page
func PhantomKit(ctx *context.Context) {
	ctx.Data["Title"] = "PhantomKit API Keys"
	ctx.Data["PageIsSettingsPhantomKit"] = true

	loadPhantomKitData(ctx)

	ctx.HTML(http.StatusOK, tplSettingsPhantomKit)
}

// PhantomKitCreatePost handles creating a new PhantomKit API key
func PhantomKitCreatePost(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.PhantomKitKeyForm)
	ctx.Data["Title"] = "PhantomKit API Keys"
	ctx.Data["PageIsSettingsPhantomKit"] = true

	if ctx.HasError() {
		loadPhantomKitData(ctx)
		ctx.HTML(http.StatusOK, tplSettingsPhantomKit)
		return
	}

	// Generate a secure API key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		ctx.ServerError("GenerateKey", err)
		return
	}
	apiKey := "pkit_" + hex.EncodeToString(keyBytes)

	// Create the API key
	key := &PhantomKitKey{
		UserID:      ctx.Doer.ID,
		Name:        form.Name,
		Description: form.Description,
		Key:         apiKey,
	}

	if err := db.Insert(ctx, key); err != nil {
		ctx.ServerError("CreatePhantomKitKey", err)
		return
	}

	ctx.Flash.Success("PhantomKit API key created successfully!")
	ctx.Flash.Info("Your API key: " + apiKey)

	ctx.Redirect(setting.AppSubURL + "/user/settings/phantomkit")
}

// PhantomKitDeletePost handles deleting a PhantomKit API key
func PhantomKitDeletePost(ctx *context.Context) {
	keyID := ctx.FormInt64("id")
	
	// Delete the key (simplified for now)
	_, err := db.GetEngine(ctx).ID(keyID).Delete(&PhantomKitKey{})
	if err != nil {
		ctx.Flash.Error("Failed to delete API key: " + err.Error())
	} else {
		ctx.Flash.Success("PhantomKit API key deleted successfully!")
	}

	ctx.JSONRedirect(setting.AppSubURL + "/user/settings/phantomkit")
}

// PhantomKitToggleVisibility toggles key visibility
func PhantomKitToggleVisibility(ctx *context.Context) {
	// This would typically update a visibility flag in the database
	// For now, we'll just return success
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}

func loadPhantomKitData(ctx *context.Context) {
	// Load existing API keys
	var keys []PhantomKitKey
	err := db.GetEngine(ctx).Where("user_id = ?", ctx.Doer.ID).OrderBy("created_unix DESC").Find(&keys)
	if err != nil {
		// If table doesn't exist yet, just use empty slice
		if strings.Contains(err.Error(), "no such table") {
			keys = []PhantomKitKey{}
		} else {
			ctx.ServerError("ListPhantomKitKeys", err)
			return
		}
	}

	// Add display flags for template
	for i := range keys {
		keys[i].ShowKey = false // Keys are hidden by default
		keys[i].HasUsed = keys[i].LastUsedUnix > 0
		keys[i].HasRecentActivity = time.Now().Unix()-keys[i].LastUsedUnix < 86400 // 24 hours
	}

	ctx.Data["PhantomKitKeys"] = keys

	// Load usage statistics
	stats := &PhantomKitStats{
		TotalKeys: int64(len(keys)),
		// TODO: Implement actual usage tracking
		TotalUploads:    0,
		TotalExecutions: 0,
	}
	ctx.Data["Stats"] = stats
}

// ValidatePhantomKitKey validates a PhantomKit API key
func ValidatePhantomKitKey(ctx *context.Context, apiKey string) (*PhantomKitKey, error) {
	if len(apiKey) < 10 {
		return nil, db.ErrNotExist{}
	}

	key := &PhantomKitKey{}
	has, err := db.GetEngine(ctx).Where("key = ?", apiKey).Get(key)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, db.ErrNotExist{}
	}

	// Update last used timestamp
	key.LastUsedUnix = time.Now().Unix()
	if _, err := db.GetEngine(ctx).ID(key.ID).Update(key); err != nil {
		// Log error but don't fail the request
		ctx.ServerError("UpdateLastUsed", err)
	}

	return key, nil
}
