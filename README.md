
<h1 align="center">
  <br>
  <a href="http://github.com/cakehappens/frosting"><img src="assets/cupcake.png" alt="playing card" width="200px" /></a>
  <br>
  Frosting
  <br>
</h1>

<h4 align="center">Enhance the operational tasks of your application with a little <i>frosting</i> 🧁</h4>

<p align="center">
  <a href="https://saythanks.io/to/ghostsquad">
      <img src="https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg">
  </a>
  <a href="https://www.paypal.me/WMcNamee">
    <img src="https://img.shields.io/badge/$-donate-ff69b4.svg?maxAge=2592000&amp;style=flat">
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

## Introduction

`frost` is a workflow tool. Inspired by [Make](https://www.gnu.org/software/make/), [Mage](https://magefile.org/), [Task](https://taskfile.dev/) and others, but with some important differences.

| Feature              | Frost | Make | Mage | Go-Task |
|----------------------|-------|------|------|---------|
| Fast                 | ✔️     | ✔️    | ✔️    | ✔️       |
| Plugins              | ✔️     |      | ✔️    |         |
| Bash Support         | ✔️     | ✔️    | ✔️    | ✔️       |
| Target Flags         | ✔️     | ✔️§   | ✔️§   | ✔️§§     |
| Autocomplete         | ✔️     | ✔️    |      | ✔️       |
| Parallelism          | ✔️     |      | ✔️†   |         |
| Imports (Namespaced) | ✔️     | ✔️    | ✔️‡   | ✔️‡      |
| *File                | Go    | Make | Go   | Yaml    |

§ Target behavior can only be modified with environment variables

§§ Target behavior can be modified with environment and can also can be templated with go templating prior to evaluation and execution.

† Doesn't support limiting of parallelism ([#273](https://github.com/magefile/mage/pull/273))

‡ Mage doesn't currently support nested namespaces ([#152](https://github.com/magefile/mage/issues/152)), and Task support for imports is still experimental

---

Make is fast, but for any sort of complex logic, I found myself writing bash scripts, and I agree with Mage in regards to make/bash:

> Makefiles are hard to read and hard to write. Mostly because makefiles are essentially fancy bash scripts with significant white space and additional make-related syntax. Go is superior to bash for any non-trivial task involving branching, looping, anything that’s not just straight line execution of commands.

Mage makes heavy use of code generation in order to create the resulting binary, and I found that to be a high barrier to entry for contributing to the project.

I've been itching to write my own CLI for awhile now, and I think I finally have a reason to.

## Install

```go

```

## How To Use

```bash

```

## Credits

- spf13/cobra
- spf13/viper
- hashicorp/go-plugins
- hashicorp/go-versions

## Support

<a href="https://www.buymeacoffee.com/50onA1pjc" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important; width: auto !important;" /></a>

<a href="https://www.paypal.me/WMcNamee" target="_blank"><img src="https://user-images.githubusercontent.com/903488/58933498-4cde9380-871c-11e9-88a5-455ed1d14380.png" alt="Paypal Donate" style="height: auto !important;width: auto !important;" /></a>

_Icons made by [Freepik](https://www.freepik.com) from [www.flaticon.com](http://www.flaticon.com) is licensed by [CC 3.0 BY](http://creativecommons.org/licenses/by/3.0/)_

## Related

## License