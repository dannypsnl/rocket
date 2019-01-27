package rocket

type Guard interface {
	VerifyRequest() error
}
