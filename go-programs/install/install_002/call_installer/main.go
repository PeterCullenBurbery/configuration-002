package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/PeterCullenBurbery/go_functions_002/date_time_functions"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("❌ Usage: call_installer.exe <base-directory>")
		os.Exit(1)
	}

	baseDir := os.Args[1]

	// Compute paths
	canonicalDir := filepath.Join(baseDir, "canonical-file-structure")
	whatPath := filepath.Join(canonicalDir, "what-to-install.yaml")
	installPath := filepath.Join(canonicalDir, "the-following-lines-are-available-to-install.yaml")
	logsDir := filepath.Join(canonicalDir, "logs")

	// Create logs directory if needed
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create logs directory: %v", err)
	}

	// Generate timestamp for log file
	timestamp, err := date_time_functions.Date_time_stamp()
	if err != nil {
		log.Fatalf("❌ Could not generate timestamp: %v", err)
	}
	safeTimestamp := date_time_functions.Safe_time_stamp(timestamp, 1)
	logPath := filepath.Join(logsDir, safeTimestamp+".log")

	// Locate winget_install.exe
	wingetExe := filepath.Join(canonicalDir, "go-programs", "winget_install", "winget_install.exe")
	if _, err := os.Stat(wingetExe); os.IsNotExist(err) {
		log.Fatalf("❌ winget_install.exe not found at: %s", wingetExe)
	}

	// Prepare and run the command
	cmd := exec.Command(wingetExe, "--what", whatPath, "--install", installPath, "--log", logPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("🚀 Running winget_install.exe\n")
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ winget_install.exe failed: %v", err)
	}

	fmt.Println("✅ Installation complete.")
}
