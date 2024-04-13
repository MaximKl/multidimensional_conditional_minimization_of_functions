package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type result struct {
	name   string
	points [][3]float64
}

type sphere struct {
	a float64
	b float64
	r float64
}

func calcFunction(x1, x2 float64) float64 {
	return math.Pow(x1, 2) - x2
}

func firstDerivativeX1(x1 float64) float64 {
	return 2 * x1
}

func firstDerivativeX2() float64 {
	return -1
}

func isAccurate(f, fPrev, e float64) bool {
	return math.Abs(f-fPrev) <= e
}

func (s *sphere) calcX1(z1, z2 float64) float64 {
	return s.a + (z1-s.a)/(math.Sqrt(math.Pow(z1-s.a, 2)+math.Pow(z2-s.b, 2)))*s.r
}

func (s *sphere) calcX2(z1, z2 float64) float64 {
	return s.b + (z2-s.b)/(math.Sqrt(math.Pow(z1-s.a, 2)+math.Pow(z2-s.b, 2)))*s.r
}

func main() {
	fmt.Println("Current function: x1^2 - x2")
	sphere := sphere{a: 5.0, b: 2.0, r: 2.0}
	reader := bufio.NewReader(os.Stdin)
	gradientReceiver := make(chan result)
	stop := ""

	for stop != "n\r\n" && stop != "N\r\n" {
		x1, x2, e, err := readUserInput(reader)
		if err != nil {
			fmt.Println(err)
			continue
		}

		points := [][3]float64{{x1, x2, calcFunction(x1, x2)}}
		go gradientProjectionMethod(e, sphere, points, gradientReceiver, 1000)

		for runtime.NumGoroutine() != 1 {
			select {
			case result := <-gradientReceiver:
				getBestPoint("Gradient Projection method", result)
				writePoints(fmt.Sprintf("Gradient_Projection_method%s", result.name), result.points)
			}
		}
		fmt.Print("Continue? [y/n]: ")
		stop, _ = reader.ReadString('\n')
	}
}

func gradientProjectionMethod(e float64, s sphere, points [][3]float64, receiver chan result, maxIter int) {
	fPrev := 0.0
	for k := 0; (k < maxIter) && !isAccurate(points[k][2], fPrev, e); k++ {
		z1 := points[k][0] - firstDerivativeX1(points[k][0])
		z2 := points[k][1] - firstDerivativeX2()

		x1 := s.calcX1(z1, z2)
		x2 := s.calcX2(z1, z2)
		fPrev = points[k][2]
		points = append(points, [3]float64{x1, x2, calcFunction(x1, x2)})
	}
	receiver <- result{
		name:   fmt.Sprintf("(%v,%v)", points[0][0], points[0][1]),
		points: points}
}

func readUserInput(reader *bufio.Reader) (float64, float64, float64, error) {
	fmt.Print("Enter X1: ")
	x1, _ := reader.ReadString('\n')
	fmt.Print("Enter X2: ")
	x2, _ := reader.ReadString('\n')
	fmt.Print("Enter Accuracy: ")
	e, _ := reader.ReadString('\n')

	x1f, err1 := strconv.ParseFloat(replace(x1), 64)
	x2f, err2 := strconv.ParseFloat(replace(x2), 64)
	ef, err3 := strconv.ParseFloat(replace(e), 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, fmt.Errorf("wrong input")
	}
	return x1f, x2f, ef, nil
}

func replace(s string) string {
	return strings.Replace(s, "\r\n", "", -1)
}

func getBestPoint(methodName string, result result) {
	fmt.Printf("----Best results of %s-----\n", methodName)
	fmt.Printf("%s with starting point%s X1 and X2: (%v, %v)\n", methodName, result.name, result.points[len(result.points)-1][0], result.points[len(result.points)-1][1])
	fmt.Printf("%s with starting point%s F: %g\n", methodName, result.name, result.points[len(result.points)-1][2])
	fmt.Printf("%s with starting point%s K: %v\n", methodName, result.name, len(result.points)-1)
}

func writePoints(fileName string, points [][3]float64) {
	stringPoints := ""
	fullFileName := fmt.Sprintf("output/%s.txt", fileName)
	for _, p := range points {
		stringPoints = fmt.Sprintf("%s(%v, %v) | %v\n", stringPoints, p[0], p[1], p[2])
	}
	err := os.WriteFile(fullFileName, []byte(stringPoints), 0777)
	if err != nil {
		fmt.Println("Error while was writing into a file")
		return
	}
	fmt.Printf("All intermediate results have been successfully written to %s\n\n", fullFileName)
}
