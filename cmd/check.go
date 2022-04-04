package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func findDirsWithFiles() []string {
	var toReturn []string

	var searchPath string
	if strings.HasPrefix(tfDir, "~/") {
		usr, _ := user.Current()
		homeDir := usr.HomeDir

		searchPath = filepath.Join(homeDir, tfDir[2:])
	} else {
		searchPath = tfDir
	}

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			return nil
		}

		if !info.IsDir() && info.Name() == ".terraform.lock.hcl" {
			toReturn = append(toReturn, filepath.Dir(path))
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return toReturn
}

func check() {
	dirSet := mapset.NewSet()

	if findLockFiles {
		for _, dir := range findDirsWithFiles() {
			dirSet.Add(dir)
		}
	} else {
		if !lockFileExists(tfDir) {
			fmt.Println("No .terraform.lock.hcl found. Exiting")
			os.Exit(1)
		}
		dirSet.Add(tfDir)
	}

	execPath, err := exec.LookPath("terraform")
	if err != nil {
		fmt.Println("Terraform executable not found on path, install terraform")
		os.Exit(1)
	}

	updatesAvailable := false

	dirIterator := dirSet.Iterator()
	for dir := range dirIterator.C {
		fmt.Println("Found lockfile in: " + dir.(string))
		tf, err := tfexec.NewTerraform(dir.(string), execPath)
		if err != nil {
			fmt.Printf("Error accessing Terraform: %s\n", err)
			os.Exit(1)
		}

		_, providerVersions, err := tf.Version(context.Background(), true)
		if err != nil {
			fmt.Printf("Error running terraform version: %s\n", err)
			os.Exit(1)
		}

		for provider, version := range providerVersions {
			updates := checkVersion(provider, version)
			if updates {
				updatesAvailable = true
			}
		}
	}

	if errorOnUpdate {
		if updatesAvailable {
			os.Exit(1)
		}
	}

	os.Exit(0)
}

// RegistryResponse represents the JSON return by the HC Registry as a struct
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
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Couldn't get read response body: %s\n", err)
	}

	var response RegistryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Could not unmarshal response JSON: %s\n", err)
	}

	remoteVersion, err := version.NewVersion(response.Version)
	if err != nil {
		fmt.Printf("Could not parse remote version: %s\n", err)
	}

	if localVersion.LessThan(remoteVersion) {
		fmt.Println("Update of", color.HiYellowString(providerName), "available", color.RedString(localVersion.String()), "<", color.BlueString(remoteVersion.String()))
		return true
	}

	return false
}
