# How Datakeeper Earnings Work

**Last update:** 2026-09-29

Datakeeper rewards are not minted from nothing — they are exactly the fees storage users pay, distributed by stored volume. Rewards are paid in the same token a user paid with, credited by contract to the wallet the license belongs to. This guide follows a single payment through the protocol: from the moment a user makes it to the moment it lands in your wallet.

## Table of Contents

- [Where Rewards Come From](#where-rewards-come-from)
- [The Deposit](#the-deposit)
- [The Reward Fund](#the-reward-fund)
- [The Storage Cycle](#the-storage-cycle)
- [How One Cycle Pays](#how-one-cycle-pays)
- [Why Payouts Are Gradual](#why-payouts-are-gradual)
- [A Worked Example](#a-worked-example)
- [Where Your Rewards Appear](#where-your-rewards-appear)
- [Verify Your Payouts (Coming Soon)](#verify-your-payouts-coming-soon)
- [FAQ](#faq)

## Where Rewards Come From

A storage user pays in **TBY tokens** for keeping their files safe. TBY is a utility token for paying for network services and distributing rewards for the provided storage resources. Datakeeper nodes do the storing, grouped into **pools** of up to 32 nodes; each user's data lives in one specific pool. Your node is rewarded by the users of *its* pool — that single fact explains almost all variation in earnings you will ever see.

No bank or intermediary holds rewards: the rules that move them are written into smart contracts, which validate storage proofs and apply rewards and penalties automatically.

## The Deposit

When a user pays, the tokens move into the pool's smart contract and stay there as the **user's balance**. Nothing is paid to nodes at this moment. The balance is a ceiling: it marks how much the protocol is allowed to charge this user over time, and anything not yet charged stays refundable at any time.

## The Reward Fund

Fees are transferred from user balances into the pool's shared **reward fund** in *draws* when the user performs one of the following actions:

| Event | What happens |
|---|---|
| User uploads / deletes files | A draw runs along with the node's report |
| User tops up or claims funds back | A draw runs |
| User's balance runs out | The draw takes the remainder; the data is no longer funded storage |
| One week passes with no activity | An automatic draw runs |

Storage is billed on elapsed time at **1 TBY per terabyte per year**, with one rule: accounts below 100 GB are billed as 100 GB. Nothing is charged between draws — if a draw is skipped, the next one simply collects the whole passed period at once.

Once a fee crosses into the fund, it is no longer anybody's balance — the cycles ahead will distribute it to the Datakeepers.

## The Storage Cycle

A pool works in repeated **cycles**. Each cycle has four stages:

1. The node reports which users uploaded or removed their data.
2. The node reports how many files it stores.
3. The period in which the node continues to store user data.
4. The node proves the stored files are still intact — payouts happen, and a new cycle starts.

A node that fails to prove its data in a cycle gets nothing for that cycle and takes one **penalty**. Ten penalties in a row mean leaving the pool — its share of the distribution goes to the remaining nodes. One successful proof resets the counter. **Rewards already earned are never taken away.**

## How One Cycle Pays

At the end of each cycle, the pool pays out **1/30 of its fund**, split among the nodes that proved their data:

> **reward = (fund / 30) × (your confirmed MB / all confirmed MB in the pool)**

Two numbers decide everything:

- **How much the fund gives up this cycle** — set by what the pool's users are paying.
- **Your share of the pool's confirmed data** — *confirmed* means reported and proven in the current cycle.

Note what is *not* in the formula: the size of your disk. You are paid for the data you actually store and prove, megabyte by megabyte — empty or purchased-but-unused space earns nothing.

## Why Payouts Are Gradual

The fund fills continuously, yet it pays out in slices of 1/30 per cycle. A deposit of 300 TBY entering the fund, for instance, pays out like this:

| Cycle | 1 | 2 | 5 | 10 | 30 | 100 |
|---|---|---|---|---|---|---|
| Paid out of this fund, TBY | 10.0 | 9.7 | 8.8 | 7.4 | 3.6 | 0.3 |
| Paid out in total, TBY | 10 | 20 | 47 | 86 | 192 | 290 |

This shape does three things:

- **Smooths payouts.** One user top-up keeps paying the pool's nodes for tens of cycles, so nodes do not live from deposit to deposit.
- **Keeps payouts honest.** A node collects its share only cycle by cycle, and only while its proofs keep passing — being paid for storage that can no longer be proven is impossible.
- It smooths income instead of paying on demand. A deposit keeps feeding payouts for about 30 cycles (~1.5 days), so a wave of user top-ups lifts them briefly and a quiet week dents them. These are short-term ripples around a stable average: over any longer stretch, what a node earns tracks what its pool's users pay.

## A Worked Example

Let's assume the protocol runs a cycle every 90 minutes (cycle length is a protocol setting and gets adjusted), so about **16 cycles a day**. Data is kept in at least 3 copies, so a pool storing 10 TB of user data confirms about 30 TB.

**Pool A** — 10 TB of data from many users:

| Account | Occupied | Billed as |
|---|---|---|
| Alice | 4 TB | 4 TB |
| Bob | 4 TB | 4 TB |
| Anna | 1 TB | 1 TB |
| 20 small accounts × 50 GB | 1 TB | 2 TB (100 GB minimum each) |
| **Total** | **10 TB** | **11 TB** |

11 TBY a year flows into the fund ≈ 0.0019 TBY per cycle, and in steady state the pool pays the same ≈ 0.0019 TBY per cycle. A node that proves 2 of the pool's 30 confirmed TB gets 2/30 of each payout: **≈ 0.00013 TBY per cycle, ≈ 0.0020 TBY per day**.

These are averages: fees enter the fund in draws rather than every cycle, so the actual reward per cycle and per day varies around them.

**Pool B** — the same 10 TB, but from one user paying for exactly 10 TB. The same node now earns **≈ 0.0018 TBY per day** from identical hardware and identical data volume.

| | Pool A | Pool B |
|---|---|---|
| Data occupied | 10 TB | 10 TB |
| Confirmed (×3 copies) | 30 TB | 30 TB |
| Billed as | 11 TB | 10 TB |
| Users pay per day | 0.030 TBY | 0.027 TBY |
| Your node (2 TB proven) per day | 0.0020 TBY | 0.0018 TBY |

The takeaway: the same terabyte can pay slightly differently from pool to pool, because a pool's rewards are the sum of its users' fees. Neither you nor the protocol can steer this — it is decided by where the users happen to be. In every pool the two sides stay equal: users pay per day exactly what nodes receive.

## Where Your Rewards Appear

Each cycle the pool pays its share¹ straight to the wallet that your license belongs to. In the **Datakeeper Console** this arrives as **Current Rewards**, which you can withdraw at any time.

¹ Per the protocol rules a fee may be retained as part of this distribution, and the same fee applies when TBY is issued or redeemed; the payout figures in this guide are shown before any such fee.

## Verify Your Payouts (Coming Soon)

Reward accounting is executed on-chain: the pool's fund and the confirmed amounts are recorded in the contracts, so payouts can in principle be re-checked for any past cycle. A comfortable interface for that doesn't exist yet — soon the **Datakeeper Console** will let you see and recalculate your earnings from the on-chain data yourself. Until then, counting, charging and distributing are handled by the contract code rather than by hand.

## FAQ

**Does a user's payment go to the nodes right away?**
No. The tokens sit in the pool contract as the user's balance and only cross into the reward fund in draws, strictly as the user's data occupies time.

**Is the "payment cycle" the same as the storage cycle?**
Yes — the protocol runs in cycles, accounts are drawn on events, and every cycle ends in one payout.

**Why do my rewards jump and sag while my stored data doesn't change?**
Because the fund drains by 1/30 per cycle, payouts mostly reflect the pool's recent deposits, and payouts are made per pool, not per network. Different days and different nodes get different rates even for identical bytes. A "network average" is a number no real node actually receives.

**Does my disk size affect my reward?**
No. Only data reported by your node and proven in the current cycle counts. Empty disk earns nothing.

**What happens if I miss a proof?**
You get nothing for that cycle and take one penalty; 10 penalties in a row remove you from the pool. One successful proof resets the counter, and rewards you already earned are never taken away.

**Are rewards guaranteed?** 
No. They depend on the storage fees paid by users in your pool, your share of proven data and passed proofs; they can be low or zero, and the TBY/USD rate changes.
