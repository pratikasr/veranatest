package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"veranatest/x/validatorregistry/keeper"
	"veranatest/x/validatorregistry/types"
)

func createNValidator(keeper keeper.Keeper, ctx context.Context, n int) []types.Validator {
	items := make([]types.Validator, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)
		items[i].MemberId = strconv.Itoa(i)
		items[i].OperatorAddress = strconv.Itoa(i)
		items[i].ConsensusPubkey = strconv.Itoa(i)
		items[i].Status = strconv.Itoa(i)
		items[i].TermEnd = uint64(i)
		_ = keeper.Validator.Set(ctx, items[i].Index, items[i])
	}
	return items
}

func TestValidatorQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNValidator(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetValidatorRequest
		response *types.QueryGetValidatorResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetValidatorRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetValidatorResponse{Validator: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetValidatorRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetValidatorResponse{Validator: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetValidatorRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetValidator(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestValidatorQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNValidator(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllValidatorRequest {
		return &types.QueryAllValidatorRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListValidator(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Validator), step)
			require.Subset(t, msgs, resp.Validator)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListValidator(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Validator), step)
			require.Subset(t, msgs, resp.Validator)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListValidator(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Validator)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListValidator(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
