// +build unittest

package colx

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func Test_GetAddrDescFromAddress_Mainnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{address: "DRM8TaiY38qcHbgdytp8oETreobBLHtpeE"},
			want:    "76a914dda91c0396050d660f9c0e38f78064486bbfcb2c88ac",
			wantErr: false,
		},
	}
	parser := NewColxParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_GetAddressesFromAddrDesc(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want2   bool
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{script: "76a914dda91c0396050d660f9c0e38f78064486bbfcb2c88ac"},
			want:    []string{"DRM8TaiY38qcHbgdytp8oETreobBLHtpeE"},
			want2:   true,
			wantErr: false,
		},
		{
			name:    "pubkey",
			args:    args{script: "210251c5555ff3c684aebfca92f5329e2f660da54856299da067060a1bcf5e8fae73ac"},
			want:    []string{"DKL3QzCbJqrHpRKAHvEqsomsDhkQPvVzZg"},
			want2:   false,
			wantErr: false,
		},
	}

	parser := NewColxParser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.script)
			got, got2, err := parser.GetAddressesFromAddrDesc(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddressesFromAddrDesc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetAddressesFromAddrDesc() = %v, want %v", got2, tt.want2)
			}
		})
	}
}

var (
	// regular transaction
	testTx1       bchain.Tx
	testTxPacked1 = "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4b02102704fe2fd15908fabe6d6de11302ffffffff978c82600e136c1290a81b27526449d7d26f936ce957487c0c01000000000000003ffffffd4a0000000d2f6e6f64655374726174756d2f000000000100671a4c370000001976a914963c3306f96c2d2b70c89023480c6bbad7d6f0f788ac00000000"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:      "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4b02102704fe2fd15908fabe6d6de11302ffffffff978c82600e136c1290a81b27526449d7d26f936ce957487c0c01000000000000003ffffffd4a0000000d2f6e6f64655374726174756d2f000000000100671a4c370000001976a914963c3306f96c2d2b70c89023480c6bbad7d6f0f788ac00000000",
		Txid:     "835302f89d29e0f3b2a788d513964d937fa04144bdc6ae1c007bd47e3423e0c2",
		LockTime: 0,
		Vin: []bchain.Vin{
			{
				Coinbase: "02102704fe2fd15908fabe6d6de11302ffffffff978c82600e136c1290a81b27526449d7d26f936ce957487c0c01000000000000003ffffffd4a0000000d2f6e6f64655374726174756d2f",
				Sequence: 0,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(2375),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a914963c3306f96c2d2b70c89023480c6bbad7d6f0f788ac",
					Addresses: []string{
						 "DJqU1kmZhbaVJVwQ57kxnpFLCwsk72KuJ6",
					},
				},
			},
		},
		Blocktime: 1506881534,
		Time:      1506881534,
	}
}

func Test_PackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *ColxParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "colx-1",
			args: args{
				tx:        testTx1,
				height:    10000,
				blockTime: 1506881534,
				parser:    NewColxParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func Test_UnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *ColxParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "colx-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewColxParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   10000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type testBlock struct {
	size int
	time int64
	txs  []string
}

var testParseBlockTxs = map[int]testBlock{
	10000: {
		size: 241,
		time: 1506881534,
		txs: []string{
			"835302f89d29e0f3b2a788d513964d937fa04144bdc6ae1c007bd47e3423e0c2"
		},
	},
}

func helperLoadBlock(t *testing.T, height int) []byte {
	name := fmt.Sprintf("block_dump.%d", height)
	path := filepath.Join("testdata", name)

	d, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	d = bytes.TrimSpace(d)

	b := make([]byte, hex.DecodedLen(len(d)))
	_, err = hex.Decode(b, d)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestParseBlock(t *testing.T) {
	p := NewColxParser(GetChainParams("main"), &btc.Configuration{})

	for height, tb := range testParseBlockTxs {
		b := helperLoadBlock(t, height)

		blk, err := p.ParseBlock(b)
		if err != nil {
			t.Fatal(err)
		}

		if blk.Size != tb.size {
			t.Errorf("ParseBlock() block size: got %d, want %d", blk.Size, tb.size)
		}

		if blk.Time != tb.time {
			t.Errorf("ParseBlock() block time: got %d, want %d", blk.Time, tb.time)
		}

		if len(blk.Txs) != len(tb.txs) {
			t.Errorf("ParseBlock() number of transactions: got %d, want %d", len(blk.Txs), len(tb.txs))
		}

		for ti, tx := range tb.txs {
			if blk.Txs[ti].Txid != tx {
				t.Errorf("ParseBlock() transaction %d: got %s, want %s", ti, blk.Txs[ti].Txid, tx)
			}
		}
	}
}
