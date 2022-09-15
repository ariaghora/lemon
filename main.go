package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	lm "github.com/ariaghora/lemon/pkg"
	"github.com/maja42/ember"
	"github.com/tidwall/gjson"
)

var marker = "~~MagicMarker for maja42/ember/v1~~"

// A dirty little trick to prevent marker being removed by go compiler's optomizer
func VOID(x ...interface{}) {}

func GetBuildMode() (*ember.Attachments, string) {
	attachments, err := ember.Open()
	if err != nil {
		return nil, "none"
	}

	if attachments.Count() == 0 {
		return attachments, "none"
	} else {
		configReader := attachments.Reader("game.json")
		if configReader == nil {
			return attachments, "none"
		} else {
			jsonStr, err := io.ReadAll(configReader)
			if err != nil {
				return attachments, "none"
			}
			return attachments, gjson.GetBytes(jsonStr, "build_mode").String()
		}
	}
}

func newGameWithConfig(rootDir string, configBuf []byte, ignoreBuildMode bool) *lm.LemonGame {
	screenWidth := gjson.GetBytes(configBuf, "screen_width").Float()
	screenHeight := gjson.GetBytes(configBuf, "screen_height").Float()
	targetFPS := gjson.GetBytes(configBuf, "target_fps").Float()
	title := gjson.GetBytes(configBuf, "title").String()
	startingScene := gjson.GetBytes(configBuf, "starting_scene").String()
	buildMode := "none"
	if !ignoreBuildMode {
		buildMode = gjson.GetBytes(configBuf, "build_mode").String()
	}

	lemonGame := lm.NewLemonGame(
		rootDir,
		startingScene,
		int(screenWidth),
		int(screenHeight),
		title,
		float32(targetFPS),
		buildMode,
	)

	return lemonGame
}

func main() {
	VOID(marker)

	attachments, buildMode := GetBuildMode()
	if attachments != nil {
		defer attachments.Close()
	}

	if buildMode == "standalone" {
		configBuf, err := io.ReadAll(attachments.Reader("game.json"))
		if err != nil {
			panic(err)
		}
		lemonGame := newGameWithConfig("", configBuf, false)
		lemonGame.BuildMode = buildMode
		defer lemonGame.Close()

		lemonGame.Run()
	} else {
		// Okay, okay, gonna need a better argparser. Now I am too tired
		// and sad to do so.
		if len(os.Args) < 3 {
			fmt.Println("not enough arguments")
			os.Exit(1)
		}

		cmd := os.Args[1]
		if !(cmd == "run" || cmd == "build") {
			panic("invalid command: " + cmd)
		}

		rootDir := os.Args[2]
		if cmd == "run" {
			if _, err := os.Stat(path.Join(rootDir, "game.json")); errors.Is(err, os.ErrNotExist) {
				fmt.Println("config file game.json not found")
				os.Exit(1)
			}

			f, err := os.Open(path.Join(rootDir, "game.json"))
			if err != nil {
				fmt.Println("cannot open game.json")
				os.Exit(1)
			}
			buf, err := io.ReadAll(f)
			if err != nil {
				fmt.Println("cannot read game.json")
				os.Exit(1)
			}
			lemonGame := newGameWithConfig(rootDir, buf, true)
			lemonGame.BuildMode = "none"
			defer lemonGame.Close()
			lemonGame.Run()
		} else if cmd == "build" {
			lm.Build(rootDir)
		}
	}
}
