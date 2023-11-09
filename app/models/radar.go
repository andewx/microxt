package models

type ADCData struct {
	if1a []int16
	if2b []int16
	if1b []int16
}

func NewADCData() *ADCData {
	return &ADCData{if1a: make([]int16, 256), if2b: make([]int16, 256), if1b: make([]int16, 256)}
}

type FFTData struct {
	data   []int16
	thresh []int16
}

func NewFFTData() *FFTData {
	return &FFTData{data: make([]int16, 256), thresh: make([]int16, 256)}
}

type PDAT_struct struct {
	distance uint16
	speed    uint16
	angle    int16
	dbs      uint16
}

func NewPDAT_struct() *PDAT_struct {
	return &PDAT_struct{}
}

type PDATData struct {
	data []DDAT_struct
}

func NewPDATData() *PDATData {
	return &PDATData{data: make([]DDAT_struct, 96)}
}

type TDAT_struct struct {
	distance uint16
	speed    uint16
	angle    int16
	dbs      uint16
}

type TDATData struct {
	data []TDAT_struct
}

func NewTDATData() *TDATData {
	return &TDATData{data: make([]TDAT_struct, 96)}
}

type DDAT_struct struct {
	detect uint8
	micro  uint8
	angle  uint8
	dir    uint8
	_range uint8
	speed  uint8
}

type DDATData struct {
	data []DDAT_struct
}

func NewDDATData() *DDATData {
	return &DDATData{data: make([]DDAT_struct, 96)}
}

type RadarParams struct {
	software_version   [19]byte
	base_freq          uint8
	max_speed          uint8
	max_range          uint8
	thresh_off         uint8
	tracking_filter    uint8
	vibration_suppress uint8
	min_detect_dist    uint8
	max_detect_dist    uint8
	min_detect_angle   int8
	max_detect_angle   int8
	min_detect_speed   uint8
	max_detect_speed   uint8
	detect_dir         uint8
	range_thresh       uint8
	angle_tresh        uint8
	speed_thresh       uint8
	d0                 uint8
	d1                 uint8
	d2                 uint8
	hold_time          uint16
	micro_retrigger    uint8
	micro_sensitivity  uint8
}
