![logo](README_DATA/logo.png) MusicPly
============

*A simple web-ui to show local music playlists and play audio files.*  
*Also it looks like Win98, for whatever reason ¯\\\_(ツ)\_/¯*

## Screenshot

![screenshot](README_DATA/screenshot.png)

## Installation (via git clone && make)

```
$> git clone https://github.com/Mikescher/musicply
$> make build-quick
$> SOURCE=".." ./build/musicply
```

## Installation (via docker)

```
$> docker pull mikescher/musicply
$> docker run                   \
          --volume "..."        \
          --publish "8000:8000" \
          --env "SOURCE=..."    \
          --name=musicply       \
          "mikescher/musicply:latest"
```

## Usage / Configuration

MusicPly is configured via environment variables.
The most important config is `SOURCE`, which specified the source directories.  
You can either
 - supply a single (json5) array in the `SOURCE` environment variable, which contains a list of all sources in the schema:   
   `{name: "...", path: "...", recursive: ...}` 
 - supplying a single path to a json5 file in the `SOURCE` environment variable
 - supply multiple `SOURCE` variables by appending indizes (aka `SOURCE_01`, `SOURCE_02`, `SOURCE_03`, ...) to the environment variable, where each variable contains a single source-object

An example config file (`config.json`) could look like:
```
[
  {name: "Back in Black",    path: "/data/ACDC/BackInBlack", recursive: true}, // with recursive:true we also iterate through all subfolders.
  {name: "Hotel California", path: "/data/Eagles/HotelCalifornia"},            // if recursive is not specified, teh default value ist false

  // values are shown in the here specified order in the web UI 
]
```

This config can be used by supplying the filepath:

```
$> SOURCE="/config.json" ./build/musicply

# or 

$> docker run                                             \
          --volume "/home/user/music:/data:ro"            \
          --volume "$(pwd)/config.json:/config.json:ro"   \
          --publish "8000:8000"                           \
          --env "SOURCE=/config.json"                     \
          "mikescher/musicply:latest"
```

Or you can provide the json5 directly in the environment variable:

```
$> SOURCE="$(cat /config.json)" ./build/musicply

# or 

$> docker run                                             \
          --volume "/home/user/music:/data:ro"            \
          --publish "8000:8000"                           \
          --env "SOURCE=$(cat /config.json)"              \
          "mikescher/musicply:latest"
```

Or you can provide the sources individually:

```
$> SOURCE_1='{name: "Back in Black",    path: "/data/ACDC/BackInBlack", recursive: true}' \
   SOURCE_2='{name: "Hotel California", path: "/data/Eagles/HotelCalifornia"}' \
   ./build/musicply

# or 

$> docker run                                                                                          \
          --volume "/home/user/music:/data:ro"                                                         \
          --publish "8000:8000"                                                                        \
          --env 'SOURCE_1={name: "Back in Black",    path: "/data/ACDC/BackInBlack", recursive: true}' \
          --env 'SOURCE_2={name: "Hotel California", path: "/data/Eagles/HotelCalifornia"}'            \
          "mikescher/musicply:latest"
```

## Additional Configuration

The following environment variables can also be used to configure the application:

- `SERVER_IP` The interface the webserver binds to (default: 0.0.0.0)
- `SERVER_PORT` The  webserver port (default: 8000)
- `CORS` Enable CORS headers (default: true)
- `LOGLEVEL`, `CUSTOM_404`, `GIN_DEBUG`, `RETURN_RAW_ERRORS` Enable more logoutput (default: 'WARN', false, false, false)

You can also show additional buttons under the playlist-control by supplying `/FOOTERLINK_[0-9]/` env variables.  
The variables must contain 3, semicolon-seperated values: `${icon-path};${Tooltip};${Link}`