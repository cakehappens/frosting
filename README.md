
<h1 align="center">
  <br>
  <a href="http://github.com/cakehappens/frosting"><img src="./assets/cupcake.png" alt="playing card" width="200px" /></a>
  <br>
  Frosting
  <br>
</h1>

<h4 align="center">Enhance the operational tasks of your application with a little <i>frosting</i> 🧁</h4>

<p align="center">
  <a href="https://pkg.go.dev/github.com/cakehappens/frosting">
    <img src="https://img.shields.io/badge/godoc-reference-5272B4.svg">
  </a>
  <!-- <a href="https://goreportcard.com/badge/github.com/cakehappens/frosting">
    <img src="https://goreportcard.com/report/github.com/cakehappens/frosting">
  </a> -->
  <a href="https://saythanks.io/to/ghostsquad">
    <img src="https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg">
  </a>
  <a href="buymeacoff.ee/50onA1pjc">
    <img src="https://img.shields.io/badge/buymeacoffee-%24-orange">
  </a>
</p>

<p align="center">
  <a href="#introduction">Introduction</a> •
  <a href="#install">Install</a> •
  <a href="#how-to-use">How To Use</a> •
  <a href="#credits">Credits</a> •
  <a href="#support">Support</a> •
  <a href="#related">Related</a> •
  <a href="#license">License</a>
</p>

## 👋 Introduction

`frosting` is library that lets you quickly and easily create a CLI for your code repositories (like how a `Makefile` enables you to run `make build`). Inspired by [Make][make], [Mage][mage], [Task][taskfile] and others.

## 🎯 Features

| Feature                          | Frost | Make | Mage | Go-Task |
|----------------------------------|-------|------|------|---------|
| *File                            | Go    | Make | Go   | Yaml    |
| Bash Support                     | 🧁    | 🐮   | 🧙  | 🐹     |
| Target-Specific Vars             | 🧁    | 🐮   | 🧙  | 🐹     |
| Namespaces                       | 🧁    | 🐮   |     | 🐹      |
| Imported Targets                 | 🧁    | 🐮   | 🧙  | 🐹‡    |
| Bash/Zsh Autocomplete            | 🧁    | 🐮   |     | 🐹      |
| Parallelism                      | 🧁    |      | 🧙† |         |
| No Custom DSL to Learn           | 🧁    |      | 🧙  |         |
| Target Args					   | 🧁    |      |      |         |
| Target Flags                     | 🧁    |      |      |         |
| Target-Specific Help             | 🧁    |      |      |         |
| [Color Support][color]           | 🧁    |      |      |         |
| Doc Generation                   | 🧁    |      |      |         |
| [Interative Terminal UI][tview]  | 🧁    |      |      |         |
| [go-prompt][gpt] Integration     | 🧁    |      |      |         |
| [Progress Bars][pb]              | 🧁    |      |      |         |
| [Spinners][spin]                 | 🧁    |      |      |         |

† Yes, but with some limitations ([#273](https://github.com/magefile/mage/pull/273))

‡ Task support for imports is still experimental

## 💡 Philosophy

Mage puts it nicely in regards to make/bash:

> Makefiles are hard to read and hard to write. Mostly because makefiles are essentially fancy bash scripts with significant white space and additional make-related syntax. Go is superior to bash for any non-trivial task involving branching, looping, anything that’s not just straight line execution of commands.

Mage makes heavy use of code generation in order to create the resulting binary, and I found that to be a high barrier to entry for contributing to the project.

I've been itching to write my own CLI for awhile now, and I think I finally landed on a great use case.

## ⚡️ Quickstart

```go
package main

import (
	"context"
	"fmt"
	"github.com/cakehappens/frosting"
	"github.com/cakehappens/frosting/ingredient"
)

func NewBuildIngredient() *ingredient.Ingredient {
	return ingredient.MustNew(
		"build",
		func(ctx context.Context) error {
			fmt.Println("Building...")
			return nil
		},
		ingredient.WithDependencies(NewTestIngredient),
	)
}

func NewTestIngredient() *ingredient.Ingredient {
	return ingredient.MustNew(
		"test",
		func(ctx context.Context) error {
			fmt.Println("Testing...")
			return nil
		},
	)
}

func main() {
	f := frosting.New("frost")
	f.MustAddIngredientGroups(
		ingredient.MustNewGroup(
			"",
			"Main Stuff:",
			ingredient.Includes(
				NewBuildIngredient,
				NewTestIngredient,
			),
		),
	)

	f.Execute("foo")
}

```

```bash
go build -o frost
frost # prints help
frost build # runs build ingredient
frost build --help # prints build-target help
```

## 🎓 Docs

see [docs](docs) for more info, or look at the [godocs](https://pkg.go.dev/github.com/cakehappens/frosting)

## 👀 Examples

see [examples](examples)

## 🌟 Contribute

I'll definitely get some templates/guidelines setup soon...

## 🤗 Support

<a href="buymeacoff.ee/50onA1pjc">
    <img src="https://img.shields.io/badge/buymeacoffee-%24-orange">
</a>

## 📖 Reading

- https://dave.cheney.net/2017/06/11/go-without-package-scoped-variables
- https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
- https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
- https://dave.cheney.net/tag/logging
- https://stackoverflow.com/questions/2214575/passing-arguments-to-make-run

## 💕 Related

- [Mage][mage]
- [Taskfile][taskfile]
- [Make][make]

## 📜 Credits

- [spf13/cobra][cobra]
- [spf13/viper][viper]
- [tivo/tview][tview]
- [theckman/yackspin][spin]
- [c-bata/go-prompt][gpt]
- [gosuri/uiprogress][pb]
- [fatih/color][color]
- Icons made by [Freepik](https://www.freepik.com) from [www.flaticon.com](http://www.flaticon.com) is licensed by [CC 3.0 BY](http://creativecommons.org/licenses/by/3.0/)

## ⚖️ License

[Apache License, Version 2.0, http://www.apache.org/licenses/](LICENSE)

[make]: https://www.gnu.org/software/make/
[taskfile]: https://taskfile.dev/
[mage]: https://magefile.org/
[cobra]: https://github.com/spf13/cobra
[viper]: https://github.com/spf13/viper
[gpt]: https://github.com/c-bata/go-prompt
[spin]: https://github.com/theckman/yacspin
[pb]: https://github.com/gosuri/uiprogress
[color]: https://github.com/fatih/color
[tview]: https://github.com/rivo/tview
