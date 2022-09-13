package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	lm "github.com/ariaghora/lemon/pkg"
	"github.com/tidwall/gjson"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please specify project directory")
		os.Exit(1)
	}
	rootDir := os.Args[1]
	if _, err := os.Stat(path.Join(rootDir, "game.json")); errors.Is(err, os.ErrNotExist) {
		fmt.Println(path.Join(rootDir, "game.json"))
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
	screenWidth := gjson.GetBytes(buf, "screen_width").Float()
	screenHeight := gjson.GetBytes(buf, "screen_height").Float()
	targetFPS := gjson.GetBytes(buf, "target_fps").Float()
	title := gjson.GetBytes(buf, "title").String()
	startingScene := gjson.GetBytes(buf, "starting_scene").String()

	lemonGame := lm.NewLemonGame(
		rootDir,
		startingScene,
		int(screenWidth),
		int(screenHeight),
		title,
		float32(targetFPS),
	)
	defer lemonGame.Close()

	lemonGame.Run()
}
