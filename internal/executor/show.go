package executor

import (
	"bytes"
	"context"
	"github.com/snowmerak/jetti/v2/internal/executor/check"
	"github.com/snowmerak/jetti/v2/lib/parser"
	"log"
	"os"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"path/filepath"
)

func ShowImports(root string) error {
	moduleName, err := check.GetModuleName(root)
	if err != nil {
		return err
	}

	links := []*check.DependencyLink(nil)

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		pkg, err := parser.ParseFile(path)
		if err != nil {
			return err
		}

		list, err := check.GetImports(root, moduleName, path, pkg)
		if err != nil {
			return err
		}

		links = append(links, list)

		return nil
	}); err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("direction: right\n")
	for _, link := range links {
		for _, from := range link.From {
			buffer.WriteString(from)
			buffer.WriteString(" --> \"")
			buffer.WriteString(link.To)
			buffer.WriteString("\"\n")
		}
	}

	ruler, _ := textmeasure.NewRuler()
	defaultLayout := func(ctx context.Context, g *d2graph.Graph) error {
		return d2dagrelayout.Layout(ctx, g, nil)
	}

	diagram, _, _ := d2lib.Compile(context.Background(), buffer.String(), &d2lib.CompileOptions{
		Layout: defaultLayout,
		Ruler:  ruler,
	})

	out, _ := d2svg.Render(diagram, &d2svg.RenderOpts{
		Center:        true,
		SetDimensions: true,
		Sketch:        true,
		Pad:           300,
		ThemeID:       d2themescatalog.GrapeSoda.ID,
	})

	f, err := os.Create("imports.svg")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}(f)

	if _, err := f.Write(out); err != nil {
		return err
	}

	return nil
}
