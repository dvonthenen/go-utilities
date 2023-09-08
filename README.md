# go-utilities

This is just a collection of utilities I have created and need to put under source control...

- `cmd/cli/file-distribute`: A tool to copy a folder containing MP3s onto a USB drive so that Honda factory radios can access/play all of the files. Honda radios have a limitation to only have 255 files/folders at any level in the directory structure. So this will create 255 folders at the root of the drive and distribute them into those 255 folders without changing or modifiing the source. That means you need to have free space equivalent to the size of the source folder.
