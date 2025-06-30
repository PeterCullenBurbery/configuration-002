package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/PeterCullenBurbery/go_functions_002/yaml_functions"
	"github.com/PeterCullenBurbery/go_functions_002/system_management_functions"
)

type ProgramEntry struct {
	Name         string   `yaml:"name"`
	Alternatives []string `yaml:"alternatives"`
	WingetID     string   `yaml:"winget id,omitempty"`
}

// sanitize trims the Winget ID and all Alternatives
func (p *ProgramEntry) sanitize() {
	p.WingetID = strings.TrimSpace(p.WingetID)
	for i, alt := range p.Alternatives {
		p.Alternatives[i] = strings.TrimSpace(alt)
	}
}

type InstallYaml struct {
	Install map[string]map[string]ProgramEntry `yaml:"install"`
}

func main() {
	whatPath := flag.String("what", "", "Path to what-to-install.yaml (required)")
	installPath := flag.String("install", "", "Path to install.yaml (required)")
	logPath := flag.String("log", "", "Path to log file (required)")
	flag.Parse()

	if *whatPath == "" || *installPath == "" || *logPath == "" {
		fmt.Println("❌ --what, --install, and --log are required.")
		flag.Usage()
		os.Exit(1)
	}

	logFile, err := os.OpenFile(*logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	logFile.WriteString("\xEF\xBB\xBF")
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Load install.yaml
	var installData InstallYaml
	rawInstallData, err := os.ReadFile(*installPath)
	if err != nil {
		log.Fatalf("❌ Failed to read install.yaml: %v", err)
	}
	if err := yaml.Unmarshal(rawInstallData, &installData); err != nil {
		log.Fatalf("❌ Failed to parse install.yaml: %v", err)
	}

	// Build lookup maps
	altToCanonical := make(map[string]string)
	canonicalToMeta := make(map[string]ProgramEntry)
	canonicalToCategory := make(map[string]string)

	for category, programs := range installData.Install {
		for canonical, meta := range programs {
			canonicalTrimmed := strings.TrimSpace(canonical)
			meta.sanitize()
			canonicalToMeta[canonicalTrimmed] = meta
			canonicalToCategory[canonicalTrimmed] = category
			altToCanonical[strings.ToLower(canonicalTrimmed)] = canonicalTrimmed
			for _, alt := range meta.Alternatives {
				altToCanonical[strings.ToLower(alt)] = canonicalTrimmed
			}
		}
	}

	// Load what-to-install.yaml
	whatData := make(map[string]interface{})
	rawWhatData, err := os.ReadFile(*whatPath)
	if err != nil {
		log.Fatalf("❌ Failed to read what-to-install.yaml: %v", err)
	}
	if err := yaml.Unmarshal(rawWhatData, &whatData); err != nil {
		log.Fatalf("❌ Failed to parse what-to-install.yaml: %v", err)
	}

	installSection := yaml_functions.GetCaseInsensitiveMap(whatData, "install")
	if installSection == nil {
		log.Fatal("❌ Missing 'install' section in what-to-install.yaml.")
	}
	requested := yaml_functions.GetCaseInsensitiveList(installSection, "programs to install")

	for _, req := range requested {
		lookup := strings.ToLower(strings.TrimSpace(req))
		canonical, ok := altToCanonical[lookup]
		if !ok {
			log.Printf("❌ Unsupported program: %s (skipped)", req)
			continue
		}

		category := canonicalToCategory[canonical]
		meta := canonicalToMeta[canonical]

		if category != "winget" {
			log.Printf("⏩ Skipping non-winget program: %s (category: %s)", canonical, category)
			continue
		}

		if meta.WingetID == "" {
			log.Printf("⚠️ Missing Winget ID for %s", canonical)
			continue
		}

		log.Printf("🚀 Installing winget program: %s (%s)", canonical, meta.WingetID)
		if err := system_management_functions.Winget_install(canonical, meta.WingetID); err != nil {
			log.Printf("❌ Failed to install %s via winget: %v", canonical, err)
		} else {
			log.Printf("✅ Installed %s via winget.", canonical)
		}
	}

	log.Println("🎉 Winget-only installation process finished.")
}
