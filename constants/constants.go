package constants

import "time"

const Workers = 5
const ChannelBufferSize = 100

const TimeIntervalToFeedJobs = 5 * time.Second
const MinAmountApproveBySystem = 10000
const MaxAmountApproveBySystem = 500000
