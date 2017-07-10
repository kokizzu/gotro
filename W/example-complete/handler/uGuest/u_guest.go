package uGuest

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"github.com/kokizzu/gotro/X"
	"github.com/wcharczuk/go-chart"
	"os"
	"time"
)

// 2017-06-28 Prayogo
func timeFormatter(v interface{}) string {
	const dateFormat = `01-02 15:04`
	if val, ok := v.(time.Time); ok {
		return val.Format(dateFormat)
	}
	if val, ok := v.(int64); ok {
		return time.Unix(0, val).Format(dateFormat)
	}
	if val, ok := v.(float64); ok {
		return time.Unix(0, int64(val)).Format(dateFormat)
	}
	return fmt.Sprintf(`%# v`, v)
}

// 2017-06-28 Prayogo
func floatFormatter(v interface{}) string {
	if typed, isTyped := v.(float64); isTyped {
		return fmt.Sprintf(`%.3f`, typed)
	}
	return fmt.Sprintf(`%# v`, v)
}

// 2017-06-28 Prayogo
func CreatePDF(full_path string) {
	pdf := gofpdf.New(`P`, `mm`, `A4`, ``)
	pdf.AddPage()
	pdf.Image(full_path+`.png`, 10, 10, 190, 0, false, ``, 0, ``)
	err := pdf.OutputFileAndClose(full_path)
	L.IsError(err, `CreatePDF.OutputFileAndClose`)
}

// 2017-06-28 Prayogo
func CreateChart(result A.MSX, agg, full_path, sta_name string) {
	time_arr := []time.Time{}
	float_arr := []float64{}
	for _, row := range result {
		submitted_at := X.ToS(row[`submitted_at`])
		time_data := time.Time{}
		var err error
		if agg == `hour` || agg == `hh` || agg == `minute` || agg == `mm` {
			time_data, err = time.Parse(`2006-01-02 15:04`, submitted_at)
		} else if agg == `none` {
			time_data, err = time.Parse(`2006-01-02 15:04:05`, submitted_at)
		} else {
			time_data = time.Unix(S.ToI(submitted_at), 0)
		}
		if L.IsError(err, `Invalid format: `+submitted_at) {
			continue
		}
		time_arr = append(time_arr, time_data)
		float_arr = append(float_arr, X.ToF(row[`level_sensor`]))
	}
	graph := chart.Chart{
		Title:      sta_name,
		TitleStyle: chart.StyleTextDefaults(),
		XAxis: chart.XAxis{
			ValueFormatter: timeFormatter,
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			ValueFormatter: floatFormatter,
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: time_arr,
				YValues: float_arr,
			},
		},
	}
	var file, err = os.OpenFile(full_path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if L.IsError(err, `CreateChart.OpenFile: %s`, full_path) {
		return
	}
	defer file.Close()
	graph.Render(chart.PNG, file)
}
