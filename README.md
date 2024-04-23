# find-unused

Find unused translation strings in a big codebase. Given a JSON translation file and a source directory.

## Install

```sh
# download and install from source
git clone https://github.com/eyeseast/find-unused && cd find-unused
go install .
```

# Usage

```sh
find-unused lang.json ./src
```

This will print keys that look like they're unused. Be sure to check (and use version control) before deleting anything.
