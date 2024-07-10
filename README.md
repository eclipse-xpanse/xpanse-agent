# xpanse-agent
xpanse agent to poll and execute config change requests.

## How to use

### Build from source

```shell
make build
```

### Start the xpanse-agent

```shell
./xpanse-agent -h

                                                              _
__  ___ __   __ _ _ __  ___  ___        __ _  __ _  ___ _ __ | |_
\ \/ / '_ \ / _' | '_ \/ __|/ _ \_____ / _' |/ _' |/ _ \ '_ \| __|
 >  <| |_) | (_| | | | \__ \  __/_____| (_| | (_| |  __/ | | | |_
/_/\_\ .__/ \__,_|_| |_|___/\___|      \__,_|\__, |\___|_| |_|\__|
     |_|                                     |___/

Usage:
  xpanse-agent [flags]
  xpanse-agent [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  start       To start xpanse-agent
  version     To view the version of the xpanse-agent

Flags:
  -h, --help   help for xpanse-agent

Use "xpanse-agent [command] --help" for more information about a command.
```

## Environment Variables

All arguments can also be read from environment variables. 
Variables in environment variables must be all in uppercase and 
must be prefixed with `XPANSE_AGENT_`.