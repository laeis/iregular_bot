package help

//CheckErr check if err not nil call panic
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}