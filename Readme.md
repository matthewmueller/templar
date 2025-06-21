# templar

[![Go Reference](https://pkg.go.dev/badge/github.com/matthewmueller/templar.svg)](https://pkg.go.dev/github.com/matthewmueller/templar)

Extensions for [Templ](https://github.com/a-h/templ).

## Extensions

- [scopedcss](./scopedcss): [styled-jsx](https://github.com/vercel/styled-jsx) for Templ. Add support for `<style scoped>` to your Templ templates. See the [testdata](./scopedcss/testdata/) for examples.
- [devmode](./devmode): updates generated `*_templ.go` files to make it easier to support livereload with third-party build systems and file watchers. It does this by not relying on environment variables and only modifying the Go file if it's a non-text change. More details in [this issue](https://github.com/a-h/templ/issues/1108).

## Install

```sh
go get github.com/matthewmueller/templar
```

## Development

First, clone the repo:

```sh
git clone https://github.com/matthewmueller/templar
cd templar
```

Next, install dependencies:

```sh
go mod tidy
```

Finally, try running the tests:

```sh
go test ./...
```

## License

MIT
