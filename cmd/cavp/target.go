package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func run(target string) error {
	info, err := os.Stat(target)
	if err != nil {
		return err
	}

	switch {
	case info.IsDir():
		requests, err := requestFiles(target)
		if err != nil {
			return err
		}
		_, err = processRequests(requests)
		return err

	case info.Mode().IsRegular() && strings.EqualFold(filepath.Ext(target), ".zip"):
		output, err := processZip(target)
		if err != nil {
			return err
		}
		log.Printf("response directory: %s", output)
		return nil

	default:
		return fmt.Errorf("target must be a directory or ZIP file: %s", target)
	}
}

func requestFiles(root string) ([]string, error) {
	var requests []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ".req") {
			requests = append(requests, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(requests)
	if len(requests) == 0 {
		return nil, fmt.Errorf("no .req files found in %s", root)
	}
	return requests, nil
}

func processRequests(requests []string) ([]string, error) {
	responses := make([]string, 0, len(requests))
	for idx, path := range requests {
		log.Printf("[%d / %d] %s", idx+1, len(requests), path)
		err := process(path, filepath.Base(path))
		if errors.Is(err, errUnknownAlgorithm) {
			log.Printf("skip %s: %v", path, err)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		responses = append(responses, responsePath(path))
	}
	if len(responses) == 0 {
		return nil, errors.New("no supported CAVP request files found")
	}
	return responses, nil
}

func processZip(input string) (string, error) {
	output := strings.TrimSuffix(input, filepath.Ext(input))
	if output == "" || output == "." {
		return "", fmt.Errorf("cannot determine output directory from ZIP file: %s", input)
	}
	if _, err := os.Lstat(output); err == nil {
		return "", fmt.Errorf("output directory already exists: %s", output)
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	tempDir, err := os.MkdirTemp(filepath.Dir(input), ".cavp-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	if err := extractZip(input, tempDir); err != nil {
		return "", err
	}
	requests, err := requestFiles(tempDir)
	if err != nil {
		return "", err
	}
	_, err = processRequests(requests)
	if err != nil {
		return "", err
	}
	for _, request := range requests {
		if err := os.Remove(request); err != nil {
			return "", fmt.Errorf("remove request file %s: %w", request, err)
		}
	}

	if err := os.Rename(tempDir, output); err != nil {
		return "", err
	}
	return output, nil
}

func extractZip(input, destination string) error {
	r, err := zip.OpenReader(input)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, entry := range r.File {
		if entry.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("ZIP contains a symbolic link: %s", entry.Name)
		}
		name := filepath.Clean(filepath.FromSlash(entry.Name))
		if name == "." {
			continue
		}
		if filepath.IsAbs(name) || name == ".." || strings.HasPrefix(name, ".."+string(filepath.Separator)) {
			return fmt.Errorf("ZIP entry escapes destination: %s", entry.Name)
		}
		path := filepath.Join(destination, name)

		if entry.FileInfo().IsDir() {
			continue
		}
		if !strings.EqualFold(filepath.Ext(name), ".req") {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := extractZipFile(entry, path); err != nil {
			return err
		}
	}
	return nil
}

func extractZipFile(entry *zip.File, path string) error {
	src, err := entry.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(dst, src)
	closeErr := dst.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
