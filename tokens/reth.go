package tokens

import (
    "fmt"
    "math/big"
    "sync"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"

    "github.com/rocket-pool/rocketpool-go/rocketpool"
    "github.com/rocket-pool/rocketpool-go/utils/eth"
)


// Get the rETH contract ETH balance
func GetRETHContractETHBalance(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*big.Int, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return nil, err
    }
    return contractETHBalance(rp, rocketTokenRETH, opts)
}


// Get rETH total supply
func GetRETHTotalSupply(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*big.Int, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return nil, err
    }
    return totalSupply(rocketTokenRETH, "rETH", opts)
}


// Get rETH balance
func GetRETHBalance(rp *rocketpool.RocketPool, address common.Address, opts *bind.CallOpts) (*big.Int, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return nil, err
    }
    return balanceOf(rocketTokenRETH, "rETH", address, opts)
}


// Get the ETH value of an amount of rETH
func GetETHValueOfRETH(rp *rocketpool.RocketPool, rethAmount *big.Int, opts *bind.CallOpts) (*big.Int, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return nil, err
    }
    ethValue := new(*big.Int)
    if err := rocketTokenRETH.Call(opts, ethValue, "getEthValue", rethAmount); err != nil {
        return nil, fmt.Errorf("Could not get ETH value of rETH amount: %w", err)
    }
    return *ethValue, nil
}


// Get the rETH value of an amount of ETH
func GetRETHValueOfETH(rp *rocketpool.RocketPool, ethAmount *big.Int, opts *bind.CallOpts) (*big.Int, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return nil, err
    }
    rethValue := new(*big.Int)
    if err := rocketTokenRETH.Call(opts, rethValue, "getRethValue", ethAmount); err != nil {
        return nil, fmt.Errorf("Could not get rETH value of ETH amount: %w", err)
    }
    return *rethValue, nil
}


// Get the current ETH : rETH exchange rate
func GetRETHExchangeRate(rp *rocketpool.RocketPool, opts *bind.CallOpts) (float64, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return 0, err
    }
    exchangeRate := new(*big.Int)
    if err := rocketTokenRETH.Call(opts, exchangeRate, "getExchangeRate"); err != nil {
        return 0, fmt.Errorf("Could not get rETH exchange rate: %w", err)
    }
    return eth.WeiToEth(*exchangeRate), nil
}


// Get the total amount of ETH collateral available for rETH trades
func GetRETHTotalCollateral(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*big.Int, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return nil, err
    }
    totalCollateral := new(*big.Int)
    if err := rocketTokenRETH.Call(opts, totalCollateral, "getTotalCollateral"); err != nil {
        return nil, fmt.Errorf("Could not get rETH total collateral: %w", err)
    }
    return *totalCollateral, nil
}


// Get the rETH collateralization rate
func GetRETHCollateralRate(rp *rocketpool.RocketPool, opts *bind.CallOpts) (float64, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return 0, err
    }
    collateralRate := new(*big.Int)
    if err := rocketTokenRETH.Call(opts, collateralRate, "getCollateralRate"); err != nil {
        return 0, fmt.Errorf("Could not get rETH collateral rate: %w", err)
    }
    return eth.WeiToEth(*collateralRate), nil
}


// Transfer rETH
func TransferRETH(rp *rocketpool.RocketPool, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*types.Receipt, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return nil, err
    }
    return transfer(rp.Client, rocketTokenRETH, "rETH", to, amount, opts)
}


// Burn rETH for ETH
func BurnRETH(rp *rocketpool.RocketPool, amount *big.Int, opts *bind.TransactOpts) (*types.Receipt, error) {
    rocketTokenRETH, err := getRocketTokenRETH(rp)
    if err != nil {
        return nil, err
    }
    txReceipt, err := rocketTokenRETH.Transact(opts, "burn", amount)
    if err != nil {
        return nil, fmt.Errorf("Could not burn rETH: %w", err)
    }
    return txReceipt, nil
}


// Get contracts
var rocketTokenRETHLock sync.Mutex
func getRocketTokenRETH(rp *rocketpool.RocketPool) (*rocketpool.Contract, error) {
    rocketTokenRETHLock.Lock()
    defer rocketTokenRETHLock.Unlock()
    return rp.GetContract("rocketTokenRETH")
}

