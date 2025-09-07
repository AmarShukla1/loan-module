package constants

import "time"

const Workers = 5
const ChannelBufferSize = 100

const TimeIntervalToFeedJobs = 5 * time.Second
const MinAmountApproveBySystem = 10000
const MaxAmountApproveBySystem = 500000

const DefaultPage = 1
const DefaultPageSize = 10
const DefaultMinPage = 1
const DefaultMaxPageSize = 10
const DefaultMinPageSize = 1
