package pkg

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	lua "github.com/yuin/gopher-lua"
)

type TextureInfo struct {
	TextureID   int
	LTable      *lua.LTable
	TextureData rl.Texture2D
}

func NewTexture(fileName string, gameRef *LemonGame) *TextureInfo {
	o := &TextureInfo{
		TextureID:   gameRef.textureCount,
		TextureData: rl.LoadTexture(fileName),
	}
	gameRef.textureDict[gameRef.textureCount] = o // keep track

	// calling prototype function to generate initial sprite table structure
	proto_new_texture_fn := gameRef.lState.GetGlobal("__proto_new_texture").(*lua.LFunction)
	gameRef.lState.CallByParam(
		lua.P{
			Fn:      proto_new_texture_fn,
			NRet:    1,
			Protect: true,
		},
		lua.LNumber(gameRef.textureCount),
	)
	table := gameRef.lState.Get(-1).(*lua.LTable)
	o.LTable = table

	gameRef.textureCount += 1
	return o
}
