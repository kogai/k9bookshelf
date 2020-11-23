package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/Yamashou/gqlgenc/clientgen"
	"github.com/Yamashou/gqlgenc/config"
	"github.com/Yamashou/gqlgenc/generator"
)

func mutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		if !regexp.MustCompile("Input$").MatchString(model.Name) {
			continue
		}
		for _, field := range model.Fields {
			// this is the logic to add omitempty
			omit := strings.TrimSuffix(field.Tag, `"`)
			field.Tag = fmt.Sprintf(`%v,omitempty"`, omit)
		}
	}
	return b
}

func main() {
	ctx := context.Background()
	cfg, err := config.LoadConfig(".gqlgenc.yml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err.Error())
		os.Exit(2)
	}
	clientPlugin := clientgen.New(cfg.Query, cfg.Client)
	p := modelgen.Plugin{
		MutateHook: mutateHook,
	}
	if err := generator.Generate(ctx, cfg, api.NoPlugins(), api.AddPlugin(&p), api.AddPlugin(clientPlugin)); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err.Error())
		os.Exit(4)
	}
}
