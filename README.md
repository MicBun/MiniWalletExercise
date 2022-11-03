# MiniWalletExercise
a Mini Wallet Exercise

## Initialize my account for wallet
`POST` `http://localhost/api/v1/init`

Example: BODY | formdata
```
customer_xid ea0212d3-abd6-406f-8c67-868e814a2436
```

Before using the wallet service, an account needs to be created. This endpoint is used to create an account as well as getting the token for the other API endpoints.
## Enable my wallet
`POST` `http://localhost/api/v1/wallet`

As a customer, I can enable my wallet using this API endpoint. The wallet stores the virtual money that customers can apply for approval. The virtual money can be used to make purchases with our partners, who are various online merchants.
If the wallet is already enabled, this endpoint would fail. This endpint should be usable again only if the wallet is disabled.
Before enabling the wallet, the customer cannot view, add, or use its virtual money.

Example:
HEADERS | Authorization

`Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238`

*Token comes from `Initialize my account for wallet`

## View my wallet balance
`GET` `http://localhost/api/v1/wallet`

With this API endpoint, a customer can view the current balance of virtual money. After adding or using virtual money, it is not expected to have the balance immediately updated. The maximum delay for updating the balance is 5 seconds.
Example: 
HEADERS | Authorization

`Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238`

## Add virtual money to my wallet
`POST` `http://localhost/api/v1/wallet/deposits`

With this API endpoint, a customer can add virtual money to the wallet balance as a deposit once the wallet is enabled. Reference ID passed must be unique for every deposit.

HEADERS | Authorization

`Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238`

BODY | formdata

`amount 100000`

`reference_id 50535246-dcb2-4929-8cc9-004ea06f5241`

## Use virtual money from my wallet
`POST` `http://localhost/api/v1/wallet/withdrawals`

With this API endpoint, a customer can use the virtual money from the wallet balance as a withdrawal once the wallet is enabled. The amount being used must not be more than the current balance. Reference ID passed must be unique for every withdrawal.

HEADERS | Authorization

`Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238`

BODY | formdata

`amount 60000`

`reference_id 4b01c9bb-3acd-47dc-87db-d9ac483d20b2`

## Disable my wallet

`PATCH` `http://localhost/api/v1/wallet`

With this API endpoint, a customer's wallet can be disabled as determined by the service. Once disabled, the customer cannot view, add, or use its virtual money.

HEADERS | Authorization

`Token 6b3f7dc70abe8aed3e56658b86fa508b472bf238`

BODY | formdata

`is_disabled true`