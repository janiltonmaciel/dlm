Installation
The easiest way to install the latest binaries for Linux and Mac is to run this in a shell:

curl -sSf https://moncho.github.io/dry/dryup.sh | sudo sh
sudo chmod 755 /usr/local/bin/dry
Binaries
If you dont like to curl | sh, binaries are provided.

darwin 386 / amd64
freebsd 386 / amd64
linux 386 / amd64
windows 386 / amd64
Mac OS X / Homebrew
If you're on OS X and want to use homebrew:

brew tap moncho/dry
brew install dry

Copyright and license
Code released under the MIT license. See LICENSE for the full license text.

Credits
Built on top of:

termbox
termui
Docker
Docker CLI

# dockerfile-gen

Generator Dockerfile

## Installation

#### Binaries

- **darwin (macOS)** [amd64](https://github.com/janiltonmaciel/dockerfile-gen/releases/download/1.10.0/dockerfile-gen_1.10.0_macOS_amd64.tar.gz)
- **linux** [386](https://github.com/janiltonmaciel/dockerfile-gen/releases/download/1.10.0/dockerfile-gen_1.10.0_linux_386.tar.gz) / [amd64](https://github.com/janiltonmaciel/dockerfile-gen/releases/download/1.10.0/dockerfile-gen_1.10.0_linux_amd64.tar.gz)
- **windows** [386](https://github.com/janiltonmaciel/dockerfile-gen/releases/download/1.10.0/dockerfile-gen_1.10.0_windows_386.zip) / [amd64](https://github.com/janiltonmaciel/dockerfile-gen/releases/download/1.10.0/dockerfile-gen_1.10.0_windows_amd64.zip)

#### Via Homebrew (macOS)
```bash
$ brew tap janiltonmaciel/homebrew-tap
$ brew install dockerfile-gen
```

#### Via Go

```bash
$ go get github.com/janiltonmaciel/dockerfile-gen
```

#### Running with Docker

```bash
$ docker run -it --rm \
    -v $(pwd):/app \
    janilton/dockerfile-gen
```

## Usage
Creating Dockerfile

![](https://github.com/janiltonmaciel/dockerfile-gen/blob/master/assets/img/dc-gen-create.gif)

---
Building docker image


![](https://github.com/janiltonmaciel/dockerfile-gen/blob/master/assets/img/dc-gen-build.gif)

---
Running docker image


![](https://github.com/janiltonmaciel/dockerfile-gen/blob/master/assets/img/dc-gen-run.gif)
