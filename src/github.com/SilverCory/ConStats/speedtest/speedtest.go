package speedtest

import (
	"encoding/json"
	"os/exec"
)

// SpeedTest the speed test instance
type SpeedTest struct {
	Command string
	Args    []string
}

// TestResult the results of the test.
type TestResult struct {
	Upload    float32 `json:"upload"`
	Download  float32 `json:"download"`
	Ping      float32 `json:"ping"`
	TimeStamp string  `json:"timestamp"`
}

// Create creates an instance of SpeedTest
func Create() *SpeedTest {
	return &SpeedTest{
		Command: "speedtest",
		Args:    []string{"--secure", "--json"},
	}
}

// Test run the actual test.
func (s *SpeedTest) Test() (*TestResult, error) {

	out, err := exec.Command(s.Command, s.Args...).Output()
	if err != nil {
		return nil, err
	}

	result := &TestResult{}
	err = json.Unmarshal(out, result)

	return result, err

}
