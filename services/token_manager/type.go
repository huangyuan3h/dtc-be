package token_manager

type Token struct {
	TokenId    string `json:"tokenId" dynamodbav:"tokenId"`
	ExpireAt   string `json:"expireAt" dynamodbav:"expireAt"`
	ConsumedBy string `json:"consumedBy" dynamodbav:"consumedBy"`
	IsConsumed string `json:"isConsumed" dynamodbav:"isConsumed"`
}
