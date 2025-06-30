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
	ChocoID      string   `yaml:"choco id,omitempty"`
}

func (p *ProgramEntry) sanitize() {
	p.WingetID = strings.TrimSpace(p.WingetID)
	p.ChocoID = strings.TrimSpace(p.ChocoID)
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
	logFile.WriteString("\xEF\xBB\xBF") // Write UTF-8 BOM
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	log.Println("📦 Ensuring Chocolatey is installed...")
	if err := system_management_functions.Install_choco(); err != nil {
		log.Fatalf("❌ Failed to install Chocolatey: %v", err)
	}

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

	for category, programs := range installData.Install {
		if strings.ToLower(category) != "choco" {
			continue
		}
		for canonical, meta := range programs {
			canonicalTrimmed := strings.TrimSpace(canonical)
			meta.sanitize()
			canonicalToMeta[canonicalTrimmed] = meta
			altToCanonical[strings.ToLower(canonicalTrimmed)] = canonicalTrimmed
			for _, alt := range meta.Alternatives {
				altToCanonical[strings.ToLower(alt)] = canonicalTrimmed
			}
		}
	}

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

	// Process each requested program
	for _, req := range requested {
		lookup := strings.ToLower(strings.TrimSpace(req))
		canonical, ok := altToCanonical[lookup]
		if !ok {
			log.Printf("❌ Unsupported program for choco: %s (skipped)", req)
			continue
		}
		meta := canonicalToMeta[canonical]
		if meta.ChocoID == "" {
			log.Printf("⚠️  Missing Choco ID for %s", canonical)
			continue
		}
		log.Printf("🔧 Installing %s via Chocolatey (ID: %s)", canonical, meta.ChocoID)
		if err := system_management_functions.Choco_install(meta.ChocoID); err != nil {
			log.Printf("❌ Chocolatey install failed for %s: %v", canonical, err)
		} else {
			log.Printf("✅ Successfully installed %s via Chocolatey.", canonical)
		}
	}

	log.Println("🎉 Chocolatey installation tasks finished.")
}