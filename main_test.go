package main

import "testing"

type mathTest struct {
    arg1, arg2, expected float64
}

var addTests = []mathTest{
    mathTest{2, 3, 5},
    mathTest{4.2, 8.1, 12.3},
    mathTest{-6, 9, 3},
}

var subtractTests = []mathTest{
    mathTest{5, 3, 2},
    mathTest{13.2, 8.1, 5.1},
    mathTest{-6, 9, -15},
}

var multiplyTests = []mathTest{
    mathTest{6, 3, 18},
    mathTest{13.2, 8.1, 106.92},
    mathTest{-6, 9, -54},
}

var divideTests = []mathTest{
    mathTest{6, 3, 2},
    mathTest{13.2, 8.1, 1.63},
    mathTest{-6, 9, -0.67},
}

func TestAdd(t *testing.T){

    for _, test := range addTests{
        if output := Add(test.arg1, test.arg2); output != test.expected {
            t.Errorf("Output %f not equal to expected %f", output, test.expected)
        }
    }
}

func TestSubtract(t *testing.T){

    for _, test := range subtractTests{
        if output := Subtract(test.arg1, test.arg2); output != test.expected {
            t.Errorf("Output %f not equal to expected %f", output, test.expected)
        }
    }
}

func TestMultiply(t *testing.T){

    for _, test := range multiplyTests{
        if output := Multiply(test.arg1, test.arg2); output != test.expected {
            t.Errorf("Output %f not equal to expected %f", output, test.expected)
        }
    }
}

func TestDivide(t *testing.T){

    for _, test := range divideTests{
        if output := Divide(test.arg1, test.arg2); output != test.expected {
            t.Errorf("Output %f not equal to expected %f", output, test.expected)
        }
    }
}