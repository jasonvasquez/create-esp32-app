# Create ESP32 App

Create ESP32 apps using the [ESP-IDF](https://github.com/espressif/esp-idf) with easier dev environment configuration.
All dependencies and build tools live within an easily updateable Docker image, resulting in minimal impact on workstation
setup.

Create ESP32 App has only been tested on Linux so far, other platforms are welcome.

## Getting Started

### Installation

* Install Docker, if you don't have it already
* Install `direnv` if you wish, to make accessing the resulting `esp32` script a little easier
* Install createESP32App  (TBD distribution method, or build it from this repo :))

### Creating an App

To create a new app, run:

```sh
createESP32App -appName myESP32App
cd myESP32App
esp32 build
esp32 flash
esp32 monitor
```

It will create a directory called `myESP32App` inside the current folder.<br>
Inside that directory, it will generate the initial project structure.

```
myESP32App
├── .envrc
├── .gitignore
├── scripts
│   └── esp32
└── src
    └── Makefile
    └── sdkconfig
    └── main
        └── component.mk
        └── main.c
```
