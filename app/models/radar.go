package models

type ADCData struct {
	If1a []int16 `json:"if1a"`
	If2b []int16 `json:"if2b"`
	If1b []int16 `json:"if1b"`
}

func NewADCData() *ADCData {
	return &ADCData{If1a: make([]int16, 256), If2b: make([]int16, 256), If1b: make([]int16, 256)}
}

type FFTData struct {
	Data   []int16 `json:"data"`
	Thresh []int16 `json:"thresh"`
}

func NewFFTData() *FFTData {
	return &FFTData{Data: make([]int16, 256), Thresh: make([]int16, 256)}
}

type PDAT_struct struct {
	Distance uint16 `json:"distance"`
	Speed    uint16 `json:"speed"`
	Angle    int16  `json:"angle"`
	Dbs      uint16 `json:"dbs"`
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
	Distance uint16 `json:"distance"`
	Speed    uint16 `json:"speed"`
	Angle    int16  `json:"angle"`
	Dbs      uint16 `json:"dbs"`
}

type TDATData struct {
	data []TDAT_struct
}

func NewTDATData() *TDATData {
	return &TDATData{data: make([]TDAT_struct, 96)}
}

type DDAT_struct struct {
	Detect uint8 `json:"detect"`
	Micro  uint8 `json:"micro"`
	Angle  uint8 `json:"angle"`
	Dir    uint8 `json:"dir"`
	Range  uint8 `json:"range"`
	Speed  uint8 `json:"speed"`
}

type DDATData struct {
	data []DDAT_struct
}

func NewDDATData() *DDATData {
	return &DDATData{data: make([]DDAT_struct, 96)}
}

type RadarParams struct {
	software_version   [19]byte
	Base_freq          uint8  `json:"base_freq"`
	Max_speed          uint8  `json:"Max_speed"`
	Max_range          uint8  `json:"Max_range"`
	Thresh_off         uint8  `json:"thresh_off"`
	Tracking_filter    uint8  `json:"tracking_filter"`
	Vibration_suppress uint8  `json:"vibration_suppress"`
	Min_detect_dist    uint8  `json:"Min_detect_dist"`
	Max_detect_dist    uint8  `json:"Max_detect_dist"`
	Min_detect_angle   int8   `json:"Min_detect_angle"`
	Max_detect_angle   int8   `json:"Max_detect_angle"`
	Min_detect_speed   uint8  `json:"Min_detect_speed"`
	Max_detect_speed   uint8  `json:"Max_detect_speed"`
	Detect_dir         uint8  `json:"detect_dir"`
	Range_thresh       uint8  `json:"range_thresh"`
	Angle_tresh        uint8  `json:"angle_tresh"`
	Speed_thresh       uint8  `json:"speed_thresh"`
	D0                 uint8  `json:"d0"`
	D1                 uint8  `json:"d1"`
	D2                 uint8  `json:"d2"`
	Hold_time          uint16 `json:"hold_time"`
	Micro_retrigger    uint8  `json:"micro_retrigger"`
	Micro_sensitivity  uint8  `json:"micro_sensitivity"`
}
