# GoSearcher

## Roadblocks

- Single Instance (https://stackoverflow.com/questions/23162986/restricting-to-single-instance-of-executable-with-golang)
- Global Hotkeys (https://github.com/fyne-io/fyne/issues/2304)

## TODO

- Fix icons not being bundled causing the program to crash

### Requested

- Search field focus on open
- Enter to open item
- Configurable Hotkeys
- Escape to Close

### Hotkeys

- Escape to Close
- Arrows for navigation

## Configuration Location

MacOS:
`$HOME/Library/Application Support/GoSearcher/config.yaml`

Windows:
`%AppData%\GoSearcher\config.yaml`

## Building

MacOS:
`fyne-cross darwin --arch=arm64 -app-id GoSearcher -icon "./icons/search-512.png" -name GoSearcher`

Windows:
`fyne-cross windows --arch=386 -icon "./icons/search-512.png" -name GoSearcher`
