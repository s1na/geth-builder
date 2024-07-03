package builder

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/s1na/geth-builder/config"
	"github.com/s1na/geth-builder/transform"
)

func BuildGeth(config *config.Config) (string, error) {
	// Clone the specified Geth repository and branch
	gethDir := "./go-ethereum"
	err := CloneRepo(config.GethRepo, config.GethBranch, gethDir, config.Verbose())
	if err != nil {
		return "", err
	}
	// Use the AbsolutePath method to get the absolute path
	absTracerPath, err := config.AbsolutePath()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(gethDir, "eth", "tracers", "native")
	// Copy the local tracer to the Geth tracers directory
	err = copyLocalTracer(absTracerPath, dest)
	if err != nil {
		return "", err
	}
	newTracerPath := filepath.Join(dest, filepath.Base(absTracerPath))
	// Remove go.mod and go.sum files from tracing package if available.
	if err := os.RemoveAll(filepath.Join(newTracerPath, "go.mod")); err != nil {
		return "", err
	}
	if err := os.RemoveAll(filepath.Join(newTracerPath, "go.sum")); err != nil {
		return "", err
	}

	gethMainPath := filepath.Join(gethDir, "cmd", "geth", "main.go")
	pkgName := filepath.Base(absTracerPath)
	importPath := filepath.Join("github.com/ethereum/go-ethereum", "eth", "tracers", "native", pkgName)
	err = transform.AddImportToFile(gethMainPath, importPath)
	if err != nil {
		log.Fatalf("Error modifying main.go: %v\n", err)
	}

	cmd := exec.Command("make", "geth")
	cmd.Dir = gethDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GOBIN="+filepath.Join(config.OutputDir))
	return gethDir, cmd.Run()
}

func CloneRepo(repoURL, branch, destDir string, verbose bool) error {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "--depth", "1", "-b", branch, repoURL, destDir)
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		return cmd.Run()
	}
	log.Printf("Directory %s already exists, skipping clone.\n", destDir)
	return nil
}

func copyLocalTracer(srcDir, destDir string) error {
	// Get the base name of the source directory
	base := filepath.Base(srcDir)

	// Create a new directory in the destination with the same name as the source directory
	destDir = filepath.Join(destDir, base)
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())
		if entry.IsDir() {
			err = os.MkdirAll(destPath, entry.Type())
			if err != nil {
				return err
			}
			err = copyLocalTracer(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	return os.Chmod(dest, srcInfo.Mode())
}
