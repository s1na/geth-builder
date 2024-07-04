package builder

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/s1na/geth-builder/config"
	"github.com/s1na/geth-builder/transform"
)

type Builder struct {
	config *config.Config
	arch   *string
}

func NewBuilder(config *config.Config, arch *string) *Builder {
	return &Builder{config, arch}
}

func (b *Builder) Build() (string, error) {
	gethDir := "./go-ethereum"
	if err := b.prepareSource(gethDir); err != nil {
		return "", err
	}

	return gethDir, b.build(gethDir, "./cmd/geth")

}

func (b *Builder) build(gethDir, pkg string) error {
	args := []string{"run", "build/ci.go", "install"}
	if b.arch != nil {
		args = append(args, "--arch", *b.arch)
	}
	if pkg != "" {
		args = append(args, pkg)
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = gethDir
	if b.config.Verbose() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

func (b *Builder) Archive(typ string) error {
	gethDir := "./go-ethereum"
	if err := b.prepareSource(gethDir); err != nil {
		return fmt.Errorf("failed to fetch and transform source: %v", err)
	}
	// Geth archive builder needs all the tooling binaries.
	if err := b.build(gethDir, ""); err != nil {
		return fmt.Errorf("failed to build binaries: %v", err)
	}

	args := []string{"run", "build/ci.go", "archive"}
	if b.arch != nil {
		args = append(args, "--arch", *b.arch)
	}
	args = append(args, "--type", typ)
	args = append(args, "./cmd/geth")
	cmd := exec.Command("go", args...)
	cmd.Dir = gethDir
	// Hack to extract archive paths from the output.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if b.config.Verbose() {
			fmt.Print(stdout.String())
			fmt.Fprint(os.Stderr, stderr.String())
		}
		return fmt.Errorf("geth archive cmd failed: %v", err)
	}
	// Still need to log to user in case verbose flag is set.
	if b.config.Verbose() {
		fmt.Print(stdout.String())
		fmt.Fprint(os.Stderr, stderr.String())
	}
	output, err := b.config.AbsoluteOutputDir()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`geth-.*\.(zip|tar.gz)`)
	matches := re.FindAllString(stdout.String(), -1)
	if len(matches) == 0 {
		return fmt.Errorf("no archive files found")
	}
	fmt.Printf("matches: %v\n", matches)
	for _, match := range matches {
		if err := CopyFile(filepath.Join(gethDir, match), filepath.Join(output, match)); err != nil {
			return fmt.Errorf("failed to copy archive: %v", err)
		}
	}

	return nil
}

func (b *Builder) prepareSource(gethDir string) error {
	// Clone the specified Geth repository and branch
	err := CloneRepo(b.config.GethRepo, b.config.GethBranch, gethDir, b.config.Verbose())
	if err != nil {
		return err
	}
	// Use the AbsolutePath method to get the absolute path
	absTracerPath, err := b.config.AbsolutePath()
	if err != nil {
		return err
	}

	dest := filepath.Join(gethDir, "eth", "tracers", "native")
	// Copy the local tracer to the Geth tracers directory
	err = copyLocalTracer(absTracerPath, dest)
	if err != nil {
		return err
	}
	newTracerPath := filepath.Join(dest, filepath.Base(absTracerPath))
	// Remove go.mod and go.sum files from tracing package if available.
	if err := os.RemoveAll(filepath.Join(newTracerPath, "go.mod")); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(newTracerPath, "go.sum")); err != nil {
		return err
	}

	gethMainPath := filepath.Join(gethDir, "cmd", "geth", "main.go")
	pkgName := filepath.Base(absTracerPath)
	importPath := filepath.Join("github.com/ethereum/go-ethereum", "eth", "tracers", "native", pkgName)
	err = transform.AddImportToFile(gethMainPath, importPath)
	if err != nil {
		return err
	}
	return nil
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
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return err
			}
			if err := copyLocalTracer(srcPath, destDir); err != nil {
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
