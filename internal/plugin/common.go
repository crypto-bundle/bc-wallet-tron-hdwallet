package plugin

import "time"

type generateMnemonicFunc func() (string, error)
type validateMnemonicFunc func(mnemonic string) bool

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
