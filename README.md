# SROC - a very simple CORS proxy

used to allow [aligator.dev](https://aligator.dev) to load stl files for 
[GoSlice](https://github.com/aligator/goslice) running as Webassembly
from any external URL.

## Usage
If the proxy runs at `https://cors.aligator.dev` then you can call it with 
`https://cors.aligator.dev/?target=https://cdn.thingiverse.com/assets/7d/fc/6e/33/fe/3DBenchy.stl`
and it will just proxy that request and setting the cors headers as needed.

See `--help` for all possible options.

## Limitation
This proxy is really simple and just built for the aligator.dev webpage. It may not handle all cases.
Also it is currently only allowed to set one possible origin. All other than the set one will be rejected.
