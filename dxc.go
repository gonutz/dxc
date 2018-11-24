package dxc

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	dll        = syscall.NewLazyDLL("D3DCompiler_47.dll")
	d3DCompile = dll.NewProc("D3DCompile")
)

func Compile(
	sourceCode string,
	entryPoint string,
	target string,
	compileFlags uint,
	effectFlags uint,
) ([]byte, error) {
	sourceCodeBytes := []byte(sourceCode)
	entryPointBytes := append([]byte(entryPoint), 0)
	targetBytes := append([]byte(target), 0)
	var output, err *blob
	ret, _, _ := d3DCompile.Call(
		uintptr(unsafe.Pointer(&sourceCodeBytes[0])),
		uintptr(len(sourceCodeBytes)),
		0, // source name
		0, // defines
		0, // include handler
		uintptr(unsafe.Pointer(&entryPointBytes[0])),
		uintptr(unsafe.Pointer(&targetBytes[0])),
		uintptr(compileFlags),
		uintptr(effectFlags),
		uintptr(unsafe.Pointer(&output)),
		uintptr(unsafe.Pointer(&err)),
	)
	if ret == 0 {
		return output.bytes(), nil
	} else {
		return nil, errors.New(string(err.bytes()))
	}
}

type blob struct {
	vtbl *blobVtbl
}

type blobVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetBufferPointer uintptr
	GetBufferSize    uintptr
}

func (b *blob) GetBufferPointer() uintptr {
	ret, _, _ := syscall.Syscall(
		b.vtbl.GetBufferPointer,
		1,
		uintptr(unsafe.Pointer(b)),
		0,
		0,
	)
	return ret
}

func (b *blob) GetBufferSize() uint32 {
	ret, _, _ := syscall.Syscall(
		b.vtbl.GetBufferSize,
		1,
		uintptr(unsafe.Pointer(b)),
		0,
		0,
	)
	return uint32(ret)
}

func (b *blob) bytes() []byte {
	size := b.GetBufferSize()
	buf := make([]byte, size)
	ptr := b.GetBufferPointer()
	for i := range buf {
		buf[i] = *((*byte)(unsafe.Pointer(ptr)))
		ptr++
	}
	return buf
}

// compiler flags
const (
	// DEBUG directs the compiler to insert debug file/line/type/symbol
	// information into the output code.
	DEBUG = 1 << 0

	// SKIP_VALIDATION directs the compiler not to validate the generated code
	// against known capabilities and constraints. We recommend that you use
	// this constant only with shaders that have been successfully compiled in
	// the past. DirectX always validates shaders before it sets them to a
	// device.
	SKIP_VALIDATION = 1 << 1

	// SKIP_OPTIMIZATION directs the compiler to skip optimization steps during
	// code generation. We recommend that you set this constant for debug
	// purposes only.
	SKIP_OPTIMIZATION = 1 << 2

	// PACK_MATRIX_ROW_MAJOR directs the compiler to pack matrices in row-major
	// order on input and output from the shader.
	PACK_MATRIX_ROW_MAJOR = 1 << 3

	// PACK_MATRIX_COLUMN_MAJOR directs the compiler to pack matrices in
	// column-major order on input and output from the shader. This type of
	// packing is generally more efficient because a series of dot-products can
	// then perform vector-matrix multiplication.
	PACK_MATRIX_COLUMN_MAJOR = 1 << 4

	// PARTIAL_PRECISION directs the compiler to perform all computations with
	// partial precision. If you set this constant, the compiled code might run
	// faster on some hardware.
	PARTIAL_PRECISION = 1 << 5

	// FORCE_VS_SOFTWARE_NO_OPT directs the compiler to compile a vertex shader
	// for the next highest shader profile. This constant turns debugging on and
	// optimizations off.
	FORCE_VS_SOFTWARE_NO_OPT = 1 << 6

	// FORCE_PS_SOFTWARE_NO_OPT directs the compiler to compile a pixel shader
	// for the next highest shader profile. This constant also turns debugging
	// on and optimizations off.
	FORCE_PS_SOFTWARE_NO_OPT = 1 << 7

	// NO_PRESHADER directs the compiler to disable Preshaders. If you set this
	// constant, the compiler does not pull out static expression for
	// evaluation.
	NO_PRESHADER = 1 << 8

	// AVOID_FLOW_CONTROL directs the compiler to not use flow-control
	// constructs where possible.
	AVOID_FLOW_CONTROL = 1 << 9

	// PREFER_FLOW_CONTROL directs the compiler to use flow-control constructs
	// where possible.
	PREFER_FLOW_CONTROL = 1 << 10

	// ENABLE_STRICTNESS forces strict compile, which might not allow for legacy
	// syntax. By default, the compiler disables strictness on deprecated
	// syntax.
	ENABLE_STRICTNESS = 1 << 11

	// ENABLE_BACKWARDS_COMPATIBILITY directs the compiler to enable older
	// shaders to compile to 5_0 targets.
	ENABLE_BACKWARDS_COMPATIBILITY = 1 << 12

	// IEEE_STRICTNESS forces the IEEE strict compile.
	IEEE_STRICTNESS = 1 << 13

	// OPTIMIZATION_LEVEL0 directs the compiler to use the lowest optimization
	// level. If you set this constant, the compiler might produce slower code
	// but produces the code more quickly. Set this constant when you develop
	// the shader iteratively.
	OPTIMIZATION_LEVEL0 = 1 << 14

	// OPTIMIZATION_LEVEL1 directs the compiler to use the second lowest
	// optimization level.
	OPTIMIZATION_LEVEL1 = 0

	// OPTIMIZATION_LEVEL2 directs the compiler to use the second highest
	// optimization level.
	OPTIMIZATION_LEVEL2 = 1<<14 | 1<<15

	// OPTIMIZATION_LEVEL3 directs the compiler to use the highest optimization
	// level. If you set this constant, the compiler produces the best possible
	// code but might take significantly longer to do so. Set this constant for
	// final builds of an application when performance is the most important
	// factor.
	OPTIMIZATION_LEVEL3 = 1 << 15

	// WARNINGS_ARE_ERRORS directs the compiler to treat all warnings as errors
	// when it compiles the shader code. We recommend that you use this constant
	// for new shader code, so that you can resolve all warnings and lower the
	// number of hard-to-find code defects.
	WARNINGS_ARE_ERRORS = 1 << 18

	// RESOURCES_MAY_ALIAS directs the compiler to assume that unordered access
	// views (UAVs) and shader resource views (SRVs) may alias for cs_5_0.
	RESOURCES_MAY_ALIAS = 1 << 19

	// ENABLE_UNBOUNDED_DESCRIPTOR_TABLES directs the compiler to enable
	// unbounded descriptor tables.
	ENABLE_UNBOUNDED_DESCRIPTOR_TABLES = 1 << 20

	// ALL_RESOURCES_BOUND directs the compiler to ensure all resources are
	// bound.
	ALL_RESOURCES_BOUND = 1 << 21
)

// effect flags
const (
	// EFFECT_CHILD_EFFECT compile the effects (.fx) file to a child effect.
	// Child effects have no initializers for any shared values because these
	// child effects are initialized in the master effect (the effect pool).
	EFFECT_CHILD_EFFECT = 1 << 0

	// EFFECT_ALLOW_SLOW_OPS disables performance mode and allows for mutable
	// state objects. By default, performance mode is enabled. Performance mode
	// disallows mutable state objects by preventing non-literal expressions
	// from appearing in state object definitions.
	EFFECT_ALLOW_SLOW_OPS = 1 << 1
)
