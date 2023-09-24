package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/gonutz/dxc"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	target := flag.String("T", "", "target profile, e.g. vs_2_0, ps_4_1, fx_5_0 etc.")
	entryPoint := flag.String("E", "main", "entrypoint name")
	debug := flag.Bool("Zi", false, "enable debug information in output")
	noValidation := flag.Bool("Vd", false, "disable validation")
	noOptimization := flag.Bool("Od", false, "disable optimizations, only use this for debug builds")
	rowMajor := flag.Bool("Zpr", false, "pack matrices in row-major order")
	colMajor := flag.Bool("Zpc", false, "pack matrices in column-major order")
	partialPrec := flag.Bool("Gpp", false, "force partial precision")
	noPreshader := flag.Bool("Op", false, "disable preshaders")
	noFlowCtrl := flag.Bool("Gfa", false, "avoid flow control constructs")
	wantFlowCtrl := flag.Bool("Gfp", false, "prefer flow control constructs")
	strict := flag.Bool("Ges", false, "enable strict mode")
	compatibility := flag.Bool("Gec", false, "enable backwards compatibility mode")
	ieeeStrict := flag.Bool("Gis", false, "force IEEE strictness")
	optLevel := flag.Int("O", 1, "optimization level 0..3")
	warnToErr := flag.Bool("WX", false, "treat warnings as errors")
	resMayAlias := flag.Bool("res_may_alias", false, "assume that UAVs/SRVs may alias for cs_5_0+")

	childEffect := flag.Bool("Gch", false, "compile as a child effect for FX 4.x targets")
	noEffectPerf := flag.Bool("Gdp", false, "disable effect performance mode")

	flag.Parse()

	if len(flag.Args()) > 1 {
		msg := "Invalid arguments: " + strings.Join(flag.Args(), " ")
		os.Stderr.WriteString(msg)
		return errors.New(msg)
	}

	var compileFlags uint
	compileFlag := func(b bool, flag uint) {
		if b {
			compileFlags |= flag
		}
	}
	compileFlag(*debug, dxc.DEBUG)
	compileFlag(*noValidation, dxc.SKIP_VALIDATION)
	compileFlag(*noOptimization, dxc.SKIP_OPTIMIZATION)
	compileFlag(*rowMajor, dxc.PACK_MATRIX_ROW_MAJOR)
	compileFlag(*colMajor, dxc.PACK_MATRIX_COLUMN_MAJOR)
	compileFlag(*partialPrec, dxc.PARTIAL_PRECISION)
	compileFlag(*noPreshader, dxc.NO_PRESHADER)
	compileFlag(*noFlowCtrl, dxc.AVOID_FLOW_CONTROL)
	compileFlag(*wantFlowCtrl, dxc.PREFER_FLOW_CONTROL)
	compileFlag(*strict, dxc.ENABLE_STRICTNESS)
	compileFlag(*compatibility, dxc.ENABLE_BACKWARDS_COMPATIBILITY)
	compileFlag(*ieeeStrict, dxc.IEEE_STRICTNESS)
	compileFlag(*optLevel == 0, dxc.OPTIMIZATION_LEVEL0)
	compileFlag(*optLevel == 1, dxc.OPTIMIZATION_LEVEL1)
	compileFlag(*optLevel == 2, dxc.OPTIMIZATION_LEVEL2)
	compileFlag(*optLevel == 3, dxc.OPTIMIZATION_LEVEL3)
	compileFlag(*warnToErr, dxc.WARNINGS_ARE_ERRORS)
	compileFlag(*resMayAlias, dxc.RESOURCES_MAY_ALIAS)

	var effectsFlags uint
	if *noEffectPerf {
		effectsFlags |= dxc.EFFECT_ALLOW_SLOW_OPS
	}
	if *childEffect {
		effectsFlags |= dxc.EFFECT_CHILD_EFFECT
	}

	var code bytes.Buffer
	if _, err := io.Copy(&code, os.Stdin); err != nil {
		os.Stderr.WriteString(err.Error())
		return err
	}

	output, err := dxc.Compile(
		code.Bytes(),
		*entryPoint,
		*target,
		compileFlags,
		effectsFlags,
	)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return err
	}

	if _, err := os.Stdout.Write(output); err != nil {
		os.Stderr.WriteString(err.Error())
		return err
	}

	return nil
}
