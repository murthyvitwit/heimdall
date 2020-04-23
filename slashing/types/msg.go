package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgUnjail{}

// MsgUnjail - struct for unjailing jailed validator
type MsgUnjail struct {
	ValidatorAddr types.HeimdallAddress `json:"validatorAddr"`
}

// NewMsgUnjail creates a new MsgUnjail instance
func NewMsgUnjail(validatorAddr types.HeimdallAddress) MsgUnjail {
	return MsgUnjail{
		ValidatorAddr: validatorAddr,
	}
}

//nolint
func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) Type() string  { return "unjail" }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.ValidatorAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgUnjail) ValidateBasic() sdk.Error {
	if msg.ValidatorAddr.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid from %v", msg.ValidatorAddr.String())
	}

	return nil
}

// Tick Msg

// TickMsg - struct for unjailing jailed validator
type MsgTick struct {
	Proposer         types.HeimdallAddress `json:"proposer"`
	SlashingInfoHash types.HeimdallHash    `json:"accountRootHash"`
}

func NewMsgtick(proposer types.HeimdallAddress, slashingInfoHash types.HeimdallHash) MsgTick {
	return MsgTick{
		Proposer:         proposer,
		SlashingInfoHash: slashingInfoHash,
	}
}

// Type returns message type
func (msg MsgTick) Type() string {
	return "tick"
}

func (msg MsgTick) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgTick) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.Proposer)}
}

func (msg MsgTick) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgTick) ValidateBasic() sdk.Error {
	if bytes.Equal(msg.SlashingInfoHash.Bytes(), helper.ZeroHash.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid slashing info hash %v", msg.SlashingInfoHash.String())
	}

	if msg.Proposer.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.Proposer.String())
	}
	return nil
}

//
// Msg Tick Ack
//

type MsgTickAck struct {
	From     types.HeimdallAddress `json:"from"`
	TxHash   types.HeimdallHash    `json:"tx_hash"`
	LogIndex uint64                `json:"log_index"`
}

func NewMsgtickAck(from types.HeimdallAddress, txHash types.HeimdallHash, logIndex uint64) MsgTickAck {
	return MsgTickAck{
		From:     from,
		TxHash:   txHash,
		LogIndex: logIndex,
	}
}

// Type returns message type
func (msg MsgTickAck) Type() string {
	return "tick-ack"
}

func (msg MsgTickAck) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgTickAck) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgTickAck) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgTickAck) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid from %v", msg.From.String())
	}

	return nil
}
