package cpf

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

func cryptoRandInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

// UnformatCPF removes all non-digit characters from the input string.
func UnformatCPF(cpfStr string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(cpfStr, "")
}

// IsRepeated checks if the string is composed entirely of the same character.
func IsRepeated(s string) bool {
	if len(s) == 0 {
		return false
	}
	firstChar := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != firstChar {
			return false
		}
	}
	return true
}

// FormatCPF formats an 11-digit CPF string as ###.###.###-##.
func FormatCPF(cpfStr string) (string, error) {
	digits := UnformatCPF(cpfStr)
	if len(digits) != 11 {
		return "", fmt.Errorf("invalid CPF number (must have 11 digits)")
	}
	return fmt.Sprintf("%s.%s.%s-%s",
		digits[0:3],
		digits[3:6],
		digits[6:9],
		digits[9:11],
	), nil
}

// calc implements the internal sum with the CPF weighting logic.
func calc(nums []int) int {
	total := 0
	for i, num := range nums {
		total += num * (9 - (i % 10))
	}
	return total
}

// getCD computes the 2 check digits (DV) from the first 9 digits of a CPF.
func getCD(digits9 []int) ([2]int, error) {
	if len(digits9) != 9 {
		return [2]int{}, fmt.Errorf("invalid digits length: expected 9, got %d", len(digits9))
	}

	reversed := make([]int, 9)
	for i := 0; i < 9; i++ {
		reversed[i] = digits9[9-1-i]
	}

	cd1 := calc(reversed) % 11 % 10

	secondInput := append([]int{0}, reversed...)
	cd2 := (calc(secondInput) + cd1*9) % 11 % 10

	return [2]int{cd1, cd2}, nil
}

// ValidateCPF checks if the provided CPF string is valid.
func ValidateCPF(cpfStr string, byLength bool) bool {
	unformatted := UnformatCPF(cpfStr)
	if len(unformatted) != 11 {
		return false
	}
	if IsRepeated(unformatted) {
		return false
	}

	if byLength {
		return true
	}

	number9 := unformatted[:9]
	dv2 := unformatted[9:11]

	var digits9 []int
	for _, ch := range number9 {
		d, err := strconv.Atoi(string(ch))
		if err != nil {
			return false
		}
		digits9 = append(digits9, d)
	}

	cd, err := getCD(digits9)
	if err != nil {
		return false
	}
	trueDV := fmt.Sprintf("%d%d", cd[0], cd[1])

	return dv2 == trueDV
}

// GenerateCPF creates a random CPF number.
func GenerateCPF(formatted, invalid bool) (string, error) {
	digits9 := make([]int, 9)
	for i := 0; i < 9; i++ {
		digit, err := cryptoRandInt(10)
		if err != nil {
			return "", fmt.Errorf("failed to generate random digit: %w", err)
		}
		digits9[i] = digit
	}

	var dv [2]int
	if invalid {
		d1, err := cryptoRandInt(10)
		if err != nil {
			return "", fmt.Errorf("failed to generate random digit: %w", err)
		}
		d2, err := cryptoRandInt(10)
		if err != nil {
			return "", fmt.Errorf("failed to generate random digit: %w", err)
		}
		dv[0] = d1
		dv[1] = d2
	} else {
		correctDV, err := getCD(digits9)
		if err != nil {
			return "", fmt.Errorf("failed to generate check digits: %w", err)
		}
		dv = correctDV
	}

	allDigits := append(digits9, dv[0], dv[1])
	cpfStr := strings.Builder{}
	for _, d := range allDigits {
		cpfStr.WriteString(strconv.Itoa(d))
	}

	if formatted {
		result, err := FormatCPF(cpfStr.String())
		if err != nil {
			return "", fmt.Errorf("failed to format CPF: %w", err)
		}
		return result, nil
	}
	return cpfStr.String(), nil
} 