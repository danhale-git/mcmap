package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"log"

	"github.com/danhale-git/mcmap/colors"

	"github.com/danhale-git/mine/world"
)

const worldDirPath = `C:\Users\danha\AppData\Local\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds\VsgSYaaGAAA=`
const texturePath = `C:\Users\danha\go\src\github.com\danhale-git\mcmap\textures\blocks`

func main() {
	c, err := colors.MapColors(texturePath)
	if err != nil {
		log.Fatalf("mapping texture pngs to block colors: %s", err)
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func walkTerrain(c map[string]color.Color) {
	w, err := world.New(worldDirPath)
	if err != nil {
		log.Fatalf("getting world: %s", err)
	}

	size := 32
	yheight := 100

	for x := 0; x < size; x++ {
		for z := 0; z < size; z++ {
		yloop:
			for y := yheight; y >= 0; y-- {
				b, err := w.GetBlock(x, y, z, 0)
				if err != nil {
					if errors.Is(err, &world.SubChunkNotSavedError{}) {
						continue
					} else {
						log.Fatal(err)
					}
				}

				if b.ID != "minecraft:air" {
					fmt.Println(x, y, z, " - ", b.ID)
					break yloop
				}
			}
		}
	}
}

/*func getImageNames() {
	dir, err := os.Stat(texturePath)
	if err != nil {
		log.Fatal(err)
	}

	if !dir.IsDir() {
		log.Fatalf("%s is not a directory", texturePath)
	}

	files, err := os.ReadDir(texturePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}*/
