package params

const (
	ProposalPeriod uint64 = 604800      // Seconds of a week, the protect time of a proposal, after which can propose a new one.
	InitOutTime    uint64 = 32503651200 // Time of 1/1/3000, for contract lock time check.

	DepositForProducer  uint64 = 1000000 * 1000000 * 1000000 // Deposit should be locked being a producer.
	LockTimeForDeposit  uint64 = 72 * 60 * 60                // Time be locked of deposit after producer unreg.
	ProducerSize        uint64 = 21                          // Default size of producers used in consensus
	MinimumProducerSize uint64 = 4                           // Minimum producer size

	StakeForVote     uint64 = 1000000 * 1000000 * 1000000 // Stake should be locked voting a producer/proxy.
	LockTimeForStake uint64 = 72 * 60 * 60                // Time be locked of stake after user unvote.
)
