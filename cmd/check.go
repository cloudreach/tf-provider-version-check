package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func check() {
	if !lockFileExists(tfDir) {
		fmt.Println("No .terraform.lock.hcl found. Exiting")
		os.Exit(1)
	}

	execPath, err := exec.LookPath("terraform")
	if err != nil {
		fmt.Println("Terraform executable not found on path, install terraform")
		os.Exit(1)
	}

	tf, err := tfexec.NewTerraform(tfDir, execPath)
	if err != nil {
		fmt.Printf("Error accessing Terraform: %s\n", err)
		os.Exit(1)
	}

	_, providerVersions, err := tf.Version(context.Background(), true)
	if err != nil {
		fmt.Printf("Error running terraform version: %s\n", err)
		os.Exit(1)
	}

	updatesAvailable := false
	for provider, version := range providerVersions {
		updates := checkVersion(provider, version)
		if updates {
			updatesAvailable = true
		}
	}

	if errorOnUpdate {
		if updatesAvailable {
			os.Exit(1)
		}
	}

	os.Exit(0)
}

type RegistryResponse struct {
	ID          string    `json:"id"`
	Owner       string    `json:"owner"`
	Namespace   string    `json:"namespace"`
	Name        string    `json:"name"`
	Alias       string    `json:"alias"`
	Version     string    `json:"version"`
	Tag         string    `json:"tag"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
	Downloads   int       `json:"downloads"`
	Tier        string    `json:"tier"`
	LogoURL     string    `json:"logo_url"`
	Versions    []string  `json:"versions"`
	Docs        []struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Path        string `json:"path"`
		Slug        string `json:"slug"`
		Category    string `json:"category"`
		Subcategory string `json:"subcategory"`
	} `json:"docs"`
}

func lockFileExists(path string) bool {
	lockFilePath := path + "/.terraform.lock.hcl"
	_, err := os.Stat(lockFilePath)
	return !errors.Is(err, os.ErrNotExist)
}

func checkVersion(provider string, localVersion *version.Version) bool {
	providerName := strings.SplitN(provider, "/", 2)[1]

	resp, err := http.Get("https://registry.terraform.io/v1/providers/" + providerName)
	if err != nil {
		fmt.Printf("Couldn't get provider details from HC: %s\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Couldn't get read response body: %s\n", err)
		os.Exit(1)
	}

	var response RegistryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Could not unmarshal response JSON: %s\n", err)
		os.Exit(1)
	}

	remoteVersion, err := version.NewVersion(response.Version)
	if err != nil {
		fmt.Printf("Could not parse remote version: %s\n", err)
		os.Exit(1)
	}

	if localVersion.LessThan(remoteVersion) {
		fmt.Println("Update of", providerName, "available", localVersion, "<", remoteVersion)
		return true
	}

	return false
}
