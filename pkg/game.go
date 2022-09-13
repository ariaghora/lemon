package pkg

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	lua "github.com/yuin/gopher-lua"
)

//go:embed lua/core.lua
var coreScript string

//go:embed lua/sprite.lua
var spriteScript string

//go:embed lua/texture.lua
var textureScript string

//go:embed lua/keymaps.lua
var keymapsScript string

var loadedModules = []string{}

func isin(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func getOverriderRequireFuncScript(root string) string {
	overrideRequire := `
	oldreq = require
	require = function(s) 
		if (
			(s == "package") or
			(s == "table") or
			(s == "io") or
			(s == "os") or
			(s == "string") or
			(s == "math") or
			(s == "debug") or
			(s == "channel") or
			(s == "coroutine") 
		)then
			return oldreq(s)
		end
		return myrequire(s, "%s")
	end
	`
	return fmt.Sprintf(overrideRequire, root)
}

func lua_myrequire(l *lua.LState) int {
	modname := l.ToString(1)
	root := l.ToString(2)

	if isin(modname, loadedModules) {
		return 0
	}

	loadedModules = append(loadedModules, modname)
	filepath := path.Join(strings.Split(modname, ".")...) + ".lua"
	filepath = path.Join(root, filepath)
	if err := l.DoFile(filepath); err != nil {
		panic(err)
	}
	return 0
}

type GameConfig struct {
	Title        string
	WindowWidth  int
	WindowHeight int
	TargetFPS    int
}

type LemonGame struct {
	RootDir            string
	lState             *lua.LState
	currentSceneLoaded bool
	spriteDict         map[string]*LemonSprite
	spriteNames        []string
	LOnloadFunc        lua.LValue
	lOnUpdateFunc      lua.LValue
	textureDict        map[int]*TextureInfo

	spriteCount  int
	textureCount int
}

func NewLemonGame(
	rootDir string,
	startingScene string,
	screenWidth int,
	screenHeight int,
	title string,
	targetFPS float32,
) *LemonGame {
	lemonGame := &LemonGame{
		RootDir:      rootDir,
		lState:       lua.NewState(),
		spriteDict:   map[string]*LemonSprite{},
		spriteCount:  0,
		textureCount: 0,
		textureDict:  map[int]*TextureInfo{},
	}

	l_Lemon := &lua.LTable{}
	funcMap := &map[string]lua.LGFunction{
		"draw_rect_fill":            lemonGame.lua_draw_rect_fill,
		"draw_text":                 lemonGame.lua_draw_text,
		"draw_texture":              lemonGame.lua_draw_texture,
		"find_sprites_by_name_like": lemonGame.lua_findSpritesByNameLike,
		"get_fps":                   lemonGame.lua_getFPS,
		"get_screen_height":         lemonGame.lua_getScreenHeight,
		"get_screen_width":          lemonGame.lua_getScreenWidth,
		"is_key_pressed":            lemonGame.lua_isKeyPressed,
		"is_key_down":               lemonGame.lua_isKeyDown,
		"is_key_up":                 lemonGame.lua_isKeyUp,
		"new_sprite":                lemonGame.lua_newSprite,
		"new_texture":               lemonGame.lua_newTexture,
		"set_scene":                 lemonGame.lua_setScene,
	}
	lemonGame.lState.SetFuncs(l_Lemon, *funcMap)

	lemonGame.lState.SetGlobal("myrequire", lemonGame.lState.NewFunction(lua_myrequire))

	// Put `Lemon` table in the global scope
	lemonGame.lState.SetGlobal("Lemon", l_Lemon)
	lemonGame.lState.DoString("L = Lemon")

	// Populate initial Lemon methods
	lemonGame.lState.DoString(coreScript)

	// Populate initial prototype methods
	lemonGame.lState.DoString(spriteScript)
	lemonGame.lState.DoString(textureScript)

	// Populate key code attributes
	lemonGame.lState.DoString(keymapsScript)

	// Override Lua's default "require" function
	lemonGame.lState.DoString(getOverriderRequireFuncScript(rootDir))

	lemonGame.SetScene(startingScene)

	// Some basic raylib initial settings
	rl.SetTraceLog(rl.LogNone)
	rl.InitWindow(int32(screenWidth), int32(screenHeight), title)
	rl.SetTargetFPS(int32(targetFPS))

	return lemonGame
}

func (lm *LemonGame) Close() {
	lm.lState.Close()
}

func (lm *LemonGame) OnDraw(dt float32) {
	for _, name := range lm.spriteNames {
		lm.spriteDict[name].DoDraw(float64(dt))
	}
}

func (lm *LemonGame) OnLoad() {
	if !lm.currentSceneLoaded {
		if lm.LOnloadFunc.Type() != lua.LTNil {
			err := lm.lState.CallByParam(
				lua.P{
					Fn:      lm.LOnloadFunc,
					NRet:    0,
					Protect: true,
				},
			)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lm.currentSceneLoaded = true
		}
	}
}

func (lm *LemonGame) OnUpdate(dt float64) {
	if lm.lOnUpdateFunc.Type() != lua.LTNil {
		lm.lState.CallByParam(
			lua.P{
				Fn:      lm.lOnUpdateFunc,
				NRet:    0,
				Protect: true,
			},
			lua.LNumber(rl.GetFrameTime()),
		)
	}

	// sprite update:
	// 1) first pass: actual update
	for _, name := range lm.spriteNames {
		lm.spriteDict[name].DoUpdate(float64(dt))
	}

	// second pass: delete sprites that should be deleted
	for _, name := range lm.spriteNames {
		lm.spriteDict[name].DoRemove()
	}
}

func (lm *LemonGame) Run() {

	for !rl.WindowShouldClose() {
		lm.OnLoad()

		rl.BeginDrawing()
		rl.ClearBackground(rl.White)
		lm.OnDraw(rl.GetFrameTime())

		rl.EndDrawing()
		lm.OnUpdate(float64(rl.GetFrameTime()))
	}
	rl.CloseWindow()
}

func (lm *LemonGame) SetScene(sceneName string) {
	sceneSourceDir := path.Join(lm.RootDir, sceneName+".lua")
	if _, err := os.Stat(sceneSourceDir); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("failed to switch to scene \"%s\" (cannot open %s)\n", sceneName, sceneSourceDir)
		os.Exit(1)
	}

	// Clear current state first
	lm.currentSceneLoaded = false
	for k := range lm.spriteDict {
		delete(lm.spriteDict, k)
	}
	lm.spriteNames = []string{}

	err := lm.lState.DoFile(sceneSourceDir)
	if err != nil {
		panic(err)
	}
	lm.LOnloadFunc = lm.lState.GetGlobal("on_load")
	lm.lOnUpdateFunc = lm.lState.GetGlobal("on_update")
}

// ==========================================
// lua function interface and implementations
// ==========================================

func (lm *LemonGame) lua_draw_rect_fill(L *lua.LState) int {
	x, y, w, h := L.ToNumber(1), L.ToNumber(2), L.ToNumber(3), L.ToNumber(4)
	color := L.ToTable(5)
	rl.DrawRectangle(
		int32(x),
		int32(y),
		int32(w),
		int32(h),
		rl.Color{
			R: uint8(color.RawGetString("r").(lua.LNumber)),
			G: uint8(color.RawGetString("g").(lua.LNumber)),
			B: uint8(color.RawGetString("b").(lua.LNumber)),
			A: uint8(color.RawGetString("a").(lua.LNumber)),
		},
	)
	return 0
}

func (lm *LemonGame) lua_draw_text(L *lua.LState) int {
	text, x, y, fs := L.ToString(1), L.ToNumber(2), L.ToNumber(3), L.ToNumber(4)
	color := L.ToTable(5)
	rl.DrawText(
		text,
		int32(x),
		int32(y),
		int32(fs),
		rl.Color{
			R: uint8(color.RawGetString("r").(lua.LNumber)),
			G: uint8(color.RawGetString("g").(lua.LNumber)),
			B: uint8(color.RawGetString("b").(lua.LNumber)),
			A: uint8(color.RawGetString("a").(lua.LNumber)),
		},
	)
	return 0
}

func (lm *LemonGame) lua_draw_texture(L *lua.LState) int {
	texture := L.ToTable(1)
	id := int(texture.RawGetString("id").(lua.LNumber))
	x, y := L.ToNumber(2), L.ToNumber(3)
	rl.DrawTexture(lm.textureDict[id].TextureData, int32(x), int32(y), rl.White)
	return 0
}

func (lm *LemonGame) lua_findSpritesByNameLike(L *lua.LState) int {
	queryStr := L.ToString(1)
	table := lua.LTable{}
	for _, name := range lm.spriteNames {
		matched := strings.Contains(name, queryStr)
		if matched {
			table.Append(&lm.spriteDict[name].LTable)
		}
	}
	L.Push(&table)
	return 1
}

func (lm *LemonGame) lua_getFPS(L *lua.LState) int {
	result := rl.GetFPS()
	L.Push(lua.LNumber(result))
	return 1
}

func (lm *LemonGame) lua_getScreenHeight(L *lua.LState) int {
	L.Push(lua.LNumber(rl.GetScreenHeight()))
	return 1
}

func (lm *LemonGame) lua_getScreenWidth(L *lua.LState) int {
	L.Push(lua.LNumber(rl.GetScreenWidth()))
	return 1
}

func (lm *LemonGame) lua_isKeyDown(L *lua.LState) int {
	keyCode := L.ToNumber(1)
	result := rl.IsKeyDown(int32(keyCode))
	L.Push(lua.LBool(result))
	return 1
}

func (lm *LemonGame) lua_isKeyUp(L *lua.LState) int {
	keyCode := L.ToNumber(1)
	result := rl.IsKeyUp(int32(keyCode))
	L.Push(lua.LBool(result))
	return 1
}

func (lm *LemonGame) lua_isKeyPressed(L *lua.LState) int {
	keyCode := L.ToNumber(1)
	result := rl.IsKeyPressed(int32(keyCode))
	L.Push(lua.LBool(result))
	return 1
}

func (lm *LemonGame) lua_newSprite(L *lua.LState) int {
	name := L.ToString(1)
	NewSprite(name, L, lm)
	return 1
}

func (lm *LemonGame) lua_newTexture(L *lua.LState) int {
	fileName := path.Join(lm.RootDir, L.ToString(1))
	tex := NewTexture(fileName, lm)
	L.Push(tex.LTable)
	return 1
}

func (lm *LemonGame) lua_setScene(L *lua.LState) int {
	sceneName := L.ToString(1)
	lm.SetScene(sceneName)
	return 0
}
