package plugin

import "time"

type walletMakerFunc func(walletUUID string,
	mnemonicDecryptedData string,
) (interface{}, error)

type ldFlagManager interface {
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() int64
	GetBuildDate() time.Time
}
