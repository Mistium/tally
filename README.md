# tally

A simple, easy to use and powerful command line tool for tracking lines of code.

## Installation

```bash
cd src
go build
cp tally /usr/local/bin
```

## Usage

```bash
tally [options]
```

### Options

| Option | Description |
| ------ | ----------- |
| `-c`   | Count lines of code in the current directory |
| `-d`   | Count lines of code in the current directory and all subdirectories |
| `-a`   | Count lines of code in the current directory and all subdirectories, including all files |
| `-h`   | Show help |

## Example

### Count lines of code in the current directory

```txt
$ tally -c       
| Extension | Lines of Code | Percentage |
| --------- | ------------- | ---------- |
| md        |            63 |      53.4% |
| gitignore |            34 |      28.8% |
| (no ext)  |            21 |      17.8% |

Total lines of code: 118
```

### Count lines of code in all subdirectories of a specific file type

```txt
$ tally -d go, md
| Extension | Lines of Code | Percentage |
| --------- | ------------- | ---------- |
| go        |           218 |      77.9% |
| md        |            62 |      22.1% |

Total lines of code: 280
```

### Count lines of code in the current directory and all subdirectories, including all files

```txt
$ tally -a
| Extension | Lines of Code | Percentage |
| --------- | ------------- | ---------- |
| go        |           218 |      69.0% |
| md        |            42 |      13.3% |
| gitignore |            32 |      10.1% |
| (no ext)  |            21 |       6.6% |
| mod       |             3 |       0.9% |
```

### Get json output

```txt
tally -c --json
{
  "files": {
    "": 21,
    ".gitignore": 34,
    ".md": 69,
  },
  "total": 124,
}
```

## License

MIT