package context

type key string

func (k key) String() string {
	return "middleware.context.key." + string(k)
}

const (
	keyUserID key = "user_id"
	keyToken  key = "token"
	keyEmail  key = "email"
	keyRole   key = "role"
)

type contextKey string

const TxKey contextKey = "tx_key"
