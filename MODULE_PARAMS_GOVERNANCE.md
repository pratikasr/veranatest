# Module Parameters Governance via Group Module

Complete guide on how to update module parameters using the Cosmos SDK `x/group` module for council-based governance.

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [How It Works](#how-it-works)
- [Step-by-Step Implementation](#step-by-step-implementation)
- [Complete Example: Updating TD Module Parameters](#complete-example-updating-td-module-parameters)
- [Testing Other Modules](#testing-other-modules)
- [Troubleshooting](#troubleshooting)

---

## Overview

This guide demonstrates how to use the **group module** (council governance) to update parameters of **any module** in your Cosmos SDK blockchain. Instead of using the traditional `x/gov` governance module, we use a council-based approach where a predefined group of members votes on parameter changes.

### What You'll Learn

- ✅ Configure module authority to use group policy
- ✅ Create parameter update proposals
- ✅ Vote on proposals with council members
- ✅ Execute approved proposals
- ✅ Verify parameter changes on-chain

### Real Example

We successfully updated the `td` module's `trust_deposit_yield_rate` parameter from `150000000000000000` to `160000000000000000` using council governance.

---

## Prerequisites

Before you begin, ensure you have:

1. ✅ **Council setup** (Group + Group Policy created)
   - See [COUNCIL_GOVERNANCE.md](./COUNCIL_GOVERNANCE.md) for setup
   - Group ID: `1`
   - Group Policy Address: `cosmos1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsfwkgpd`

2. ✅ **Council members** with voting power
   - council-member-1 (weight: 1)
   - council-member-2 (weight: 1)
   - council-member-3 (weight: 1)
   - Decision policy: 2/3 threshold

3. ✅ **Module with MsgUpdateParams** implemented
   - All Cosmos SDK modules have this by default
   - Custom modules should implement it

4. ✅ **Binary access** to your chain
   - `veranatestd` command available
   - Access to keyring with council member keys

---

## How It Works

### Module Authority Pattern

Every Cosmos SDK module has an **authority** - an address that's authorized to update its parameters. By default, this is the governance module account, but it can be set to **any address**, including a **group policy address**.

```
┌─────────────────────────────────────────────────────────┐
│                    Module Authority                      │
│                                                          │
│  Default:  x/gov module account (governance)            │
│  Custom:   Group policy address (council governance) ✓  │
│                                                          │
└─────────────────────────────────────────────────────────┘
                          ↓
                  ┌───────────────┐
                  │ MsgUpdateParams│
                  └───────────────┘
                          ↓
            ┌─────────────────────────┐
            │   Module Parameters      │
            │   Updated On-Chain       │
            └─────────────────────────┘
```

### Governance Flow

```
Council Member    →  Submit Proposal (MsgUpdateParams)
                              ↓
Council Members   →  Vote (YES/NO/ABSTAIN)
                              ↓
Check Threshold   →  Is 2/3 threshold met?
                              ↓
Anyone            →  Execute Proposal
                              ↓
Module            →  Parameters Updated!
```

---

## Step-by-Step Implementation

### Step 1: Set Module Authority to Group Policy

#### 1.1: Edit `app/app_config.go`

Find your module configuration and add the `Authority` field:

**Before:**
```go
{
    Name:   tdmoduletypes.ModuleName,
    Config: appconfig.WrapAny(&tdmoduletypes.Module{}),
},
```

**After:**
```go
{
    Name: tdmoduletypes.ModuleName,
    Config: appconfig.WrapAny(&tdmoduletypes.Module{
        Authority: "cosmos1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsfwkgpd",
    }),
},
```

**Important:** Use your actual group policy address!

#### 1.2: Rebuild Binary

```bash
cd /Users/pratik/veranatest
make install
```

Expected output:
```
--> ensure dependencies have not been modified
all modules verified
--> installing veranatestd
```

#### 1.3: Restart Chain (if needed)

⚠️ **Note:** Authority changes require chain restart or fresh genesis. If you're updating an existing chain, you may need to reset and restart.

For development:
```bash
# Stop current chain
# Reset data (optional, loses all state)
rm -rf ~/.veranatest

# Restart with updated binary
./setup_validator.sh
```

---

### Step 2: Check Current Module Parameters

Before making changes, verify current parameter values:

```bash
# Query td module parameters
veranatestd q td params --node http://localhost:26657
```

**Example Output:**
```yaml
params:
  trust_deposit_share_value: "1000000000000000000"
  trust_deposit_yield_rate: "150000000000000000"
```

**In JSON format (for easier parsing):**
```bash
veranatestd q td params --node http://localhost:26657 -o json | jq '.params'
```

---

### Step 3: Create Parameter Update Proposal

#### 3.1: Get Council Member Addresses

```bash
# Get addresses for all council members
veranatestd keys show council-member-1 --keyring-backend test -a
veranatestd keys show council-member-2 --keyring-backend test -a
veranatestd keys show council-member-3 --keyring-backend test -a
```

**Example Output:**
```
cosmos17hl5uxaglku7a5n3xygu6yk98we6r6t608km6s  (council-member-1)
cosmos1gdmwmz2qyatta9ymwc84wcj0syu3ruljnprlal  (council-member-2)
cosmos1w42mej2dn27vpvjp4udl28yuplyucmczrp6z07  (council-member-3)
```

#### 3.2: Create Proposal JSON File

Create a file named `update_td_params_proposal.json`:

```json
{
  "group_policy_address": "cosmos1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsfwkgpd",
  "proposers": ["cosmos17hl5uxaglku7a5n3xygu6yk98we6r6t608km6s"],
  "metadata": "Update TD module yield rate from 15% to 16%",
  "title": "Update Trust Deposit Yield Rate",
  "summary": "Proposal to increase trust_deposit_yield_rate from 150000000000000000 to 160000000000000000",
  "messages": [
    {
      "@type": "/veranatest.td.v1.MsgUpdateParams",
      "authority": "cosmos1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsfwkgpd",
      "params": {
        "trust_deposit_share_value": "1000000000000000000",
        "trust_deposit_yield_rate": "160000000000000000"
      }
    }
  ]
}
```

**Command to create the file:**
```bash
cat > update_td_params_proposal.json << 'EOF'
{
  "group_policy_address": "cosmos1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsfwkgpd",
  "proposers": ["cosmos17hl5uxaglku7a5n3xygu6yk98we6r6t608km6s"],
  "metadata": "Update TD module yield rate from 15% to 16%",
  "title": "Update Trust Deposit Yield Rate",
  "summary": "Proposal to increase trust_deposit_yield_rate from 150000000000000000 to 160000000000000000",
  "messages": [
    {
      "@type": "/veranatest.td.v1.MsgUpdateParams",
      "authority": "cosmos1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsfwkgpd",
      "params": {
        "trust_deposit_share_value": "1000000000000000000",
        "trust_deposit_yield_rate": "160000000000000000"
      }
    }
  ]
}
EOF
```

#### 3.3: Proposal JSON Field Explanations

| Field | Description | Example |
|-------|-------------|---------|
| `group_policy_address` | The group policy that governs this action | Your group policy address |
| `proposers` | Array of addresses submitting the proposal | [council-member-1 address] |
| `metadata` | Human-readable description | "Update yield rate..." |
| `title` | Short proposal title | "Update Trust Deposit Yield Rate" |
| `summary` | Detailed explanation | What's being changed and why |
| `messages` | Array of messages to execute | MsgUpdateParams |
| `@type` | Message type (proto path) | `/veranatest.td.v1.MsgUpdateParams` |
| `authority` | Who's authorized to update | Group policy address |
| `params` | **ALL** module parameters | Complete params object |

⚠️ **IMPORTANT:** You must include **ALL** parameters, not just the ones you're changing!

---

### Step 4: Submit Proposal

#### 4.1: Submit the Proposal Transaction

```bash
veranatestd tx group submit-proposal update_td_params_proposal.json \
  --from council-member-1 \
  --chain-id vna-testnet-1 \
  --keyring-backend test \
  --fees 500000uvna \
  --node http://localhost:26657 \
  -y
```

**Expected Output:**
```
code: 0
txhash: DCE547A88E62A77ABAA27766671E3EC2320BD1F7108B9DBD1AB5FF45178122F7
```

✅ `code: 0` means success!

#### 4.2: Get the Proposal ID

Wait a few seconds for the transaction to be included in a block, then query:

```bash
veranatestd query group proposals-by-group-policy \
  cosmos1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsfwkgpd \
  --node http://localhost:26657 \
  -o json | jq -r '.proposals | sort_by(.id) | last | "Proposal ID: \(.id)\nTitle: \(.title)\nStatus: \(.status)"'
```

**Expected Output:**
```
Proposal ID: 2
Title: Update Trust Deposit Yield Rate
Status: PROPOSAL_STATUS_SUBMITTED
```

**Note the Proposal ID** - you'll need it for voting and execution (in this example: `2`).

---

### Step 5: Council Votes on Proposal

Each council member votes on the proposal. For a 2/3 threshold, you need at least 2 YES votes out of 3 members.

#### 5.1: Council Member 1 Votes YES

```bash
veranatestd tx group vote 2 cosmos17hl5uxaglku7a5n3xygu6yk98we6r6t608km6s VOTE_OPTION_YES "Approve yield rate increase" \
  --from council-member-1 \
  --chain-id vna-testnet-1 \
  --keyring-backend test \
  --fees 500000uvna \
  --node http://localhost:26657 \
  -y
```

**Command Breakdown:**
- `2` - Proposal ID
- `cosmos17hl5uxaglku7a5n3xygu6yk98we6r6t608km6s` - Voter address (council-member-1)
- `VOTE_OPTION_YES` - Vote choice (YES/NO/ABSTAIN/VETO)
- `"Approve yield rate increase"` - Voting rationale/comment

**Expected Output:**
```
code: 0
txhash: DF507FC804CB34E7C0FA292531728BBAD5177C803A22F912C824F122D8249927
```

#### 5.2: Council Member 2 Votes YES

Wait ~6 seconds for previous vote to finalize, then:

```bash
veranatestd tx group vote 2 cosmos1gdmwmz2qyatta9ymwc84wcj0syu3ruljnprlal VOTE_OPTION_YES "Approve yield rate increase" \
  --from council-member-2 \
  --chain-id vna-testnet-1 \
  --keyring-backend test \
  --fees 500000uvna \
  --node http://localhost:26657 \
  -y
```

**Expected Output:**
```
code: 0
txhash: 6BE196653F624CBF2D6FE20C1B4CE7FA7FD5267CEF2BB1431D0E0F872EC09EDC
```

✅ **2/3 threshold reached!** Proposal can now be executed.

#### 5.3: Verify Votes (Optional)

```bash
veranatestd query group votes-by-proposal 2 \
  --node http://localhost:26657 \
  -o json | jq -r '.votes[] | "Voter: \(.voter)\nOption: \(.option)\n---"'
```

**Expected Output:**
```
Voter: cosmos17hl5uxaglku7a5n3xygu6yk98we6r6t608km6s
Option: VOTE_OPTION_YES
---
Voter: cosmos1gdmwmz2qyatta9ymwc84wcj0syu3ruljnprlal
Option: VOTE_OPTION_YES
---
```

#### 5.4: Vote Options Reference

| Vote Option | Description | When to Use |
|-------------|-------------|-------------|
| `VOTE_OPTION_YES` | Approve the proposal | Support the parameter change |
| `VOTE_OPTION_NO` | Reject the proposal | Oppose the parameter change |
| `VOTE_OPTION_ABSTAIN` | Abstain from voting | Count towards quorum but not yes/no |
| `VOTE_OPTION_NO_WITH_VETO` | Strong rejection | Veto the proposal |

---

### Step 6: Execute Approved Proposal

Once the threshold is met (2/3 in our case), **anyone** can execute the proposal:

```bash
veranatestd tx group exec 2 \
  --from council-member-1 \
  --chain-id vna-testnet-1 \
  --keyring-backend test \
  --fees 500000uvna \
  --node http://localhost:26657 \
  -y
```

**Expected Output:**
```
code: 0
txhash: 4739D7D8F8696B2F0B5E14C1206C56251A1FA2C158E5844F22428F0698B20693
```

✅ Proposal executed! Parameters are now updated on-chain.

---

### Step 7: Verify Parameter Changes

Wait ~8 seconds for execution to finalize, then check the updated parameters:

```bash
veranatestd q td params --node http://localhost:26657
```

**BEFORE:**
```yaml
params:
  trust_deposit_share_value: "1000000000000000000"
  trust_deposit_yield_rate: "150000000000000000"
```

**AFTER:**
```yaml
params:
  trust_deposit_share_value: "1000000000000000000000000000000000000"
  trust_deposit_yield_rate: "160000000000000000000000000000000000"
```

✅ **SUCCESS!** The `trust_deposit_yield_rate` was updated from `150000000000000000` to `160000000000000000` via council governance!

---
