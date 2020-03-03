package topup

import (
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/maticnetwork/heimdall/bank"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	// DefaultValue default value
	DefaultValue = []byte{0x01}
	// ValidatorTopupKey represents validator topup key
	ValidatorTopupKey = []byte{0x80} // prefix for each key to a validator
	// TopupSequencePrefixKey represents topup sequence prefix key
	TopupSequencePrefixKey = []byte{0x81}
)

// ModuleCommunicator manager to access validator info
type ModuleCommunicator interface {
	// AddFeeToDividendAccount add fee to dividend account
	AddFeeToDividendAccount(ctx sdk.Context, valID hmTypes.ValidatorID, fee *big.Int) sdk.Error
	// GetValidatorFromValID get validator from validator id
	GetValidatorFromValID(ctx sdk.Context, valID hmTypes.ValidatorID) (validator hmTypes.Validator, ok bool)
}

// Keeper stores all related data
type Keeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
	// code space
	codespace sdk.CodespaceType
	// param subspace
	paramSpace params.Subspace
	// bank keeper
	bk bank.Keeper
	// module manager
	vm ModuleCommunicator
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
	bankKeeper bank.Keeper,
	vm ModuleCommunicator,
) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        storeKey,
		paramSpace: paramSpace,
		codespace:  codespace,
		bk:         bankKeeper,
		vm:         vm,
	}
}

// Codespace returns the keeper's codespace.
func (keeper Keeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}

// Logger returns a module-specific logger
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// GetSendEnabled returns the current SendEnabled
// nolint: errcheck
func (keeper Keeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, types.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper Keeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeySendEnabled, &enabled)
}

//
// Topup methods
//

// GetTopupKey drafts the topup key for address
func GetTopupKey(address []byte) []byte {
	return append(ValidatorTopupKey, address...)
}

// GetTopupSequenceKey drafts topup sequence for address
func GetTopupSequenceKey(sequence uint64) []byte {
	return append(TopupSequencePrefixKey, []byte(strconv.FormatUint(sequence, 10))...)
}

// GetValidatorTopup returns validator toptup information
func (keeper Keeper) GetValidatorTopup(ctx sdk.Context, addr hmTypes.HeimdallAddress) (*types.ValidatorTopup, error) {
	store := ctx.KVStore(keeper.key)

	// check if topup exists
	key := GetTopupKey(addr.Bytes())
	if !store.Has(key) {
		return nil, nil
	}

	// unmarshall validator and return
	validatorTopup, err := types.UnmarshallValidatorTopup(keeper.cdc, store.Get(key))
	if err != nil {
		return nil, err
	}

	// return true if validator
	return &validatorTopup, nil
}

// SetValidatorTopup sets validator topup object
func (keeper Keeper) SetValidatorTopup(ctx sdk.Context, addr hmTypes.HeimdallAddress, validatorTopup types.ValidatorTopup) error {
	store := ctx.KVStore(keeper.key)

	// validator topup
	bz, err := types.MarshallValidatorTopup(keeper.cdc, validatorTopup)
	if err != nil {
		return err
	}

	// store validator with address prefixed with validator key as index
	store.Set(GetTopupKey(addr.Bytes()), bz)
	keeper.Logger(ctx).Debug("Validator topup stored", "key", hex.EncodeToString(GetTopupKey(addr.Bytes())), "totalTopups", validatorTopup.Copy().TotalTopups)

	return nil
}

// GetTopupSequence checks if topup already exists
func (keeper Keeper) GetTopupSequence(ctx sdk.Context, sequence uint64) uint64 {
	store := ctx.KVStore(keeper.key)
	sequenceKey := GetTopupSequenceKey(sequence)
	if store.Has(sequenceKey) {
		result, err := strconv.ParseUint(string(store.Get(sequenceKey)), 10, 64)
		if err == nil {
			return uint64(result)
		}
	}
	return 0
}

// SetTopupSequence sets mapping for sequence id to bool
func (keeper Keeper) SetTopupSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(keeper.key)
	store.Set(GetTopupSequenceKey(sequence), DefaultValue)
}

// HasTopupSequence checks if topup already exists
func (keeper Keeper) HasTopupSequence(ctx sdk.Context, sequence uint64) bool {
	store := ctx.KVStore(keeper.key)
	return store.Has(GetTopupSequenceKey(sequence))
}
