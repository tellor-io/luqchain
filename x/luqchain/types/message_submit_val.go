package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSubmitVal = "submit_val"

var _ sdk.Msg = &MsgSubmitVal{}

func NewMsgSubmitVal(creator string, qdata string, value uint64) *MsgSubmitVal {
	return &MsgSubmitVal{
		Creator: creator,
		Qdata:   qdata,
		Value:   value,
	}
}

func (msg *MsgSubmitVal) Route() string {
	return RouterKey
}

func (msg *MsgSubmitVal) Type() string {
	return TypeMsgSubmitVal
}

func (msg *MsgSubmitVal) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSubmitVal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitVal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
