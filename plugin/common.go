/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

const tronPluginName = "tron-hdwallet-plugin"

// DO NOT EDIT THESE VARIABLES DIRECTLY. These are build-time constants
// DO NOT USE THESE VARIABLES IN APPLICATION CODE. USE commonConfig.NewLdFlagsManager SERVICE-COMPONENT INSTEAD OF IT
var (
	// ReleaseTag - release tag in TAG.SHORT_COMMIT_ID.BUILD_NUMBER.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ReleaseTag = "v0.0.0-00000000-100500"

	// CommitID - latest commit id.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	CommitID = "0000000000000000000000000000000000000000"

	// ShortCommitID - first 12 characters from CommitID.
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	ShortCommitID = "0000000"

	// BuildNumber - ci/cd build number for BuildNumber
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildNumber string = "100500"

	// BuildDateTS - ci/cd build date in time stamp
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	BuildDateTS string = "1713280105"
	// NetworkChainID - blockchain network ID, Tron blockchain HdWallet coinID = 195
	// DO NOT EDIT THIS VARIABLE DIRECTLY. These are build-time constants
	// DO NOT USE THESE VARIABLES IN APPLICATION CODE
	NetworkChainID = "195"
)

var (
	pluginChainID = TronCoinNumber

	ErrUnsupportedCoinID = errors.New("unsupported coin id value")
)

var setChainIDOnce = sync.Once{}

func init() {
	setChainIDOnce.Do(func() {
		if NetworkChainID == "" {
			pluginChainID = TronCoinNumber

			return
		}

		chainIDInt, err := strconv.Atoi(NetworkChainID)
		if err != nil {
			panic(fmt.Errorf("wrong network chainID format: %w", err))
		}

		switch chainIDInt {
		case TronCoinNumber:
			pluginChainID = chainIDInt

			err = nil
		default:
			err = ErrUnsupportedCoinID
		}

		if err != nil {
			panic(err)
		}

		return
	})
}

func GetPluginName() string {
	return tronPluginName
}

func GetPluginReleaseTag() string {
	return ReleaseTag
}

func GetPluginCommitID() string {
	return CommitID
}

func GetPluginShortCommitID() string {
	return ShortCommitID
}

func GetPluginBuildNumber() string {
	return BuildNumber
}

func GetPluginBuildDateTS() string {
	return BuildDateTS
}

func GetSupportedChainIDs() []int {
	return []int{TronCoinNumber}
}

func GetChainID() int {
	return pluginChainID
}
