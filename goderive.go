package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/nextzhou/goderive/plugin"
	"github.com/nextzhou/goderive/plugin/set"
	"github.com/nextzhou/goderive/utils"
	"github.com/spf13/cobra"
)

func main() {
	defer func() {
		if info := recover(); info != nil {
			fmt.Fprintln(os.Stderr, `╔════════════════════════════════════════════════════════════════════════════════╗`)
			fmt.Fprintln(os.Stderr, `║NOTICE: You found a bug!!!                                                      ║`)
			fmt.Fprintln(os.Stderr, `║Please report bug to https://github.com/nextzhou/goderive/issues with log below.║`)
			fmt.Fprintln(os.Stderr, `╚════════════════════════════════════════════════════════════════════════════════╝`)
			fmt.Fprintf(os.Stderr, "panic: %v\n\n%s\n", info, debug.Stack())
			os.Exit(1)
		}
	}()

	derive := NewDerive()
	derive.RegisterPlugin(set.Set{})

	if err := derive.Execute(); err != nil {
		os.Exit(1)
	}
}

type Derive struct {
	Plugins []plugin.Plugin
	Cmd     *cobra.Command
	Err     error
	Output  string
	Delete  bool
}

func NewDerive() *Derive {
	derive := new(Derive)
	derive.Cmd = &cobra.Command{
		Use: "goderive",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 && args[0] == "help" {
				return derive.Help(args[1:])
			}
			return derive.Run(args)
		},
		SilenceUsage: true,
	}
	derive.Cmd.Flags().StringVarP(&derive.Output, "output", "o", "derived.gen.go", "output file name")
	derive.Cmd.Flags().BoolVarP(&derive.Delete, "delete", "d", true, "delete existing generated file when no derived type")
	return derive
}

func (d *Derive) RegisterPlugin(plugins ...plugin.Plugin) {
	// TODO uniq
	// TODO Description check
	d.Plugins = append(d.Plugins, plugins...)
}

func (d *Derive) Execute() error {
	// set help template after plugin registration
	d.Cmd.SetHelpTemplate(d.HelpString())
	if d.Err != nil {
		return d.Err
	}
	return d.Cmd.Execute()
}

func (d Derive) HelpString() string {
	help := bytes.NewBufferString(`GoDerive

Add derive comment above your type, and generate source code for the marked type.

Comment Format:
  // derive-<plugin>
  // derive-<plugin>: flag;!negative_flag;arg=single_value; arg2=val1,val2
  type YourType struct{/* ... */}

Usage:
  goderive [flags] [path ...]
  goderive help [plugin ...]

Flags:
  -d, --delete          delete existing generated file when no derived type (default true)
  -h, --help            help for goderive
  -o, --output string   output file name (default "derived.gen.go")

Plugins:
`)
	w := utils.NewTableWriter(help)
	for _, plugin := range d.Plugins {
		desc := plugin.Describe()
		w.Append([]string{desc.Identity, desc.Effect})
	}
	w.Render()
	return help.String()
}

func ListGoFiles(path string) ([]string, error) {
	// TODO support ./...
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var files []string
	if stat.IsDir() {
		dirInfo, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, entry := range dirInfo {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
				files = append(files, filepath.Join(path, entry.Name()))
			}
		}
	} else {
		if strings.HasSuffix(path, ".go") {
			files = []string{path}
		} else {
			return nil, fmt.Errorf("%#v is not a go source file", path)
		}
	}
	return files, nil
}

func (d *Derive) Run(inputPaths []string) error {
	// scan go source file
	// TODO uniq
	if len(inputPaths) == 0 {
		inputPaths = []string{"."}
	}
	var files []string
	for _, path := range inputPaths {
		fs, err := ListGoFiles(path)
		if err != nil {
			return err
		}
		files = append(files, fs...)
	}

	// extract type info, and group them by package(path)
	groupFileInfoByPath := make(map[string]*FileInfo)
	for _, file := range files {
		path, err := filepath.Abs(filepath.Dir(file))
		if err != nil {
			panic(err)
		}
		fi, ok := groupFileInfoByPath[path]
		if !ok {
			fi = new(FileInfo)
			groupFileInfoByPath[path] = fi

		}

		src, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %#v : %s", file, err.Error())
		}
		fileInfo, err := ExtractTypes(src)
		if err != nil {
			return fmt.Errorf("%s: %s", file, err.Error())
		}
		if fileInfo.PkgName == "" || len(fileInfo.Types) == 0 {
			continue
		}
		for _, typ := range fileInfo.Types {
			for pluginID, opts := range typ.Plugins {
				if err := d.ValidatePluginOptions(pluginID, opts); err != nil {
					return fmt.Errorf("%#v: type %s: %v", file, typ.Name, err)
				}
			}
		}
		if fi.PkgName == "" {
			fi.PkgName = fileInfo.PkgName
		}
		fi.Types = append(fi.Types, fileInfo.Types...)
		groupFileInfoByPath[path] = fi
	}

	var shouldDeletedFiles []string
	for path, fileInfo := range groupFileInfoByPath {
		filename := filepath.Join(path, d.Output)
		if fileInfo.PkgName == "" {
			if d.Delete {
				shouldDeletedFiles = append(shouldDeletedFiles, filename)
			}
			continue
		}
		headBuf := bytes.NewBuffer(nil)
		headBuf.WriteString(fmt.Sprintf("package %s\n\n", fileInfo.PkgName))
		imports := make([]string, 0)
		bodyBuf := bytes.NewBuffer(nil)
		for _, typ := range fileInfo.Types {
			for pluginID, opts := range typ.Plugins {
				plugin, _ := d.GetPlugin(pluginID)
				prerequisites, err := plugin.GenerateTo(bodyBuf, typ.Name, typ.Ast, *opts)
				if err != nil {
					// TODO log file path of type
					return fmt.Errorf("failed to generate code of type %s: %v", typ.Name, err)
				}
				// FIXME: imports uniq
				imports = append(imports, prerequisites.Imports...)
			}
		}
		// TODO write file after all generating

		if len(imports) > 0 {
			headBuf.WriteString("import (\n")
			for _, i := range imports {
				headBuf.WriteString(fmt.Sprintf("\t%#v\n", i))
			}
			headBuf.WriteString(")\n")
		}

		headBuf.Write(bodyBuf.Bytes())

		err := ioutil.WriteFile(filename, headBuf.Bytes(), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write %#v error : %s\n", filename, err.Error())
			os.Exit(1)
		}
	}

	for _, file := range shouldDeletedFiles {
		// ignore errors
		os.Remove(file)
	}
	return nil
}

func (d *Derive) Help(pluginID []string) error {
	if len(pluginID) == 0 {
		fmt.Println(d.HelpString())
		return nil
	}
	help := bytes.NewBuffer(nil)
	for _, topic := range pluginID {
		plugin, err := d.GetPlugin(topic)
		if err != nil {
			return err
		}
		help.WriteString(plugin.Describe().ToHelpString())
		help.WriteByte('\n')
	}
	fmt.Println(help.String())
	return nil
}

func (d *Derive) ValidatePluginOptions(pluginID string, opts *plugin.Options) error {
	plugin, err := d.GetPlugin(pluginID)
	if err != nil {
		return err
	}
	return plugin.Describe().Validate(opts)
}

func (d *Derive) GetPlugin(pluginID string) (plugin.Plugin, error) {
	for _, plugin := range d.Plugins {
		if plugin.Describe().Identity == pluginID {
			return plugin, nil
		}
	}
	return nil, &utils.UnsupportedError{Type: "plugin", Idents: []string{pluginID}}
}
