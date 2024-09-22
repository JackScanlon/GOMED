package templates

import (
	"bufio"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"maps"
	"path/filepath"
	"regexp"
	"snomed/src/pg"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"golang.org/x/sync/errgroup"
)

const (
	tmplFileFormat string = "content/*.sql"
	tmplTokenStart string = "template"
	tmplTokenEnd   string = "endtemplate"
)

var (
	//go:embed content/*.sql
	sources embed.FS

	tmplInstance *TemplateContainer
	tmplOnce     *sync.Mutex    = &sync.Mutex{}
	tmplTokenRe  *regexp.Regexp = regexp.MustCompile(`^--!\[(\w+)\]`)
	tmplParamRe  *regexp.Regexp = regexp.MustCompile(`(\w+):("[^"\\]*(?:\\.[^"\\]*)*")`)
)

type TemplateItem struct {
	Name     string
	Filename string
	Content  *template.Template
	Params   map[string]any
}

type TemplateContainer struct {
	Children map[string]map[string]*TemplateItem
}

func GetContainer() *TemplateContainer {
	container, err := TryGetContainer()
	if err != nil {
		panic(err)
	}

	return container
}

func TryGetContainer() (*TemplateContainer, error) {
	if tmplInstance == nil {
		return nil, fmt.Errorf("template container is not initialised")
	}

	return tmplInstance, nil
}

func InitContainer(ctx context.Context) (*TemplateContainer, error) {
	err := buildTemplates(ctx)
	if err != nil {
		return nil, err
	}

	return tmplInstance, nil
}

func (t *TemplateContainer) Source(name string, options ...TmplOption) *Template {
	tmpl, err := t.TrySource(name, options...)
	if err != nil {
		panic(err)
	}

	return tmpl
}

func (t *TemplateContainer) TrySource(name string, options ...TmplOption) (*Template, error) {
	src := strings.Split(name, ":")
	if len(src) < 2 {
		return nil, fmt.Errorf("invalid name, expected '<source>:<name>' but got %s", name)
	}

	children, ok := t.Children[src[0]]
	if !ok {
		return nil, fmt.Errorf("unknown template with source '%s' from: '%s'", src[0], name)
	}

	templateItem, ok := children[src[1]]
	if !ok {
		return nil, fmt.Errorf("no known template with name '%s' from: '%s'", src[1], name)
	}

	tmpl := &Template{
		Name:    name,
		Data:    maps.Clone(templateItem.Params),
		Content: templateItem.Content,
		HasOpts: false,
	}

	for _, opt := range options {
		if err := opt(tmpl); err != nil {
			return nil, err
		}
	}

	if !tmpl.HasOpts {
		tmpl.PgOpts = pg.PgOptions{Ctx: context.TODO()}
	}

	return tmpl, nil
}

func buildTemplates(ctx context.Context) error {
	tmplOnce.Lock()
	defer tmplOnce.Unlock()

	if tmplInstance != nil {
		return nil
	}

	tmplInstance = &TemplateContainer{
		Children: map[string]map[string]*TemplateItem{},
	}

	entries, err := fs.Glob(sources, tmplFileFormat)
	if err != nil {
		return err
	}

	if len(entries) < 1 {
		return nil
	}

	group, _ := errgroup.WithContext(ctx)
	results := make(chan *TemplateItem)

	for _, entry := range entries {
		group.Go(func() error {
			file, err := sources.Open(entry)
			if err != nil {
				return err
			}

			filename := strings.TrimSuffix(filepath.Base(entry), filepath.Ext(entry))
			return buildQueries(results, filename, bufio.NewScanner(file))
		})
	}

	go func() {
		err = group.Wait()
		close(results)
	}()

	for template := range results {
		if tmplInstance.Children[template.Filename] == nil {
			tmplInstance.Children[template.Filename] = map[string]*TemplateItem{}
		}

		tmplInstance.Children[template.Filename][template.Name] = template
	}

	return group.Wait()
}

func buildQueries(resultChan chan *TemplateItem, filename string, scanner *bufio.Scanner) error {
	var (
		tk *TemplateItem
		sb strings.Builder

		lcount uint32 = 0
		tstart uint32 = 0
	)

	for scanner.Scan() {
		line := scanner.Text()
		lcount++

		indices := tmplTokenRe.FindStringSubmatchIndex(line)
		if len(indices) > 0 {
			token := strings.ToLower(line[indices[2]:indices[3]])
			switch token {
			case tmplTokenStart:
				if tk != nil {
					return fmt.Errorf(
						("new template block started on line %d;" +
							"failed to close the template started at line %d with name '%s' in file '%s'"),
						lcount, tstart, tk.Name, tk.Filename,
					)
				}

				var name string
				params := map[string]any{}

				matches := tmplParamRe.FindAllStringSubmatch(line[indices[3]:], -1)
				for _, match := range matches {
					key := strings.ToLower(match[1])
					value, err := strconv.Unquote(match[2])
					if err != nil {
						return err
					}

					if key == "name" {
						name = value
					} else {
						params[key] = value
					}
				}

				if len(name) < 1 {
					return fmt.Errorf("expected name for template tag started at line %d in file '%s'", lcount, filename)
				}

				tk = &TemplateItem{
					Name:     name,
					Filename: filename,
					Params:   params,
				}
				tstart = lcount

				continue
			case tmplTokenEnd:
				content := sb.String()
				if len(strings.TrimSpace(content)) < 1 {
					continue
				}

				tmpl, err := template.New(
					fmt.Sprintf("%s:%s", filename, tk.Name),
				).
					Parse(content)

				if err != nil {
					return err
				}

				tk.Content = tmpl
				resultChan <- tk

				sb.Reset()
				tk = nil
				continue
			default:
				break
			}
		}

		if tk != nil {
			if _, err := sb.WriteString(line + "\n"); err != nil {
				return err
			}
		}
	}

	return nil
}
