# find-unused

Find unused translation strings in a big codebase. Given a JSON translation file and a source directory, run like this:

```sh
find-unused lang.json ./src
```

This will print keys that look like they're unused. Be sure to check (and use version control) before deleting anything.
