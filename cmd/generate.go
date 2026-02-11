package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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

// runGeneration duyệt qua danh sách Entries trong cấu trúc Config mới
func runGeneration(cmd *cobra.Command) error {
	configFile, _ := cmd.Flags().GetString("config")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := helpers.LoadConfig(path.Join(cwd, configFile))
	if err != nil {
		return err
	}

	for _, entry := range config.Entries {
		sourcePath := path.Join(cwd, entry.Source)

		err = filepath.WalkDir(sourcePath, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && filepath.Ext(filePath) == ".contract" {
				ast, pErr := ParseFile(filePath)
				if pErr != nil {
					parser.PrintError(pErr, "")
					return nil
				}

				relPath, rErr := filepath.Rel(sourcePath, filePath)
				if rErr != nil {
					return rErr
				}

				baseName := relPath[:len(relPath)-len(".contract")]

				if strings.Contains(entry.Lang, "ts") {
					outputPath := filepath.Join(cwd, entry.Output, baseName+".ts")
					tsGen := generator.NewTypescriptGenerator()
					code, gErr := tsGen.Generate(ast)
					if gErr != nil {
						parser.PrintError(gErr, "")
					} else {
						helpers.WriteTextFile(outputPath, code)
						fmt.Printf("Generated TypeScript: %s -> %s\n", filePath, outputPath)
					}
				}

				if strings.Contains(entry.Lang, "go") {
					outputPath := filepath.Join(cwd, entry.Output, baseName+".go")
					goGen := generator.NewGoGenerator()
					code, gErr := goGen.Generate(ast, entry.PkgName, "")
					if gErr != nil {
						parser.PrintError(gErr, "")
					} else {
						helpers.WriteTextFile(outputPath, code)
						fmt.Printf("Generated Go:%s -> %s\n", filePath, outputPath)
					}
				}
			}
			return nil
		})
	}
	return err
}

func GenerateCommandFn(cmd *cobra.Command, args []string) error {
	fmt.Println("Initial generation started...")
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

	for _, entry := range config.Entries {
		sourcePath := path.Join(cwd, entry.Source)
		filepath.WalkDir(sourcePath, func(p string, d fs.DirEntry, err error) error {
			if d != nil && d.IsDir() {
				return watcher.Add(p)
			}
			return nil
		})
	}

	fmt.Println("Watching for changes in defined source paths...")

	var timer *time.Timer
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if (event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) && strings.HasSuffix(event.Name, ".contract") {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(100*time.Millisecond, func() {
					fmt.Printf("Change detected at %s. Regenerating...\n", time.Now().Format("15:04:05"))
					if err := runGeneration(cmd); err != nil {
						log.Printf("Generation error: %v\n", err)
					}
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher error: %v\n", err)
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
