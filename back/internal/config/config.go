package config

import (
	"docx-converter-demo/internal/types"
	"docx-converter-demo/internal/utils"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

const (
	DefaultPort        = "2000"
	DefaultMaxUploadMB = 10
)

func MustLoad() *types.Config {
	_ = godotenv.Load() // Load .env file if exists

	cfg := &types.Config{
		Port:             utils.GetEnvString("PORT", DefaultPort),
		CORSAllowOrigins: utils.GetEnvString("CORS_ALLOW_ORIGINS", ""),
		MaxUploadMB:      utils.GetEnvInt("MAX_UPLOAD_MB", DefaultMaxUploadMB),
	}

	return cfg
}

func Addr(c *types.Config) string {
	return fmt.Sprintf(":%s", c.Port)
}

func Origins(c *types.Config) string {
	if c.CORSAllowOrigins == "" {
		return ""
	}
	parts := strings.Split(c.CORSAllowOrigins, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return strings.Join(parts, ",")
}

func MaxUploadBytes(c *types.Config) int64 {
	mb := c.MaxUploadMB
	if mb <= 0 {
		mb = DefaultMaxUploadMB
	}
	return int64(mb) * 1024 * 1024
}
