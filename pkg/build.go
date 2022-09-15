package pkg

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/maja42/ember/embedding"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/gjson"
)

func Build(rootDir string) {
	if _, err := os.Stat(path.Join(rootDir, "game.json")); err != nil {
		fmt.Println("cannot open", path.Join(rootDir, "game.json"))
	}

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exe, err := os.Open(exePath)
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Open(path.Join(rootDir, "game.json"))
	if err != nil {
		panic(err)
	}
	buf, _ := io.ReadAll(jsonFile)
	gameTitle := gjson.GetBytes(buf, "title").String()

	// For some reason outPath cannot be set into subdirectories.
	// For a little hack, we move this file later.
	outPath := "__lemon__.__out__"
	os.Remove(outPath)
	out, err := os.OpenFile(outPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	filemap := map[string]string{}
	err = filepath.Walk(rootDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			memoryFileName := strings.Replace(path, rootDir, "", -1)
			memoryFileName = strings.TrimPrefix(memoryFileName, "/")
			filemap[memoryFileName] = path
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	bar := progressbar.NewOptions(
		len(filemap),
		progressbar.OptionSetDescription("Building \""+gameTitle+"\""),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(
			progressbar.Theme{
				Saucer:     "[yellow]=[reset]",
				SaucerHead: ">",
				BarStart:   "[",
				BarEnd:     "]",
			},
		),
	)
	logger := func(format string, args ...interface{}) {
		bar.Add(1)
	}
	err = embedding.EmbedFiles(out, exe, filemap, logger)
	if err != nil {
		panic(err)
	}

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	err = os.Rename(outPath, path.Join(rootDir, gameTitle)+ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n✨Game built: \"%s\"✨\n", path.Join(rootDir, gameTitle))
}
