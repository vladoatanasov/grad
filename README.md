# GRAD (GitHub Release Asset Downloader)

This is a small command line utiliy which downloads release artefacts via the [GitHub Release API](https://developer.github.com/v3/repos/releases/)

## Usage
```
-artefact string
  	Artefacts to download, comma separated
-d	Debug logging
-release string
  	Release tag
-repo string
  	user/repo
-token string
  	Git personal access token, to access private repos
-v	Print current version
```

## Example
```./grad -repo=atom/atom -release=latest -artefact=AtomSetup.msi,AtomSetup.exe``` will download the artefacts from the latest release, you can also use a specific release ```./grad -repo=atom/atom -release=v1.8.0 -artefact=AtomSetup.msi```

## Note
This utility was hacked together pretty quickly, so I haven't done a lot of testing. Please use caution if you decide to use it in production. PR's are welcome
