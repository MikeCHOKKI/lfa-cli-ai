package data

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed agents/*.md
//go:embed opencode.jsonc
//go:embed AGENTS.md
var Files embed.FS

func ReadAgent(name string) ([]byte, error) {
	return Files.ReadFile(filepath.Join("agents", name+".md"))
}

func ListAgents() ([]string, error) {
	entries, err := Files.ReadDir("agents")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func ListSkillsOnDisk(dataDir string) ([]string, error) {
	skillsDir := filepath.Join(dataDir, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		fullPath := filepath.Join(skillsDir, e.Name())
		info, err := os.Stat(fullPath)
		if err == nil && info.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func CopySkillFromDisk(dataDir, skillName, destDir string) error {
	src := filepath.Join(dataDir, "skills", skillName)
	dst := filepath.Join(destDir, skillName)
	return fs.WalkDir(os.DirFS(src), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, path)
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}
		content, err := os.ReadFile(filepath.Join(src, path))
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, content, 0644)
	})
}
