# Release process

Homebrew formula: https://github.com/ayakovlenko/homebrew-tools/blob/master/zit.rb

Run `release.sh` to tag a new release and push it to the main branch:

```sh
`./scripts/release.sh`
```

Run `homebrew.sh` to calculate SHA256 sum of the new release and update the
Homebrew formula with it.

```sh
`./scripts/homebrew.sh`
```
