package live

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/ethereum/go-ethereum/node"
	"github.com/stretchr/testify/require"
)

var (
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr   = crypto.PubkeyToAddress(testKey.PublicKey)
)

func withTracer(name, config string) func(nodeConf *node.Config, ethConf *ethconfig.Config) {
	return func(nodeConf *node.Config, ethConf *ethconfig.Config) {
		ethConf.VMTrace = name
		ethConf.VMTraceJsonConfig = config

	}
}

func newBackend(tracer, tracerConfig string) *simulated.Backend {
	genesis := core.DeveloperGenesisBlock(0, &testAddr)
	return simulated.NewDeterministicBackend(genesis.Alloc, withTracer(tracer, tracerConfig))
}

func TestSimpleTracer(t *testing.T) {
	traceOutputPath := filepath.ToSlash(t.TempDir())
	traceOutputFilename := path.Join(traceOutputPath, "result.jsonl")
	backend := newBackend("liveTracer", fmt.Sprintf(`{"path": "%s"}`, traceOutputPath))
	defer backend.Close()
	client := backend.Client()
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	_, err = sendTx(client, chainId, testKey, &types.DynamicFeeTx{
		Nonce: 0,
		Gas:   21000,
		To:    &common.Address{0xff},
		Value: big.NewInt(1),
	})
	if err != nil {
		t.Fatal(err)
	}
	backend.Commit()

	// read tracer output file and print
	result, err := parseResult(traceOutputFilename)

	expectedRaw, err := os.ReadFile("testdata/basic_tx.json")
	if err != nil {
		t.Fatal(err)
	}
	var expected []blockResult
	if err := json.Unmarshal(expectedRaw, &expected); err != nil {
		t.Fatal(err)
	}
	require.EqualValues(t, expected, result)
}

func sendTx(client simulated.Client, chainId *big.Int, key *ecdsa.PrivateKey, tx *types.DynamicFeeTx) (common.Hash, error) {
	feeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}
	tipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return common.Hash{}, err
	}
	tx.GasFeeCap = feeCap
	tx.GasTipCap = tipCap
	signed, err := types.SignTx(types.NewTx(tx), types.NewLondonSigner(chainId), key)
	if err != nil {
		return common.Hash{}, err
	}
	if err := client.SendTransaction(context.Background(), signed); err != nil {
		return common.Hash{}, err
	}
	return signed.Hash(), nil
}

func parseResult(path string) ([]blockResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var result []blockResult
	dec := json.NewDecoder(f)
	for dec.More() {
		var r blockResult
		if err := dec.Decode(&r); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}
