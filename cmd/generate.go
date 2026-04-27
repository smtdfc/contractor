package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/smtdfc/contractor/emitters"
	"github.com/smtdfc/contractor/emitters/typescript"
	"github.com/smtdfc/contractor/generator"
	"github.com/smtdfc/contractor/internal/config"
	"github.com/smtdfc/contractor/parser"
	"github.com/spf13/cobra"
)

var configPath string
var generateLang string

func init() {
	generateCmd.Flags().StringVarP(&configPath, "config", "c", "contractor.json", "Path to contractor config file")
	generateCmd.Flags().StringVarP(&generateLang, "lang", "l", "", "Generate only for this language (e.g. go, typescript, java, kotlin, csharp)")
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code from .contract files",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		targets, err := selectTargets(cfg.Targets, generateLang)
		if err != nil {
			return err
		}

		files, err := findContractFiles(cfg.SourceDir, cfg.Extension)
		if err != nil {
			return fmt.Errorf("scan source dir: %w", err)
		}

		if len(files) == 0 {
			cmd.Printf("No %s files found in %s\n", cfg.Extension, cfg.SourceDir)
			return nil
		}

		for _, filePath := range files {
			relPath, err := filepath.Rel(cfg.SourceDir, filePath)
			if err != nil {
				return fmt.Errorf("resolve relative path for %s: %w", filePath, err)
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("read %s: %w", filePath, err)
			}

			ir, err := parseProgram(filePath, string(content))
			if err != nil {
				return err
			}

			for _, target := range targets {
				emitter, ext, err := resolveEmitter(target.Language)
				if err != nil {
					return err
				}

				output, err := emitter.Emit(ir)
				if err != nil {
					return fmt.Errorf("emit %s for %s: %w", target.Language, filePath, err)
				}

				outFilePath := outputPathForTarget(target.OutDir, relPath, ext)
				if err := os.MkdirAll(filepath.Dir(outFilePath), 0o755); err != nil {
					return fmt.Errorf("create output dir for %s: %w", outFilePath, err)
				}

				if err := os.WriteFile(outFilePath, []byte(output), 0o644); err != nil {
					return fmt.Errorf("write output file %s: %w", outFilePath, err)
				}

				cmd.Printf("Generated %s\n", outFilePath)
			}
		}

		return nil
	},
}

func findContractFiles(sourceDir string, extension string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.WalkDir(sourceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.EqualFold(filepath.Ext(path), extension) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func parseProgram(filePath string, code string) (*generator.ProgramIR, error) {
	lexer := parser.NewLexer(filePath)
	tokens, err := lexer.Start(code)
	if err != nil {
		return nil, fmt.Errorf("lex %s: %w", filePath, err)
	}

	p := parser.NewParser(filePath, tokens)
	ast, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", filePath, err)
	}

	typeChecker := parser.NewTypeChecker()
	if err := typeChecker.Check(ast); err != nil {
		return nil, fmt.Errorf("type-check %s: %w", filePath, err)
	}

	irGenerator := generator.NewIRGenerator()
	ir, err := irGenerator.GenerateProgram(ast)
	if err != nil {
		return nil, fmt.Errorf("generate IR %s: %w", filePath, err)
	}

	return ir, nil
}

func resolveEmitter(language string) (emitters.ProgramEmitter, string, error) {
	switch normalizeLanguage(language) {
	// case "go", "golang":
	// 	return golang.NewGoEmitter(), ".go", nil
	case "typescript", "ts":
		return typescript.NewTypescriptEmitter(), ".ts", nil
	// case "java":
	// 	return java.NewJavaEmitter(), ".java", nil
	// case "kotlin", "kt":
	// 	return kotlin.NewKotlinEmitter(), ".kt", nil
	// case "csharp", "cs", "c#":
	// 	return csharp.NewCSharpEmitter(), ".cs", nil
	default:
		return nil, "", fmt.Errorf("unsupported target language: %s", language)
	}
}

func selectTargets(targets []config.Target, requestedLanguage string) ([]config.Target, error) {
	if strings.TrimSpace(requestedLanguage) == "" {
		return targets, nil
	}

	requestedCanonical, err := canonicalLanguage(requestedLanguage)
	if err != nil {
		return nil, fmt.Errorf("invalid --lang value: %w", err)
	}

	filtered := make([]config.Target, 0, len(targets))
	for _, target := range targets {
		targetCanonical, err := canonicalLanguage(target.Language)
		if err != nil {
			continue
		}

		if targetCanonical == requestedCanonical {
			filtered = append(filtered, target)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no target language %q found in config targets", requestedLanguage)
	}

	return filtered, nil
}

func canonicalLanguage(language string) (string, error) {
	switch normalizeLanguage(language) {
	case "go", "golang":
		return "go", nil
	case "typescript", "ts":
		return "typescript", nil
	case "java":
		return "java", nil
	case "kotlin", "kt":
		return "kotlin", nil
	case "csharp", "cs", "c#":
		return "csharp", nil
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
}

func normalizeLanguage(language string) string {
	return strings.ToLower(strings.TrimSpace(language))
}

func outputPathForTarget(outDir string, relContractPath string, outExt string) string {
	relDir := filepath.Dir(relContractPath)
	baseName := strings.TrimSuffix(filepath.Base(relContractPath), filepath.Ext(relContractPath))
	fileName := "index" + outExt

	if relDir == "." {
		return filepath.Join(outDir, baseName, fileName)
	}

	return filepath.Join(outDir, relDir, baseName, fileName)
}
