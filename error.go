package rocket

type PageNotFound string

func (e PageNotFound) Error() string {
	return string(e)
}
