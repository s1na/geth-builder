package live

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	tracers.LiveDirectory.Register("liveTracer", newLiveTracer)
}

type liveTracer struct {
	logger *lumberjack.Logger
}

type tracerConfig struct {
	Path string `json:"path"`
}

type blockResult struct {
	Number uint64      `json:"blockNumber"`
	Hash   common.Hash `json:"hash"`
}

func newLiveTracer(cfg json.RawMessage) (*tracing.Hooks, error) {
	var c tracerConfig
	if cfg != nil {
		if err := json.Unmarshal(cfg, &c); err != nil {
			return nil, err
		}
	}
	if c.Path == "" {
		return nil, errors.New("no path specified")
	}
	logger := &lumberjack.Logger{
		Filename: filepath.Join(c.Path, "result.jsonl"),
	}
	t := &liveTracer{logger: logger}
	return &tracing.Hooks{
		OnGenesisBlock: t.OnGenesisBlock,
		OnBlockStart:   t.OnBlockStart,
	}, nil
}

func (t *liveTracer) OnGenesisBlock(b *types.Block, alloc types.GenesisAlloc) {
	result := blockResult{
		Number: b.NumberU64(),
		Hash:   b.Hash(),
	}
	data, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("failed to marshal block result: %v\n", err)
		return
	}
	t.logger.Write(data)
}

func (t *liveTracer) OnBlockStart(ev tracing.BlockEvent) {
	b := ev.Block
	result := blockResult{
		Number: b.NumberU64(),
		Hash:   b.Hash(),
	}
	data, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("failed to marshal block result: %v\n", err)
		return
	}
	t.logger.Write(data)
}
