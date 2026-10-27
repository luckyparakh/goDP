package ocp

func init() {
	register("student", student)
}

func student(cents int64) int64 {
	return cents * 75 / 100
}
