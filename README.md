# Cogent ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/AbyssExplorer/cogent) ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/AbyssExplorer/Cogent/total) ![GitHub Release Date](https://img.shields.io/github/release-date/AbyssExplorer/Cogent)


A simple CLI for generating client credentials tokens for AWS Cognito App Clients

## Download
> If Go is not installed
1. Download executable from [release page](https://github.com/AbyssExplorer/Cogent/releases/), copy it into a folder.
2. Add the folder path in environment variable for easy access.
> Using go get
1. To install cogent, run `go install github.com/AbyssExplorer/Cogent`
   
## Usage
1. Run `cogent` (uses `us-east-1` as default region) or `cogent --region {{region}}`

---
> Note: You need to configure aws credentials (profile: default) so cogent can read and use them
