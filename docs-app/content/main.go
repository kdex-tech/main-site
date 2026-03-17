package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"sigs.k8s.io/yaml"
)

type Document struct {
	Children    []*Document    `json:"children,omitempty"`
	FrontMatter map[string]any `json:"frontMatter,omitempty"`
	Name        string         `json:"-"`
	Path        string         `json:"path,omitempty"`
	Title       string         `json:"title"`
	Weight      int            `json:"weight"`
}

func extractAddress(content []byte, addr string) []byte {
	addr = strings.TrimSpace(addr)
	var remainder bytes.Buffer

	found := false
	for line := range strings.Lines(string(content)) {
		if strings.Contains(line, addr) {
			found = true
		}
		if found {
			remainder.WriteString(line)
		}
	}
	if found {
		return remainder.Bytes()
	}
	return fmt.Appendf(nil, "Error: Could not find address %s in content", addr)
}

func linkConversionHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if link, ok := node.(*ast.Link); ok && entering {
		dest := string(link.Destination)
		if before, after, ok0 := strings.Cut(dest, ".md"); ok0 {
			link.Destination = []byte(before + ".html" + after)
		}
	}
	return ast.GoToNext, false
}

func main() {
	manifest, err := generateManifest("src", "out")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating manifest: %v\n", err)
		os.Exit(1)
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling manifest: %v\n", err)
		os.Exit(1)
	}

	manifestPath := "out/manifest.json"
	os.MkdirAll("out", 0755)
	err = os.WriteFile(manifestPath, manifestData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", manifestPath, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s\n", manifestPath)
}

func generateManifest(srcRoot string, outRoot string) (map[string]*Document, error) {
	manifest := make(map[string]*Document)
	fmt.Printf("Starting renderer in: %s\n", srcRoot)

	err := filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == srcRoot {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(srcRoot, path)
		parts := strings.Split(relPath, string(os.PathSeparator))
		lang := parts[0]

		if _, ok := manifest[lang]; !ok {
			manifest[lang] = &Document{
				Title: lang,
			}
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		fileName := filepath.Base(path)
		dir := filepath.Dir(path)
		targetDir := strings.TrimPrefix(dir, srcRoot)
		targetDir = outRoot + targetDir
		targetPath := targetDir + "/" + strings.TrimSuffix(fileName, ".md") + ".html"
		fmt.Printf("Processing: %s -> %s\n", path, targetPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		frontMatter, remainder := readFrontMatter(path, string(content))
		htmlContent := mdToHTML([]byte(remainder))

		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(targetPath), err)
		}

		err = os.WriteFile(targetPath, htmlContent, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		manifestPath := strings.TrimPrefix(targetPath, outRoot+"/")
		_, manifestPath, _ = strings.Cut(manifestPath, "/")

		title := fileName
		weight := 0
		if t, ok := frontMatter["title"].(string); ok {
			title = t
		}
		if w, ok := frontMatter["weight"].(float64); ok {
			weight = int(w)
		} else if w, ok := frontMatter["weight"].(int); ok {
			weight = w
		}

		curr := manifest[lang]
		// Traverse the tree to find or create the parent node for the current file
		for i := 1; i < len(parts)-1; i++ {
			part := parts[i]
			found := false
			for _, child := range curr.Children {
				if child.Name == part {
					curr = child
					found = true
					break
				}
			}
			if !found {
				name, w := parseName(part)
				newNode := &Document{
					Name:   part,
					Title:  name,
					Weight: w,
				}
				curr.Children = append(curr.Children, newNode)
				curr = newNode
			}
		}

		if fileName == "_index.md" {
			// Update the current node (directory) with properties from _index.md
			curr.FrontMatter = frontMatter
			curr.Path = manifestPath
			curr.Title = title
			curr.Weight = weight
		} else {
			// Add regular page as a child
			doc := &Document{
				FrontMatter: frontMatter,
				Path:        manifestPath,
				Title:       title,
				Weight:      weight,
			}
			curr.Children = append(curr.Children, doc)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	sortTree(manifest)
	return manifest, nil
}

func mdToHTML(md []byte) []byte {
	mdStr := strings.ReplaceAll(string(md), "<details", "<kdex-details")
	mdStr = strings.ReplaceAll(mdStr, "</details>", "</kdex-details>")
	mdStr = strings.ReplaceAll(mdStr, "<summary", "<kdex-summary")
	mdStr = strings.ReplaceAll(mdStr, "</summary>", "</kdex-summary>")

	extensions := parser.CommonExtensions |
		parser.AutoHeadingIDs |
		parser.FencedCode |
		parser.NoEmptyLineBeforeBlock |
		parser.Attributes |
		parser.Includes

	p := parser.NewWithExtensions(extensions)
	p.Opts.ReadIncludeFn = readIncludeFn
	doc := p.Parse([]byte(mdStr))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: linkConversionHook,
	}
	renderer := html.NewRenderer(opts)

	output := markdown.Render(doc, renderer)
	outStr := string(output)

	// Clean up paragraph wrappers inserted by gomarkdown treating custom tags as inline text
	outStr = regexp.MustCompile(`(?s)<p>\s*(<kdex-details[^>]*>)\s*(<kdex-summary[^>]*>.*?</kdex-summary>)\s*</p>`).ReplaceAllString(outStr, "$1\n$2")
	outStr = regexp.MustCompile(`(?s)<p>\s*(</kdex-details>)\s*</p>`).ReplaceAllString(outStr, "$1")
	outStr = regexp.MustCompile(`(?s)<p>\s*(<kdex-details[^>]*>)\s*</p>`).ReplaceAllString(outStr, "$1")
	outStr = regexp.MustCompile(`(?s)<p>\s*(<kdex-summary[^>]*>.*?</kdex-summary>)\s*</p>`).ReplaceAllString(outStr, "$1")

	// Revert tags
	outStr = strings.ReplaceAll(outStr, "kdex-details", "details")
	outStr = strings.ReplaceAll(outStr, "kdex-summary", "summary")

	return []byte(outStr)
}

func parseName(name string) (string, int) {
	if before, after, ok := strings.Cut(name, "_"); ok {
		weight, err := strconv.Atoi(before)
		if err == nil {
			return after, weight
		}
	}
	return name, 0
}

func readFrontMatter(path string, content string) (frontMatter map[string]any, remainingContent string) {
	lang := "en"
	path = strings.TrimPrefix(path, "src/")
	if before, _, ok := strings.Cut(path, "/"); ok {
		lang = before
	}

	frontMatter = map[string]any{
		"lang": lang,
	}

	if strings.HasPrefix(content, "---") {
		after := strings.TrimPrefix(content, "---")
		idx := strings.Index(after, "---")
		if idx != -1 {
			block := after[:idx]
			err := yaml.Unmarshal([]byte(block), &frontMatter)
			if err != nil {
				fmt.Printf("error unmarshaling front matter: %v\n", err)
				return frontMatter, content
			}
			return frontMatter, after[idx+3:]
		}
	}

	for l := range strings.Lines(content) {
		if strings.HasPrefix(l, "##") {
			frontMatter["title"] = strings.TrimSpace(strings.TrimPrefix(l, "##"))
		} else if strings.HasPrefix(l, "# ") {
			frontMatter["title"] = strings.TrimSpace(strings.TrimPrefix(l, "# "))
		}
		break
	}

	return frontMatter, content
}

func readIncludeFn(parentPath string, includePath string, address []byte) []byte {
	var allowedRemoteRoots = []string{
		"https://raw.githubusercontent.com/kdex-tech/",
	}

	if strings.HasPrefix(includePath, "https://") {
		allowed := false
		for _, root := range allowedRemoteRoots {
			if strings.HasPrefix(includePath, root) {
				allowed = true
				break
			}
		}

		if !allowed {
			return fmt.Appendf(nil, "Error: Remote include path %s is not in the allowed list", includePath)
		}

		resp, err := http.Get(includePath)
		if err != nil {
			return fmt.Appendf(nil, "Error: Failed to fetch remote include %s: %v", includePath, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Appendf(nil, "Error: Failed to fetch remote include %s: status %d", includePath, resp.StatusCode)
		}

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Appendf(nil, "Error: Failed to read remote include response %s: %v", includePath, err)
		}

		if len(address) > 0 {
			return extractAddress(content, string(address))
		}

		return content
	}

	cwd, _ := os.Getwd()
	fullPath := includePath
	if parentPath != "" {
		dir := filepath.Dir(parentPath)
		fullPath = filepath.Join(dir, includePath)
	}
	fullPath = filepath.Join(cwd, fullPath)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return []byte("Error: Could not include file " + includePath + ", derived path was " + fullPath)
	}

	if len(address) > 0 {
		return extractAddress(content, string(address))
	}

	return content
}

func sortTree(manifest map[string]*Document) {
	for _, root := range manifest {
		sortNode(root)
	}
}

func sortNode(node *Document) {
	if len(node.Children) == 0 {
		return
	}
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].Weight != node.Children[j].Weight {
			return node.Children[i].Weight < node.Children[j].Weight
		}
		return node.Children[i].Title < node.Children[j].Title
	})
	for _, child := range node.Children {
		sortNode(child)
	}
}
