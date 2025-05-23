// Copyright 2014 The go-ethereum Authors
// (original work)
// Copyright 2024 The Erigon Authors
// (modifications)
// This file is part of Erigon.
//
// Erigon is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Erigon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Erigon. If not, see <http://www.gnu.org/licenses/>.

// evm executes EVM code snippets.
package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/erigontech/erigon-lib/log/v3"

	"github.com/erigontech/erigon/cmd/evm/internal/t8ntool"
	"github.com/erigontech/erigon/cmd/utils/flags"
	"github.com/erigontech/erigon/params"
	cli2 "github.com/erigontech/erigon/turbo/cli"
)

var (
	app = cli2.NewApp(params.GitCommit, "the evm command line interface")

	DebugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "output full trace logs",
	}
	MemProfileFlag = cli.StringFlag{
		Name:  "memprofile",
		Usage: "creates a memory profile at the given path",
	}
	CPUProfileFlag = cli.StringFlag{
		Name:  "cpuprofile",
		Usage: "creates a CPU profile at the given path",
	}
	StatDumpFlag = cli.BoolFlag{
		Name:  "statdump",
		Usage: "displays stack and heap memory information",
	}
	CodeFlag = cli.StringFlag{
		Name:  "code",
		Usage: "EVM code",
	}
	CodeFileFlag = cli.StringFlag{
		Name:  "codefile",
		Usage: "File containing EVM code. If '-' is specified, code is read from stdin ",
	}
	GasFlag = cli.Uint64Flag{
		Name:  "gas",
		Usage: "gas limit for the evm",
		Value: 10000000000,
	}
	PriceFlag = flags.BigFlag{
		Name:  "price",
		Usage: "price set for the evm",
		Value: new(big.Int),
	}
	ValueFlag = flags.BigFlag{
		Name:  "value",
		Usage: "value set for the evm",
		Value: new(big.Int),
	}
	DumpFlag = cli.BoolFlag{
		Name:  "dump",
		Usage: "dumps the state after the run",
	}
	InputFlag = cli.StringFlag{
		Name:  "input",
		Usage: "input for the EVM",
	}
	InputFileFlag = cli.StringFlag{
		Name:  "inputfile",
		Usage: "file containing input for the EVM",
	}
	VerbosityFlag = cli.IntFlag{
		Name:  "verbosity",
		Usage: "sets the verbosity level",
	}
	BenchFlag = cli.BoolFlag{
		Name:  "bench",
		Usage: "benchmark the execution",
	}
	CreateFlag = cli.BoolFlag{
		Name:  "create",
		Usage: "indicates the action should be create rather than call",
	}
	GenesisFlag = cli.StringFlag{
		Name:  "prestate",
		Usage: "JSON file with prestate (genesis) config",
	}
	MachineFlag = cli.BoolFlag{
		Name:  "json",
		Usage: "output trace logs in machine readable format (json)",
	}
	SenderFlag = cli.StringFlag{
		Name:  "sender",
		Usage: "The transaction origin",
	}
	ReceiverFlag = cli.StringFlag{
		Name:  "receiver",
		Usage: "The transaction receiver (execution context)",
	}
	DisableMemoryFlag = cli.BoolFlag{
		Name:  "nomemory",
		Usage: "disable memory output",
	}
	DisableStackFlag = cli.BoolFlag{
		Name:  "nostack",
		Usage: "disable stack output",
	}
	DisableStorageFlag = cli.BoolFlag{
		Name:  "nostorage",
		Usage: "disable storage output",
	}
	DisableReturnDataFlag = cli.BoolFlag{
		Name:  "noreturndata",
		Usage: "disable return data output",
	}
)

var stateTransitionCommand = cli.Command{
	Name:    "transition",
	Aliases: []string{"t8n"},
	Usage:   "executes a full state transition",
	Action:  t8ntool.Main,
	Flags: []cli.Flag{
		&t8ntool.TraceFlag,
		&t8ntool.TraceDisableMemoryFlag,
		&t8ntool.TraceDisableStackFlag,
		&t8ntool.TraceDisableReturnDataFlag,
		&t8ntool.OutputBasedir,
		&t8ntool.OutputAllocFlag,
		&t8ntool.OutputResultFlag,
		&t8ntool.OutputBodyFlag,
		&t8ntool.InputAllocFlag,
		&t8ntool.InputEnvFlag,
		&t8ntool.InputTxsFlag,
		&t8ntool.ForknameFlag,
		&t8ntool.ChainIDFlag,
		&t8ntool.VerbosityFlag,
	},
}

func init() {
	app.Flags = []cli.Flag{
		&BenchFlag,
		&CreateFlag,
		&DebugFlag,
		&VerbosityFlag,
		&CodeFlag,
		&CodeFileFlag,
		&GasFlag,
		&PriceFlag,
		&ValueFlag,
		&DumpFlag,
		&InputFlag,
		&InputFileFlag,
		&MemProfileFlag,
		&CPUProfileFlag,
		&StatDumpFlag,
		&GenesisFlag,
		&MachineFlag,
		&SenderFlag,
		&ReceiverFlag,
		&DisableMemoryFlag,
		&DisableStackFlag,
		&DisableStorageFlag,
		&DisableReturnDataFlag,
	}
	app.Commands = []*cli.Command{
		&compileCommand,
		&disasmCommand,
		&runCommand,
		&stateTestCommand,
		&stateTransitionCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		code := 1
		if ec, ok := err.(*t8ntool.NumberedError); ok {
			code = ec.ExitCode()
		}
		_, printErr := fmt.Fprintln(os.Stderr, err)
		if printErr != nil {
			log.Warn("print error", "err", printErr)
		}
		os.Exit(code)
	}
}
