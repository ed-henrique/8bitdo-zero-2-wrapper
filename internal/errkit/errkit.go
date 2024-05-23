package errkit

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
