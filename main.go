package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"

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

	/*b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))*/

	walkTerrain(c, 250)
}

func walkTerrain(c map[string]color.Color, size int) {
	w, err := world.New(worldDirPath)
	if err != nil {
		log.Fatalf("getting world: %s", err)
	}

	yheight := 100

	//
	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: size, Y: size}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

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
					key := strings.Replace(b.ID, "minecraft:", "", 1)
					col, ok := c[key]
					if ok {
						img.Set(x, z, col)
					} else {
						//fmt.Println("not found:", key)
					}

					//fmt.Println(x, y, z, " - ", b.ID, col)

					break yloop
				}

			}
		}
	}

	f, _ := os.Create("test.png")
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
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
