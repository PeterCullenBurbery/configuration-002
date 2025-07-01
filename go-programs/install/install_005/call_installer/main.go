package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/PeterCullenBurbery/go_functions_002/date_time_functions"
)

func runExecutable(path string, args ...string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("🚀 Running: %s %v\n", path, args)
	return cmd.Run()
}

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

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create logs directory: %v", err)
	}

	timestamp, err := date_time_functions.Date_time_stamp()
	if err != nil {
		log.Fatalf("❌ Could not generate timestamp: %v", err)
	}
	safeTimestamp := date_time_functions.Safe_time_stamp(timestamp, 1)
	logPath := filepath.Join(logsDir, safeTimestamp+".log")

	// Prepare executables
	chocoInstaller := filepath.Join(canonicalDir, "go-programs", "install_Choco", "install_Choco.exe")
	wingetInstaller := filepath.Join(canonicalDir, "go-programs", "winget_install", "winget_install.exe")
	chocoPackages := filepath.Join(canonicalDir, "go-programs", "Choco_install", "Choco_install.exe")

	// Run install_Choco.exe
	if _, err := os.Stat(chocoInstaller); os.IsNotExist(err) {
		log.Fatalf("❌ install_Choco.exe not found at: %s", chocoInstaller)
	}
	if err := runExecutable(chocoInstaller); err != nil {
		log.Fatalf("❌ install_Choco.exe failed: %v", err)
	}

	// Run winget_install.exe
	if _, err := os.Stat(wingetInstaller); os.IsNotExist(err) {
		log.Fatalf("❌ winget_install.exe not found at: %s", wingetInstaller)
	}
	if err := runExecutable(wingetInstaller, "--what", whatPath, "--install", installPath, "--log", logPath); err != nil {
		log.Fatalf("❌ winget_install.exe failed: %v", err)
	}

	// Run Choco_install.exe
	if _, err := os.Stat(chocoPackages); os.IsNotExist(err) {
		log.Fatalf("❌ Choco_install.exe not found at: %s", chocoPackages)
	}
	if err := runExecutable(chocoPackages, "--what", whatPath, "--install", installPath, "--log", logPath); err != nil {
		log.Fatalf("❌ Choco_install.exe failed: %v", err)
	}

	fmt.Println("✅ All installation steps completed.")
}
