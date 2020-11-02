package gov

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/gov/types"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/supply"
	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
)

// Keeper governance Keeper
type Keeper struct {
	// The reference to the Paramstore to get and set gov specific params
	paramSpace subspace.Subspace

	// The SupplyKeeper to reduce the supply of the network
	supplyKeeper supply.Keeper

	// The reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	sk staking.Keeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType

	// Proposal router
	router Router
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	paramSpace subspace.Subspace,
	supplyKeeper supply.Keeper,
	sk staking.Keeper,
	codespace sdk.CodespaceType,
	rtr Router,
) Keeper {

	// ensure governance module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr.Empty() {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// It is vital to seal the governance proposal router here as to not allow
	// further handlers to be registered after the keeper is created since this
	// could create invalid or non-deterministic behavior.
	rtr.Seal()

	return Keeper{
		storeKey:     key,
		paramSpace:   paramSpace,
		supplyKeeper: supplyKeeper,
		sk:           sk,
		cdc:          cdc,
		codespace:    codespace,
		router:       rtr,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}