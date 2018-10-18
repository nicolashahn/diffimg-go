# diffimg-go
Image differentiation tool similar to my python module and command line tool `diffimg`: https://github.com/nicolashahn/diffimg

You may find it useful for testing image processing applications/services.

## CLI tool installation
1. Clone the repo: `git clone https://github.com/nicolashahn/diffimg-go`
2. `cd` into it: `cd diffimg-go`
3. Add it to your `$GOPATH/bin`: `go install`

## Usage

```
$ diffimg-go [-generate] [-ratio] [-invertalpha] IMAGE1 IMAGE2
```

`IMAGE1` and `IMAGE2` are image files. They must be the same size.

`-generate` creates a diff image at `diff.png`.

`-ratio` returns a `float` ratio instead of the sentence `Images differ by X.XX%`

`-invertalpha` inverts the alpha channel for the generated diff image. If both images are fully opaque (all the alpha channel values for all pixels in both images are the maximum value) then a simple diff would produce a fully transparent image. Use this flag if you do not want that. It does not affect the ratio, so it does nothing if the `-generate` flag is not used.
