package utils

import (
// "fmt"
// "os"
// "strings"
// "time"
)

// IsValidEngine checks if the provided engine name is valid
func IsValidEngine(engine string) bool {
	validEngines := []string{"chatgpt", "gemini", "ollama"}
	for _, e := range validEngines {
		if e == engine {
			return true
		}
	}
	return false
}

// // TruncateString truncates a string to a specified length and adds an ellipsis
// func TruncateString(s string, maxLength int) string {
// 	if len(s) <= maxLength {
// 		return s
// 	}
// 	return s[:maxLength-3] + "..."
// }

// // GetEnvOrDefault retrieves an environment variable or returns a default value
// func GetEnvOrDefault(key, defaultValue string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return defaultValue
// }

// // FormatDuration formats a duration in a human-readable format
// func FormatDuration(d time.Duration) string {
// 	d = d.Round(time.Second)
// 	h := d / time.Hour
// 	d -= h * time.Hour
// 	m := d / time.Minute
// 	d -= m * time.Minute
// 	s := d / time.Second

// 	parts := []string{}
// 	if h > 0 {
// 		parts = append(parts, fmt.Sprintf("%dh", h))
// 	}
// 	if m > 0 {
// 		parts = append(parts, fmt.Sprintf("%dm", m))
// 	}
// 	if s > 0 || len(parts) == 0 {
// 		parts = append(parts, fmt.Sprintf("%ds", s))
// 	}
// 	return strings.Join(parts, " ")
// }

// // SanitizeFilename removes or replaces characters that are unsafe for filenames
// func SanitizeFilename(filename string) string {
// 	// Replace unsafe characters with underscores
// 	unsafe := []string{"/", "\\", "?", "%", "*", ":", "|", "\"", "<", ">"}
// 	for _, char := range unsafe {
// 		filename = strings.ReplaceAll(filename, char, "_")
// 	}
// 	return filename
// }
