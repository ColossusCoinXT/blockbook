package colx

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/json"

	"github.com/golang/glog"
)

// ColxRPC is an interface to JSON-RPC bitcoind service.
type ColxRPC struct {
	*btc.BitcoinRPC
}

// NewColxRPC returns new ColxRPC instance.
func NewColxRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	s := &ColxRPC{
		b.(*btc.BitcoinRPC),
	}
	s.RPCMarshaler = btc.JSONMarshalerV1{}
	s.ChainConfig.SupportsEstimateFee = true
	s.ChainConfig.SupportsEstimateSmartFee = false

	return s, nil
}

// Initialize initializes ColxRPC instance.
func (b *ColxRPC) Initialize() error {
	ci, err := b.GetChainInfo()
	if err != nil {
		return err
	}
	chainName := ci.Chain

	glog.Info("Chain name ", chainName)
	params := GetChainParams(chainName)

	// always create parser
	b.Parser = NewColxParser(params, b.ChainConfig)

	// parameters for getInfo request
	b.Testnet = false
	b.Network = "livenet"

	glog.Info("rpc: block chain ", params.Name)

	return nil
}
