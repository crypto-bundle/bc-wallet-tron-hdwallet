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
	"log"
	"testing"
)

func TestTronPubKey(t *testing.T) {
	// test tron
	prvMagic := [4]byte{0x04, 0x88, 0xad, 0xe4}
	pubMagic := [4]byte{0x04, 0x88, 0xb2, 0x1e}

	prvKey := "xprv9y95ZDyUd1t39HcvxQV9j5UqjhYRFsRXdoUrUVyDKCZo5fjvRDept9qLbHPg9B9jCn1vp9PJm575LCbkUyqv4ybRLRgzUV32DYbXVfMaS31"
	w, err := hdWalletFromString(prvKey, prvMagic, pubMagic)
	if err != nil {
		t.Errorf("%s should have been nil", err.Error())
	}
	pub, err := w.Pub().String()
	if err != nil {
		t.Error(err)
	}

	log.Println(pub)
	if pub != "xpub6C8RxjWNTPSLMmhQ4S2A6DRaHjNufL9P12QTGtNpsY6mxU54xky5Rx9pSZQp6B3gqPyfetfjwwPFeqPceuxyzm78HS4dAPYjFUEwdcpTqGF" {
		t.Errorf("\n%s\nsupposed to be\n", pub)
	}
}
