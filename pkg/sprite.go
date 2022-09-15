package pkg

import (
	"fmt"
	"math"
	"os"
	"path"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	lua "github.com/yuin/gopher-lua"
)

type LemonSprite struct {
	name string
	// Texture       rl.Texture2D
	Texture       *TextureInfo
	LTable        lua.LTable
	GameRef       *LemonGame
	LOnLoadFunc   lua.LValue
	LOnUpdateFunc lua.LValue
	scriptLoaded  bool
	shouldRemove  bool
	animCounter   float64
}

const (
	SpriteFnSetScript          = "set_script"
	SpriteFnSetTextureFromFile = "set_texture_from_file"
	SpriteFnSetTexture         = "set_texture"
	SpriteFnRemove             = "remove"
)

func (s *LemonSprite) GetName() string {
	return s.name
}

func NewSprite(name string, L *lua.LState, gameRef *LemonGame) int {
	if gameRef.spriteDict[name] != nil {
		panic(fmt.Sprintf("name '%s' is already used, please use a different name", name))
	}
	o := &LemonSprite{
		name:         name,
		GameRef:      gameRef,
		scriptLoaded: false,
		shouldRemove: false,
		animCounter:  1,
	}
	gameRef.spriteDict[name] = o // keep track
	gameRef.spriteNames = append(gameRef.spriteNames, name)

	// calling prototype function to generate initial sprite table structure
	proto_new_sprite_fn := L.GetGlobal("__proto_new_sprite").(*lua.LFunction)
	L.CallByParam(
		lua.P{
			Fn:   proto_new_sprite_fn,
			NRet: 1, Protect: true,
		},
		lua.LString(name),
	)
	table := L.Get(-1).(*lua.LTable)
	table.RawSetString(SpriteFnSetScript, L.NewFunction(o.lua_setScript))
	table.RawSetString(SpriteFnSetTextureFromFile, L.NewFunction(o.lua_setTextureFromFile))
	table.RawSetString(SpriteFnSetTexture, L.NewFunction(o.lua_setTexture))
	table.RawSetString(SpriteFnRemove, L.NewFunction(o.lua_remove))
	gameRef.spriteDict[name].LTable = *table

	L.Push(table)

	gameRef.spriteCount += 1
	return 0
}

func (l *LemonSprite) DoDraw(dt float64) {
	x := float32(l.LTable.RawGetString("x").(lua.LNumber))
	y := float32(l.LTable.RawGetString("y").(lua.LNumber))
	width := float32(l.LTable.RawGetString("width").(lua.LNumber))
	height := float32(l.LTable.RawGetString("height").(lua.LNumber))
	frameCountX := int(l.LTable.RawGetString("frame_count_x").(lua.LNumber))
	frameCountY := int(l.LTable.RawGetString("frame_count_y").(lua.LNumber))
	framewidth := int(l.LTable.RawGetString("frame_width").(lua.LNumber))
	frameheight := int(l.LTable.RawGetString("frame_height").(lua.LNumber))
	rotation := float32(l.LTable.RawGetString("rotation").(lua.LNumber))
	originX := float32(l.LTable.RawGetString("origin_x").(lua.LNumber))
	originY := float32(l.LTable.RawGetString("origin_y").(lua.LNumber))
	shown := l.LTable.RawGetString("shown").(lua.LBool)
	playing := l.LTable.RawGetString("playing").(lua.LBool)
	dur := float64(l.LTable.RawGetString("animation_duration").(lua.LNumber))

	idx := int(l.LTable.RawGetString("frame_index").(lua.LNumber))
	if shown {
		row := 0
		col := 0
		if playing {
			l.animCounter += dt
			if l.animCounter >= dur {
				l.animCounter -= dur
			}

			idx = int(math.Floor(l.animCounter/dur*float64(frameCountX*frameCountY)) + 1)

			l.LTable.RawSetString("frame_index", lua.LNumber(idx))
		}
		// convert linear index (of frame) to row and column position in the
		// texture
		row = (idx - 1) / frameCountX
		col = (idx - 1) % frameCountX
		rl.DrawTexturePro(
			l.Texture.TextureData,
			rl.Rectangle{X: float32(col * framewidth), Y: float32(row * frameheight), Width: float32(framewidth), Height: float32(frameheight)},
			rl.Rectangle{X: x, Y: y, Width: float32(width), Height: float32(height)},
			rl.Vector2{X: originX, Y: originY},
			rotation,
			rl.White,
		)
	}
}

func (l *LemonSprite) DoLoad() {
	if l.LOnLoadFunc != nil {
		err := l.GameRef.lState.CallByParam(lua.P{
			Fn:      l.LOnLoadFunc,
			NRet:    0,
			Protect: true,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// Remove sprite from the game's sprite map. This method is not supposed to be
// invoked directly. If a sprite must be deleted, then set its shouldRemove to
// true, and it will be removed automatically before the game's next frame.
func (l *LemonSprite) DoRemove() {
	if l != nil {
		if l.shouldRemove {
			spriteName := l.LTable.RawGetString("name").String()

			// delete its entry based on sprite name
			delete(l.GameRef.spriteDict, spriteName)

			// also delte it from name slice
			for i, name := range l.GameRef.spriteNames {
				if name == spriteName {
					l.GameRef.spriteNames = append(l.GameRef.spriteNames[:i], l.GameRef.spriteNames[i+1:]...)
					break
				}
			}
		}
	}
}

func (l *LemonSprite) DoUpdate(dt float64) {
	if l.LOnUpdateFunc != nil {
		if l.LOnUpdateFunc.Type() != lua.LTNil {
			err := l.GameRef.lState.CallByParam(lua.P{
				Fn:      l.LOnUpdateFunc,
				NRet:    0,
				Protect: true,
			}, lua.LNumber(dt))
			if err != nil {
				panic(err)
			}
		}
	}
}

// }}}

// Lua function interfaces {{{

func (l *LemonSprite) lua_setScript(L *lua.LState) int {
	Argcheck(SpriteFnSetScript, 2, L)
	this := L.Get(1)
	scriptPath := path.Join(l.GameRef.RootDir, L.ToString(2))
	scriptState := lua.NewState()

	// We add a script-wide global variable (relative to the script) called `this`
	// which refers to the script holder (this sprite)
	scriptState.SetGlobal("this", this)

	// This script still need to access L variable. Thus, we link "global L"
	// to the "script-wide global L"
	scriptState.SetGlobal("L", L.GetGlobal("L"))

	if l.GameRef.BuildMode == "standalone" {
		memoryFileName := strings.TrimPrefix(scriptPath, l.GameRef.RootDir)
		buf := l.GameRef.GameAttachment[memoryFileName]
		err := scriptState.DoString(string(buf))
		if err != nil {
			panic(err)
		}
	} else {
		err := scriptState.DoFile(scriptPath)
		if err != nil {
			panic(err)
		}
	}

	// Extract on_load and on_update from the script state
	l.LOnLoadFunc = scriptState.GetGlobal("on_load")
	l.LOnUpdateFunc = scriptState.GetGlobal("on_update")

	// Immediately call on_load once it is set
	if l.LOnLoadFunc.Type() != lua.LTNil {
		err := l.GameRef.lState.CallByParam(lua.P{
			Fn:      l.LOnLoadFunc,
			NRet:    0,
			Protect: true,
		})
		if err != nil {
			panic(err)
		}
	}

	return 0
}

func (l *LemonSprite) lua_setTextureFromFile(L *lua.LState) int {
	Argcheck(SpriteFnSetTextureFromFile, 2, L)
	spriteName := l.LTable.RawGetString("name").String()
	fileName := path.Join(l.GameRef.RootDir, L.ToString(2))

	l.GameRef.spriteDict[spriteName].Texture = NewTexture(fileName, l.GameRef)

	tw := lua.LNumber(l.GameRef.spriteDict[spriteName].Texture.TextureData.Width)
	th := lua.LNumber(l.GameRef.spriteDict[spriteName].Texture.TextureData.Height)
	l.LTable.RawSetString("width", tw)
	l.LTable.RawSetString("height", th)
	l.LTable.RawSetString("origin_x", tw/2)
	l.LTable.RawSetString("origin_y", th/2)
	l.LTable.RawSetString("frame_width", tw)
	l.LTable.RawSetString("frame_height", th)
	return 0
}

func (l *LemonSprite) lua_setTexture(L *lua.LState) int {
	Argcheck(SpriteFnSetTexture, 2, L)
	spriteTable := L.ToTable(1)
	texLTable := L.ToTable(2)
	texID := int(texLTable.RawGetString("id").(lua.LNumber))
	spriteName := spriteTable.RawGetString("name").(lua.LString).String()

	l.GameRef.spriteDict[string(spriteName)].Texture = l.GameRef.textureDict[texID]

	tw := lua.LNumber(l.GameRef.spriteDict[spriteName].Texture.TextureData.Width)
	th := lua.LNumber(l.GameRef.spriteDict[spriteName].Texture.TextureData.Height)
	l.LTable.RawSetString("width", tw)
	l.LTable.RawSetString("height", th)
	l.LTable.RawSetString("origin_x", tw/2)
	l.LTable.RawSetString("origin_y", th/2)
	l.LTable.RawSetString("frame_width", tw)
	l.LTable.RawSetString("frame_height", th)
	return 0
}

func (l *LemonSprite) lua_remove(L *lua.LState) int {
	Argcheck(SpriteFnRemove, 0, L)
	l.shouldRemove = true
	return 0
}

// }}}
