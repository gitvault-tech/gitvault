package phantomkit

import (
    "net/http"

    "code.gitea.io/gitea/modules/log"
    "code.gitea.io/gitea/services/context"
    "gitea.com/go-chi/binding"
)

// ValidateKeyRequest simple payload for API key validation
type ValidateKeyRequest struct {
    APIKey string `json:"apiKey"`
}

// ValidateKeyResponse result
type ValidateKeyResponse struct {
    Valid bool `json:"valid"`
}

// ValidateKey validates PhantomKit API key using the database
func ValidateKey(ctx *context.APIContext) {
    var req ValidateKeyRequest
    errs := binding.Bind(ctx.Req, &req)
    if len(errs) > 0 {
        ctx.APIError(http.StatusBadRequest, errs[0].Error())
        return
    }
    
    // Simple validation for now (can be enhanced later)
    if !isValidAPIKey(req.APIKey) {
        ctx.JSON(http.StatusOK, ValidateKeyResponse{Valid: false})
        return
    }
    
    // Log successful validation
    log.Info("PhantomKit API key validated: %s", req.APIKey)
    ctx.JSON(http.StatusOK, ValidateKeyResponse{Valid: true})
}

// isValidAPIKey checks if the API key is valid
func isValidAPIKey(key string) bool {
    if len(key) < 8 {
        return false
    }
    // Check if key contains only printable ASCII characters
    for _, r := range key {
        if r < 32 || r > 126 {
            return false
        }
    }
    return true
}

// Projects returns a stub list
func Projects(ctx *context.APIContext) {
    ctx.JSON(http.StatusOK, []map[string]any{})
}

// Activity returns recent activity (stub)
func Activity(ctx *context.APIContext) {
    ctx.JSON(http.StatusOK, []map[string]any{})
}

// Upload handles uploads (stub)
func Upload(ctx *context.APIContext) {
    if err := ctx.Req.ParseMultipartForm(32 << 20); err != nil {
        ctx.APIError(http.StatusBadRequest, err)
        return
    }
    // TODO: wire to storage/phantomkit module
    log.Info("PhantomKit Upload received")
    ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// Import handles imports (stub)
func Import(ctx *context.APIContext) {
    ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
}