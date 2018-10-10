package main

import (
  "image"
  "fmt"
  "os"
)
import _ "image/png"

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}

func main () {
  
  f1, err := os.Open("images/mario-circle-cs.png")
  checkErr(err)
  defer f1.Close()
  im1, _, err := image.Decode(f1)
  checkErr(err)

  f2, err := os.Open("images/mario-circle-node.png")
  checkErr(err)
  defer f2.Close()
  im2, _, err := image.Decode(f2)
  checkErr(err)

  bounds1 := im1.Bounds()
  bounds2 := im2.Bounds()

  fmt.Println(bounds1, bounds2)
  if bounds1 != bounds2 {
    panic("image dimensions are different")
  }


}
