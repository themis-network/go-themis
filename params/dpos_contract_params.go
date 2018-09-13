package params

const (
	ProposalPeriod uint64 = 604800 // Seconds of a week, the protect time of a proposal, after which can propose a new one.

	DepositForProducer  uint64 = 1000000 * 1000000 * 1000000 // Deposit should be locked being a producer.
	LockTimeForDeposit  uint64 = 72 * 60 * 60                // Time be locked of deposit after producer unreg.
	ProducerSize        uint64 = 21                          // Default size of producers used in consensus
	MinimumProducerSize uint64 = 4                           // Minimum producer size
	MaxProducersSize    uint64 = 10000                       // Max producers size

	StakeForVote     uint64 = 1000000 * 1000000 * 1000000 // Stake should be locked voting a producer/proxy.
	LockTimeForStake uint64 = 72 * 60 * 60                // Time be locked of stake after user unvote.
)
