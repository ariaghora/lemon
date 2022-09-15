package pkg

import (
	"path"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
	lua "github.com/yuin/gopher-lua"
)

type TextureInfo struct {
	TextureID   int
	LTable      *lua.LTable
	TextureData rl.Texture2D
}

func NewTexture(fileName string, gameRef *LemonGame) *TextureInfo {
	var o *TextureInfo
	if gameRef.BuildMode == "standalone" {
		memoryFileName := strings.TrimPrefix(fileName, gameRef.RootDir)
		buf := gameRef.GameAttachment[memoryFileName]
		imageSize := len(buf)
		imageData := rl.LoadImageFromMemory(path.Ext(memoryFileName), buf, int32(imageSize))

		o = &TextureInfo{
			TextureID:   gameRef.textureCount,
			TextureData: rl.LoadTextureFromImage(imageData),
		}
	} else {
		o = &TextureInfo{
			TextureID:   gameRef.textureCount,
			TextureData: rl.LoadTexture(fileName),
		}
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
