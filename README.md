# helm-run :rocket:

A Helm plugin that executes scripts defined in Chart.yaml similar to `npm run`. 
This lets you define and run repetitive tasks like linting, testing, building, or validations per chart.

## Features

- Define scripts centrally in Chart.yaml and run them via `helm run`
- Different script runners:
  - Shell `sh`
  - Bash `bash`
  - Z-Shell `zsh`
  - Fish `fish`
  - Windows Powershell `powershell`
  - Powershell `pwsh`
  - Windows Commandline `cmd`
  
> [!NOTE]
> By default, the plugin uses `$SHELL` which is the parent shell of the `helm run` command on linux and macOS.
> On Windows, it uses `powershell` by default.

> [!WARNING]
> Windows Powershell `powershell` and Powershell `pwsh` are different!
> For more information, see [this article](https://learn.microsoft.com/en-us/powershell/scripting/whats-new/differences-from-windows-powershell).

## Requirements

- Helm v3
- For building from source: Go 1.23+

> [!NOTE]
> The plugin contains pre-built binaries for linux, macOS, and Windows with the amd64 and arm64 architectures.
> If you are on a different platform, the plugin will try to build the binary from source.
> This requires Go to be installed!

## Installation

### As a Helm plugin (recommended)

In the project directory:
```
helm plugin install https://github.com/benni-tec/helm-run
```

### Build from source (optional)
```sh
go build main -o helm-run
```

You can then either call `helm-run` from the local binary or add it to your `PATH`.

## Usage

### Definition :page_facing_up:

You can define scripts in your Chart.yaml either as a simple script, which is just a string, or a runner-dependent script, which is an object with a field for each runner.

> [!NOTE]
> If you define a runner-dependent script and `helm-run` can not find an exact match for the runner, it will try to guess the best matching one.
> Specifically, it will use the `sh` script for all shells (`bash`, `zsh`, `fish`) and `powershell` for `pwsh`.

```yaml
apiVersion: v2
name: example
version: 0.0.1

scripts:
  simple-script: echo "Hello world!"
  runner-dependent-script:
    sh: echo "Hello world, from shell!"
    bash: echo "Hello world, from bash!"
    zsh: echo "Hello world, from zsh!"
    fish: echo "Hello world, from fish!"
    
    powershell: echo "Hello world, from powershell!"
    pwsh: Write-Output "Hello world, from pwsh!"
    cmd: echo Hello world, from cmd!
```

> [!TIP]
> The plugin copies the script to a temporary file and then executes it.
> This means that you can use multiline scripts.

### Run a script :rocket:

Scripts can then be run via `helm run`:

```sh
helm run [--verbose] [--runner <sh|bash|zsh|fish|powershell|pwsh|cmd>] [--runner-args <args>] <chart> <script>
``` 

Arguments:
- `<chart>`: Path to the chart (e.g., `.` for the current directory)
- `<script>`: Name of the script defined in Chart.yaml
- `--runner <sh|bash|zsh|fish|powershell|pwsh|cmd>`: Overrides the default shell runner
- `--runner-args <args>`: Arguments passed to the script runner instead of the script itself
- `--verbose`: Shows additional output/debug information

## Contributing
All contributions are welcome.
- If you have found a bug, please open an issue or a pull request.
- If you want to contribute, e.g. new runners, feel free to open a pull request.