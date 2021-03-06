Compile DirectX Shaders from Go
===============================

`dxc` is a Go library to compile DirectX shaders and effects. It uses the D3DCompiler_XX.dll DLLs that come with any Windows system.

To install the library, call

	go get -u github.com/gonutz/dxc
	
where the `-u` option will get the latest version from [the Github repository](https://github.com/gonutz/dxc).

See [the godoc](https://godoc.org/github.com/gonutz/dxc) for how to use the library.

Command Line Tool
=================

The `cmd/dxc` folder contains a command line version of this tool.

```
Usage of dxc:
  -E string
    	entrypoint name (default "main")
  -Gch
    	compile as a child effect for FX 4.x targets
  -Gdp
    	disable effect performance mode
  -Gec
    	enable backwards compatibility mode
  -Ges
    	enable strict mode
  -Gfa
    	avoid flow control constructs
  -Gfp
    	prefer flow control constructs
  -Gis
    	force IEEE strictness
  -Gpp
    	force partial precision
  -O int
    	optimization level 0..3 (default 1)
  -Od
    	disable optimizations, only use this for debug builds
  -Op
    	disable preshaders
  -T string
    	target profile, e.g. vs_2_0, ps_4_1, fx_5_0 etc.
  -Vd
    	disable validation
  -WX
    	treat warnings as errors
  -Zi
    	enable debug information in output
  -Zpc
    	pack matrices in column-major order
  -Zpr
    	pack matrices in row-major order
  -res_may_alias
    	assume that UAVs/SRVs may alias for cs_5_0+
```
