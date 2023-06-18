package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"
)

//go:embed embed/ac11.csv
var ac11 []byte

//go:embed embed/qc.csv
var qc []byte

func sq(n ...float64) float64 {
	sum := 0.0
	for _, i := range n {
		sum += math.Pow(i, 2)
	}

	return sum
}

func calcTowTheta(wl, n, lc float64) float64 {
	return math.Asin(wl*math.Sqrt(n)/2/lc) * 360 / math.Pi
}

func calcAcPeak(lc, wl float64) (*[]string, error) {
	var list []string
	r := csv.NewReader(bytes.NewReader(ac11))

	for {
		rec, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		h, _ := strconv.ParseFloat(rec[0], 64)
		k, _ := strconv.ParseFloat(rec[1], 64)
		l, _ := strconv.ParseFloat(rec[2], 64)

		N := sq(h, k, l)
		tt := calcTowTheta(wl, N, lc)

		list = append(list, fmt.Sprintf("(%s, %s, %s) ~%.2f:", rec[0], rec[1], rec[2], tt))
	}

	return &list, nil
}

func calcQcPeak(lc, wl float64) (*[]string, error) {
	var list []string
	r := csv.NewReader(bytes.NewReader(qc))

	for {
		rec, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		h, _ := strconv.ParseFloat(rec[0], 64)
		k, _ := strconv.ParseFloat(rec[1], 64)
		l, _ := strconv.ParseFloat(rec[2], 64)
		m, _ := strconv.ParseFloat(rec[3], 64)
		n, _ := strconv.ParseFloat(rec[4], 64)
		o, _ := strconv.ParseFloat(rec[5], 64)

		r1 := math.Sqrt(5)*h + k + l + m + n + o
		r2 := h + math.Sqrt(5)*k + l - m - n + o
		r3 := h + k + math.Sqrt(5)*l + m - n - o
		r4 := h - k + l + math.Sqrt(5)*m + n - o
		r5 := h - k - l + m + math.Sqrt(5)*n + o
		r6 := h + k - l - m + n + math.Sqrt(5)*o

		N := sq(r1, r2, r3, r4, r5, r6) / 20
		tt := calcTowTheta(wl, N, lc)

		list = append(list, fmt.Sprintf("(%s, %s, %s, %s, %s, %s) ~%.2f:", rec[0], rec[1], rec[2], rec[3], rec[4], rec[5], tt))
	}

	return &list, nil
}
