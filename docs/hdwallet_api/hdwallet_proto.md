# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [hdwallet_api.proto](#hdwallet_api.proto)
    - [AddNewWalletRequest](#hdwallet_api.AddNewWalletRequest)
    - [AddNewWalletResponse](#hdwallet_api.AddNewWalletResponse)
    - [DerivationAddressByRangeRequest](#hdwallet_api.DerivationAddressByRangeRequest)
    - [DerivationAddressByRangeResponse](#hdwallet_api.DerivationAddressByRangeResponse)
    - [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity)
    - [DerivationAddressRequest](#hdwallet_api.DerivationAddressRequest)
    - [DerivationAddressResponse](#hdwallet_api.DerivationAddressResponse)
    - [GetEnabledWalletsRequest](#hdwallet_api.GetEnabledWalletsRequest)
    - [GetEnabledWalletsResponse](#hdwallet_api.GetEnabledWalletsResponse)
    - [GetWalletInfoRequest](#hdwallet_api.GetWalletInfoRequest)
    - [GetWalletInfoResponse](#hdwallet_api.GetWalletInfoResponse)
    - [MnemonicWalletData](#hdwallet_api.MnemonicWalletData)
    - [MnemonicWalletIdentity](#hdwallet_api.MnemonicWalletIdentity)
    - [RangeRequestUnit](#hdwallet_api.RangeRequestUnit)
    - [SignTransactionRequest](#hdwallet_api.SignTransactionRequest)
    - [SignTransactionResponse](#hdwallet_api.SignTransactionResponse)
    - [WalletBookmarks](#hdwallet_api.WalletBookmarks)
    - [WalletData](#hdwallet_api.WalletData)
    - [WalletIdentity](#hdwallet_api.WalletIdentity)
  
    - [WalletMakerStrategy](#hdwallet_api.WalletMakerStrategy)
  
    - [HdWalletApi](#hdwallet_api.HdWalletApi)
  
- [Scalar Value Types](#scalar-value-types)



<a name="hdwallet_api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## hdwallet_api.proto



<a name="hdwallet_api.AddNewWalletRequest"></a>

### AddNewWalletRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Title | [string](#string) |  |  |
| Purpose | [string](#string) |  |  |
| Strategy | [WalletMakerStrategy](#hdwallet_api.WalletMakerStrategy) |  |  |






<a name="hdwallet_api.AddNewWalletResponse"></a>

### AddNewWalletResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Wallet | [WalletData](#hdwallet_api.WalletData) |  |  |






<a name="hdwallet_api.DerivationAddressByRangeRequest"></a>

### DerivationAddressByRangeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [WalletIdentity](#hdwallet_api.WalletIdentity) |  |  |
| MnemonicIdentity | [MnemonicWalletIdentity](#hdwallet_api.MnemonicWalletIdentity) |  |  |
| Ranges | [RangeRequestUnit](#hdwallet_api.RangeRequestUnit) | repeated |  |






<a name="hdwallet_api.DerivationAddressByRangeResponse"></a>

### DerivationAddressByRangeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [WalletIdentity](#hdwallet_api.WalletIdentity) |  |  |
| MnemonicIdentity | [MnemonicWalletIdentity](#hdwallet_api.MnemonicWalletIdentity) |  |  |
| AddressIdentitiesCount | [uint64](#uint64) |  |  |
| AddressIdentities | [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity) | repeated |  |






<a name="hdwallet_api.DerivationAddressIdentity"></a>

### DerivationAddressIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AccountIndex | [uint32](#uint32) |  |  |
| InternalIndex | [uint32](#uint32) |  |  |
| AddressIndex | [uint32](#uint32) |  |  |
| Address | [string](#string) |  |  |






<a name="hdwallet_api.DerivationAddressRequest"></a>

### DerivationAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [WalletIdentity](#hdwallet_api.WalletIdentity) |  |  |
| MnemonicIdentity | [MnemonicWalletIdentity](#hdwallet_api.MnemonicWalletIdentity) |  |  |
| AddressIdentity | [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity) |  |  |






<a name="hdwallet_api.DerivationAddressResponse"></a>

### DerivationAddressResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [WalletIdentity](#hdwallet_api.WalletIdentity) |  |  |
| MnemonicIdentity | [MnemonicWalletIdentity](#hdwallet_api.MnemonicWalletIdentity) |  |  |
| AddressIdentity | [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity) |  |  |






<a name="hdwallet_api.GetEnabledWalletsRequest"></a>

### GetEnabledWalletsRequest







<a name="hdwallet_api.GetEnabledWalletsResponse"></a>

### GetEnabledWalletsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletsCount | [uint32](#uint32) |  |  |
| Wallets | [WalletData](#hdwallet_api.WalletData) | repeated |  |






<a name="hdwallet_api.GetWalletInfoRequest"></a>

### GetWalletInfoRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [WalletIdentity](#hdwallet_api.WalletIdentity) |  |  |






<a name="hdwallet_api.GetWalletInfoResponse"></a>

### GetWalletInfoResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [WalletIdentity](#hdwallet_api.WalletIdentity) |  |  |
| WalletInfo | [WalletData](#hdwallet_api.WalletData) |  |  |






<a name="hdwallet_api.MnemonicWalletData"></a>

### MnemonicWalletData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Identity | [MnemonicWalletIdentity](#hdwallet_api.MnemonicWalletIdentity) |  |  |
| IsHot | [bool](#bool) |  |  |






<a name="hdwallet_api.MnemonicWalletIdentity"></a>

### MnemonicWalletIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletUUID | [string](#string) |  |  |
| WalletHash | [string](#string) |  |  |






<a name="hdwallet_api.RangeRequestUnit"></a>

### RangeRequestUnit



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AccountIndex | [uint32](#uint32) |  |  |
| InternalIndex | [uint32](#uint32) |  |  |
| AddressIndexFrom | [uint32](#uint32) |  |  |
| AddressIndexTo | [uint32](#uint32) |  |  |






<a name="hdwallet_api.SignTransactionRequest"></a>

### SignTransactionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletUUID | [string](#string) |  |  |
| MnemonicWalletUUID | [string](#string) |  |  |
| AddressIdentity | [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity) |  |  |
| CreatedTxData | [bytes](#bytes) |  |  |






<a name="hdwallet_api.SignTransactionResponse"></a>

### SignTransactionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletIdentity | [WalletIdentity](#hdwallet_api.WalletIdentity) |  |  |
| MnemonicIdentity | [MnemonicWalletIdentity](#hdwallet_api.MnemonicWalletIdentity) |  |  |
| TxOwnerIdentity | [DerivationAddressIdentity](#hdwallet_api.DerivationAddressIdentity) |  |  |
| SignedTxData | [bytes](#bytes) |  |  |






<a name="hdwallet_api.WalletBookmarks"></a>

### WalletBookmarks



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| HotWalletIndex | [uint32](#uint32) |  |  |
| NonHotWalletIndexes | [uint32](#uint32) | repeated |  |






<a name="hdwallet_api.WalletData"></a>

### WalletData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Identity | [WalletIdentity](#hdwallet_api.WalletIdentity) |  |  |
| Title | [string](#string) |  |  |
| Purpose | [string](#string) |  |  |
| Strategy | [WalletMakerStrategy](#hdwallet_api.WalletMakerStrategy) |  |  |
| MnemonicWalletCount | [uint32](#uint32) |  |  |
| Bookmarks | [WalletBookmarks](#hdwallet_api.WalletBookmarks) |  |  |
| MnemonicWallets | [MnemonicWalletData](#hdwallet_api.MnemonicWalletData) | repeated |  |






<a name="hdwallet_api.WalletIdentity"></a>

### WalletIdentity



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| WalletUUID | [string](#string) |  |  |





 


<a name="hdwallet_api.WalletMakerStrategy"></a>

### WalletMakerStrategy


| Name | Number | Description |
| ---- | ------ | ----------- |
| PLACEHOLDER_STRATEGY_TYPE | 0 |  |
| SINGLE_MNEMONIC_STRATEGY_TYPE | 1 |  |
| MULTIPLE_MNEMONIC_STRATEGY_TYPE | 2 |  |


 

 


<a name="hdwallet_api.HdWalletApi"></a>

### HdWalletApi


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddNewWallet | [AddNewWalletRequest](#hdwallet_api.AddNewWalletRequest) | [AddNewWalletResponse](#hdwallet_api.AddNewWalletResponse) |  |
| GetWalletInfo | [GetWalletInfoRequest](#hdwallet_api.GetWalletInfoRequest) | [GetWalletInfoResponse](#hdwallet_api.GetWalletInfoResponse) |  |
| GetEnabledWallets | [GetEnabledWalletsRequest](#hdwallet_api.GetEnabledWalletsRequest) | [GetEnabledWalletsResponse](#hdwallet_api.GetEnabledWalletsResponse) |  |
| GetDerivationAddress | [DerivationAddressRequest](#hdwallet_api.DerivationAddressRequest) | [DerivationAddressResponse](#hdwallet_api.DerivationAddressResponse) |  |
| GetDerivationAddressByRange | [DerivationAddressByRangeRequest](#hdwallet_api.DerivationAddressByRangeRequest) | [DerivationAddressByRangeResponse](#hdwallet_api.DerivationAddressByRangeResponse) |  |
| SignTransaction | [SignTransactionRequest](#hdwallet_api.SignTransactionRequest) | [SignTransactionResponse](#hdwallet_api.SignTransactionResponse) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

