package main

import (
	"log"
	"os"

	"github.com/PeterCullenBurbery/go_functions_002/system_management_functions"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("🔧 Attempting to install Chocolatey...")

	if err := system_management_functions.Install_choco(); err != nil {
		log.Fatalf("❌ Chocolatey installation failed: %v", err)
	}

	log.Println("🎉 Chocolatey is installed and ready to use.")
}
