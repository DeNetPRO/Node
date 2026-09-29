# How Datakeeper Earnings Work

**Last update:** 2026-09-29

Datakeeper rewards are not minted from nothing — they are exactly the fees storage users pay, distributed by stored volume. This guide follows a single coin through the protocol: from a user's payment to your wallet.

## Table of Contents

- [Where the Money Comes From](#where-the-money-comes-from)
- [The Deposit](#the-deposit)
- [The Reward Jar](#the-reward-jar)
- [The Storage Cycle](#the-storage-cycle)
- [How One Cycle Pays](#how-one-cycle-pays)
- [Why Payouts Are Gradual](#why-payouts-are-gradual)
- [A Worked Example](#a-worked-example)
- [Where Your Rewards Appear](#where-your-rewards-appear)
- [Verify Your Payouts (Coming Soon)](#verify-your-payouts-coming-soon)
- [FAQ](#faq)

## Where the Money Comes From

A storage user pays in **TBY tokens** for keeping their files safe. Datakeeper nodes do the storing, grouped into **pools** of up to 32 nodes; each user's data lives in one specific pool. Your node is paid by the users of *its* pool — that single fact explains almost all variation in earnings you will ever see.

## The Deposit

When a user pays, the tokens move into the pool's smart contract and stay there as the **user's balance**. Nothing is paid to nodes at this moment. The balance is a ceiling: it marks how much the protocol is allowed to charge this user over time, and anything not yet charged stays refundable at any time.

## The Reward Jar

Money is transferred from user balances into the pool's shared **reward jar** in *draws* — small charges that happen as the user's data occupies storage:

| Event | What happens |
|---|---|
| User uploads / deletes files | A draw runs along with the node's report |
| User tops up or claims funds back | A draw runs |
| One week passes with no activity | An automatic draw runs |
| User's balance runs out | The draw takes the remainder; the data is no longer funded storage |

Storage is billed on elapsed time at **1 TBY per terabyte per year**, with one rule: accounts below 100 GB are billed as 100 GB. Nothing is charged between draws — if a draw is skipped, the next one simply collects the whole passed period at once.

Once money crosses into the jar, it is no longer anybody's balance — the cycles ahead will distribute it to the Datakeepers.

## The Storage Cycle

A pool works in repeated **cycles**. Each cycle has four stages:

1. The node reports which users uploaded or removed their data.
2. The node reports how many files it stores.
3. The network waits for all reports to settle.
4. The node proves the stored files are still intact — payouts happen, and a new cycle starts.

A node that fails to prove its data in a cycle gets nothing for that cycle and takes one **penalty**. Ten penalties in a row mean leaving the pool — its share of the distribution goes to the remaining nodes. One successful proof resets the counter. **Money already earned is never taken away.**

## How One Cycle Pays

At the end of each cycle, the pool pays out **1/30 of its jar**, split among the nodes that proved their data:

> **reward = (jar / 30) × (your confirmed MB / all confirmed MB in the pool)**

Two numbers decide everything:

- **How much the jar gives up this cycle** — set by what the pool's users are paying.
- **Your share of the pool's confirmed data** — *confirmed* means reported and proven in the current cycle.

Note what is *not* in the formula: the size of your disk. You are paid for the data you actually store and prove, megabyte by megabyte — empty or purchased-but-unused space earns nothing.

## Why Payouts Are Gradual

The jar fills continuously, yet it pays out in slices of 1/30 per cycle. A deposit of 300 TBY entering the jar, for instance, pays out like this:

| Cycle | 1 | 2 | 5 | 10 | 30 | 100 |
|---|---|---|---|---|---|---|
| Paid out of this jar, TBY | 10.0 | 9.7 | 8.8 | 7.4 | 3.6 | 0.3 |
| Paid out in total, TBY | 10 | 20 | 47 | 86 | 192 | 290 |

This shape does three things:

- **Smooths income.** One user top-up keeps paying the pool's nodes for tens of cycles, so nodes do not live from deposit to deposit.
- **Prevents exit scams.** A node collects its share only cycle by cycle, and only while its proofs keep passing — grabbing the money and disappearing is impossible.
- **Creates short ripples around a stable average.** A wave of user top-ups lifts payouts briefly and a quiet week dents them, but over any longer stretch what a node earns tracks what its pool's users pay.

## A Worked Example

The protocol runs a cycle roughly every 75 minutes (cycle length is a protocol setting and gets adjusted), so about **19 cycles a day**. Data is kept in at least 3 copies, so a pool storing 10 TB of user data confirms about 30 TB.

**Pool A** — 10 TB of data from many users:

| Account | Occupied | Billed as |
|---|---|---|
| Alice | 4 TB | 4 TB |
| Bob | 4 TB | 4 TB |
| Anna | 1 TB | 1 TB |
| 20 small accounts × 50 GB | 1 TB | 2 TB (100 GB minimum each) |
| **Total** | **10 TB** | **11 TB** |

11 TBY a year flows into the jar ≈ 0.0016 TBY per cycle, and in steady state the pool pays the same ≈ 0.0016 TBY per cycle. A node that proves 2 of the pool's 30 confirmed TB gets 2/30 of each payout: **≈ 0.00011 TBY per cycle, ≈ 0.0020 TBY per day**.

**Pool B** — the same 10 TB, but from one user paying for exactly 10 TB. The same node now earns **≈ 0.0018 TBY per day** from identical hardware and identical data volume.

| | Pool A | Pool B |
|---|---|---|
| Data occupied | 10 TB | 10 TB |
| Confirmed (×3 copies) | 30 TB | 30 TB |
| Billed as | 11 TB | 10 TB |
| Users pay per day | 0.030 TBY | 0.027 TBY |
| Your node (2 TB proven) per day | 0.0020 TBY | 0.0018 TBY |

The takeaway: the same terabyte can pay slightly differently from pool to pool, because a pool's income is the sum of its users' bills. Neither you nor the protocol can steer this — it is decided by where the users happen to be. In every pool the two sides stay equal: users pay per day exactly what nodes receive.

## Where Your Rewards Appear

Each cycle the pool pays its share straight to the wallet that your license belongs to. In the **Datakeeper Console** this arrives as **Current Rewards**, which you can withdraw at any time.

## Verify Your Payouts (Coming Soon)

Everything described here lives on-chain: the pool's jar and the confirmed amounts are recorded, so payouts can in principle be re-checked for any past cycle. A comfortable interface for that doesn't exist yet — soon the **Datakeeper Console** will let you visually see and verify your earnings against the on-chain data yourself. Until then, the protocol itself counts, distributes and charges everything for you.

## FAQ

**Does a user's payment go to the nodes right away?**
No. The tokens sit in the pool contract as the user's balance and only cross into the reward jar in draws, strictly as the user's data occupies time.

**Is the "payment cycle" the same as the storage cycle?**
Yes — the protocol runs in cycles, accounts are drawn on events, and every cycle ends in one payout.

**Why do my rewards jump and sag while my stored data doesn't change?**
Because the jar drains by 1/30 per cycle, payouts mostly reflect the pool's recent deposits, and payouts are made per pool, not per network. Different days and different nodes get different rates even for identical bytes. A "network average" is a number no real node actually receives.

**Does my disk size affect my reward?**
No. Only data reported by your node and proven in the current cycle counts. Empty disk earns nothing.

**What happens if I miss a proof?**
You get nothing for that cycle and take one penalty; 10 penalties in a row remove you from the pool. One successful proof resets the counter, and money you already earned is never taken away.
