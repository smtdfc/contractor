package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/smtdfc/contractor/generator"
	"github.com/smtdfc/contractor/helpers"
	"github.com/smtdfc/contractor/parser"
	"github.com/spf13/cobra"
)

func ParseFile(filePath string) (*parser.AST, parser.BaseError) {
	content, err := helpers.ReadTextFile(filePath)
	if err != nil {
		panic(err)
	}

	lexer := parser.NewLexer()
	tokens, lexErr := lexer.Parse(content, filePath)
	if lexErr != nil {
		return nil, lexErr
	}

	p := parser.NewParser()
	ast, parseErr := p.Parse(tokens, filePath)
	if parseErr != nil {
		return nil, parseErr
	}

	typeChecker := parser.NewTypeChecker()
	typeErr := typeChecker.Check(ast)
	if typeErr != nil {
		return nil, typeErr
	}

	return ast, nil
}

func GenerateTsCode(ast *parser.AST, output string) (string, parser.BaseError) {
	tsGenerator := generator.NewTypescriptGenerator()
	code, err := tsGenerator.Generate(ast)
	if err != nil {
		return "", err
	}

	wErr := helpers.WriteTextFile(output, code)
	if wErr != nil {
		panic(wErr)
	}

	return code, nil
}

func runGeneration(cmd *cobra.Command) error {
	configFile := "contractor.config.json"
	configFlag, _ := cmd.Flags().GetString("config")
	if configFlag != "" {
		configFile = configFlag
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := helpers.LoadConfig(path.Join(cwd, configFile))
	if err != nil {
		return err
	}

	sourcePath := path.Join(cwd, config.Source)

	err = filepath.WalkDir(sourcePath, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(filePath) == ".contract" {
			info, err := d.Info()
			if err != nil || info.IsDir() {
				return nil
			}

			ast, pErr := ParseFile(filePath)
			if pErr != nil {
				parser.PrintError(pErr, "")
				return nil
			}

			relPath, rErr := filepath.Rel(sourcePath, filePath)
			if rErr != nil {
				return rErr
			}

			if slices.Contains(config.Lang, "ts") {
				outputPath := filepath.Join(config.Output, relPath)
				outputPath = outputPath[:len(outputPath)-len(".contract")] + ".ts"

				_, gErr := GenerateTsCode(ast, outputPath)
				if gErr != nil {
					parser.PrintError(gErr, "")
					return nil
				}

				fmt.Printf("Generated: %s -> %s\n", filePath, outputPath)
			}
		}
		return nil
	})

	return err
}

func GenerateCommandFn(cmd *cobra.Command, args []string) error {

	fmt.Println("Starting initial generation...")
	if err := runGeneration(cmd); err != nil {
		return err
	}

	isWatch, _ := cmd.Flags().GetBool("watch")
	if !isWatch {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	configFile, _ := cmd.Flags().GetString("config")
	cwd, _ := os.Getwd()
	config, err := helpers.LoadConfig(path.Join(cwd, configFile))
	if err != nil {
		return err
	}

	sourcePath := path.Join(cwd, config.Source)
	filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if d != nil && d.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	fmt.Printf("Watching for changes in: %s\n", sourcePath)

	var timer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if (event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) && filepath.Ext(event.Name) == ".contract" {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(100*time.Millisecond, func() {
					fmt.Printf("\nChange detected at %s. Re-generating...\n", time.Now().Format("15:04:05"))
					if err := runGeneration(cmd); err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Println("Watcher error:", err)
		}
	}
}

var GenerateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from contracts",
	Long:  "Generate code from contracts and optionally watch for changes",
	RunE:  GenerateCommandFn,
}

func InitAllCommand(root *cobra.Command) {
	root.AddCommand(InitCommand)

	GenerateCommand.Flags().StringP("config", "c", "contractor.config.json", "Path to config file")
	GenerateCommand.Flags().BoolP("watch", "w", false, "Watch for file changes and regenerate automatically")

	root.AddCommand(GenerateCommand)
}
