// back/internal/config/config.go
package config

import (
	"docx-converter-demo/internal/types"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func MustLoad() *types.Config {
	_ = godotenv.Load() // ไม่มีก็ไม่เป็นไร

	cfg := &types.Config{
		Port:        getEnv("PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "*"),
	}

	return cfg
}

func Addr(c *types.Config) string {
	return fmt.Sprintf(":%s", c.Port)
}

func Origins(c *types.Config) string {
	parts := strings.Split(c.FrontendURL, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return strings.Join(parts, ",")
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
