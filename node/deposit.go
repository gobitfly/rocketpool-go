package node

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// Estimate the gas of Deposit
func EstimateDepositGas(rp *rocketpool.RocketPool, bondAmount *big.Int, minimumNodeFee float64, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, salt *big.Int, expectedMinipoolAddress common.Address, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketNodeDeposit, err := getRocketNodeDeposit(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	return rocketNodeDeposit.GetTransactionGasInfo(opts, "deposit", bondAmount, eth.EthToWei(minimumNodeFee), validatorPubkey[:], validatorSignature[:], depositDataRoot, salt, expectedMinipoolAddress)
}

// Make a node deposit
func Deposit(rp *rocketpool.RocketPool, bondAmount *big.Int, minimumNodeFee float64, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, salt *big.Int, expectedMinipoolAddress common.Address, opts *bind.TransactOpts) (*types.Transaction, error) {
	rocketNodeDeposit, err := getRocketNodeDeposit(rp, nil)
	if err != nil {
		return nil, err
	}
	tx, err := rocketNodeDeposit.Transact(opts, "deposit", bondAmount, eth.EthToWei(minimumNodeFee), validatorPubkey[:], validatorSignature[:], depositDataRoot, salt, expectedMinipoolAddress)
	if err != nil {
		return nil, fmt.Errorf("Could not make node deposit: %w", err)
	}
	return tx, nil
}

// Get the amount of ETH in the node's deposit credit bank
func GetNodeDepositCredit(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) (*big.Int, error) {
	rocketNodeDeposit, err := getRocketNodeDeposit(rp, opts)
	if err != nil {
		return nil, err
	}

	creditBalance := new(*big.Int)
	if err := rocketNodeDeposit.Call(opts, creditBalance, "getNodeDepositCredit", nodeAddress); err != nil {
		return nil, fmt.Errorf("Could not get node deposit credit: %w", err)
	}
	return *creditBalance, nil
}

// Get contracts
var rocketNodeDepositLock sync.Mutex

func getRocketNodeDeposit(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*rocketpool.Contract, error) {
	rocketNodeDepositLock.Lock()
	defer rocketNodeDepositLock.Unlock()
	return rp.GetContract("rocketNodeDeposit", opts)
}
