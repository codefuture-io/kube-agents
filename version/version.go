/*
Copyright 2026 CodeFuture Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package version

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// These are set during build time via -ldflags.
var (
	Version   = "latest"
	GitCommit = "N/A"
	BuildDate = "N/A"
)

// Info holds the version information of kube-agents.
type Info struct {
	Version      string `json:"version"`
	GitCommit    string `json:"git_commit"`
	BuildDate    string `json:"build_date"`
	GoVersion    string `json:"go_version"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
	RuntimeCores int    `json:"runtime_cores"`
	TotalMem     int    `json:"total_mem"`
}

// GetVersion returns the version information.
func GetVersion() Info {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return Info{
		Version:      Version,
		GitCommit:    GitCommit,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		RuntimeCores: runtime.GOMAXPROCS(0),
		TotalMem:     int(memStats.TotalAlloc / 1024),
	}
}

var (
	Blue  = color.New(color.FgHiBlue, color.Bold).SprintFunc()
	Green = color.New(color.FgHiGreen, color.Bold).SprintFunc()
)

// Print outputs version information in a table format.
func Print() {
	v := GetVersion()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{
		"Version", "Git Commit", "Build Date",
		"Go Version", "Compiler", "Platform",
		"Runtime Cores", "Total Memory",
	})

	t.AppendRow([]interface{}{
		v.Version, v.GitCommit, v.BuildDate,
		v.GoVersion, v.Compiler, v.Platform,
		strconv.Itoa(v.RuntimeCores) + " cores",
		strconv.Itoa(v.TotalMem) + " KB",
	})

	t.SetStyle(table.StyleDefault)
	t.Style().Format.Header = text.FormatUpper
	t.Style().Color.Header = text.Colors{text.FgHiBlue}
	t.Style().Options.SeparateRows = true

	t.Render()
}

// Term returns the terminal banner string.
func Term() string {
	return fmt.Sprint(Blue(`
╭╮╭━╮╱╱╭╮╱╱╱╱╱╱╱╭━━━╮╱╱╱╱╱╱╱╱╭╮
┃┃┃╭╯╱╱┃┃╱╱╱╱╱╱╱┃╭━╮┃╱╱╱╱╱╱╱╭╯╰╮
┃╰╯╯╭╮╭┫╰━┳━━╮╱╱┃┃╱┃┣━━┳━━┳━╋╮╭╋━━╮
┃╭╮┃┃┃┃┃╭╮┃┃━╋━━┫╰━╯┃╭╮┃┃━┫╭╮┫┃┃━━┫
┃┃┃╰┫╰╯┃╰╯┃┃━╋━━┫╭━╮┃╰╯┃┃━┫┃┃┃╰╋━━┃
╰╯╰━┻━━┻━━┻━━╯╱╱╰╯╱╰┻━╮┣━━┻╯╰┻━┻━━╯
╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╭━╯┃
╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╰━━╯
`))
}
