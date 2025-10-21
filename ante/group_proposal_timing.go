package ante

import (
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
)

// GroupProposalTimingDecorator checks if group proposals are being executed after the voting period ends
type GroupProposalTimingDecorator struct {
	groupKeeper groupkeeper.Keeper
}

// NewGroupProposalTimingDecorator creates a new GroupProposalTimingDecorator
func NewGroupProposalTimingDecorator(groupKeeper groupkeeper.Keeper) GroupProposalTimingDecorator {
	return GroupProposalTimingDecorator{
		groupKeeper: groupKeeper,
	}
}

// AnteHandle checks if group proposal execution happens after voting period ends
func (gptd GroupProposalTimingDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		if execMsg, ok := msg.(*group.MsgExec); ok {
			// Get the proposal to check its voting period
			proposalResp, err := gptd.groupKeeper.Proposal(ctx, &group.QueryProposalRequest{
				ProposalId: execMsg.ProposalId,
			})
			if err != nil {
				return ctx, errors.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"failed to get proposal %d: %v",
					execMsg.ProposalId,
					err,
				)
			}

			proposal := proposalResp.Proposal

			// Check if current block time is after voting period end
			currentTime := ctx.BlockTime()
			votingPeriodEnd := proposal.VotingPeriodEnd

			if currentTime.Before(votingPeriodEnd) {
				return ctx, errors.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"proposal %d cannot be executed yet. Voting period ends at %s, current time is %s. Execute only after voting period ends",
					execMsg.ProposalId,
					votingPeriodEnd.Format(time.RFC3339),
					currentTime.Format(time.RFC3339),
				)
			}

			// Additional check: ensure proposal is in ACCEPTED status
			if proposal.Status != group.PROPOSAL_STATUS_ACCEPTED {
				return ctx, errors.Wrapf(
					sdkerrors.ErrInvalidRequest,
					"proposal %d is not in ACCEPTED status (current: %s). Only accepted proposals can be executed",
					execMsg.ProposalId,
					proposal.Status.String(),
				)
			}
		}
	}

	return next(ctx, tx, simulate)
}
