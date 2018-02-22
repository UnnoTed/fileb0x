# Changelog
All notable changes to this project will be documented in this file.

To update simply run:
```bash
go get -u github.com/UnnoTed/fileb0x
```

## 2018-02-22
### Added
- Avoid rewriting the main b0x file by checking a MD5 hash of the (file's modification time + cfg).
- Avoid rewriting unchanged files by comparing the Timestamp of the b0x's file and the file's modification time.
- Config option `lcf` which when enabled along with `spread` **l**ogs the list of **c**hanged **f**iles to the console.
- Message to inform that no file or cfg changes have been detecTed (not an error).
### Changed
- Config option `clean` to only remove unused b0x files instead of everything.
