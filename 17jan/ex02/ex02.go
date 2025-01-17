package ex02

// Factorial
// O fatorial de um número é o resultado da multiplicação desse número por todos
// os seus antecessores maiores que zero
// Por exemplo, o fatorial de 4 é 4 × 3 × 2 × 1, que é igual a 24.
func Factorial(i int) int {
	if i < 0 {
		return 0
	}

	// base
	if i == 0 || i == 1 {
		return 1
	}

	// if i == 2 {
	// 	return 2
	// }
	//
	// if i == 3 {
	// 	return 6
	// }
	//
	// if i == 4 {
	// 	return 24
	// }

	// Recursão
	return i * Factorial(i-1)
}
