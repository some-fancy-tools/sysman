package sysman

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type SysManager struct {
	mux *http.ServeMux
}

func New() (*SysManager, error) {
	return &SysManager{
		mux: http.NewServeMux(),
	}, nil
}

type Governor struct {
	CPUID              string   `json:"cpuId"`
	AvailableGovernors []string `json:"availableGovernors"`
	SelectedGovernor   string   `json:"selectedGovernor"`
}

func (s *SysManager) GetScalingGovernors() (map[string]string, map[string][]string, error) {
	re, err := regexp.Compile(`cpu\d+`)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	availableGovernorsMap := map[string][]string{}
	selectedGovernorMap := map[string]string{}
	filepath.Walk("/sys/devices/system/cpu/", func(path string, d fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if path == "/sys/devices/system/cpu/" {
			return nil
		}
		if !strings.HasPrefix(d.Name(), "cpu") || !strings.HasSuffix(path, "cpufreq") {
			// log.Println("skipping")
			// fmt.Println(path)
			return nil
		}
		cpuId := re.FindString(path)

		if cpuId == "" {
			return nil
		}
		// log.Println(path)
		dd, err := os.ReadFile(filepath.Join(path, "scaling_available_governors"))
		if err != nil {
			log.Println(err)
			return err
		}

		availableGovernorsMap[cpuId] = strings.Split(strings.TrimSuffix(string(dd[:]), " \n"), " ")

		ddc, err := os.ReadFile(filepath.Join(path, "scaling_governor"))
		if err != nil {
			log.Println(err)
			return err
		}
		selectedGovernorMap[cpuId] = strings.TrimSuffix(string(ddc[:]), "\n")

		return nil
	})
	return selectedGovernorMap, availableGovernorsMap, nil
}

func (s *SysManager) SetScalingGovernor(selection []Governor) error {
	for _, governor := range selection {
		spath := filepath.Join("/sys/devices/system/cpu", governor.CPUID, "cpufreq", "scaling_governor")
		// log.Printf("writing %s to %s", governor.SelectedGovernor, spath)
		if err := os.WriteFile(spath, []byte(governor.SelectedGovernor), 0644); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (s *SysManager) GetScalingFrequencies() (map[string][]int64, error) {
	re, err := regexp.Compile(`cpu\d+`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	freqMap := map[string][]int64{}
	err = filepath.Walk("/sys/devices/system/cpu/", func(path string, d fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if path == "/sys/devices/system/cpu/" {
			return nil
		}
		if !strings.HasPrefix(d.Name(), "cpu") || !strings.HasSuffix(path, "cpufreq") {
			// log.Println("skipping")
			// fmt.Println(path)
			return nil
		}
		cpuId := re.FindString(path)

		if cpuId == "" {
			return nil
		}
		freqMap[cpuId] = []int64{}
		// log.Println(path)
		for _, p := range []string{"scaling_cur_freq", "scaling_min_freq", "scaling_max_freq"} {
			dd, err := os.ReadFile(filepath.Join(path, p))
			if err != nil {
				log.Println(err)
				return err
			}
			freq, err := strconv.ParseInt(strings.TrimSuffix(string(dd[:]), "\n"), 10, 64)
			if err != nil {
				log.Println(err)
				return err
			}
			freqMap[cpuId] = append(freqMap[cpuId], freq)
		}
		// freqMap[cpuId] =
		return nil
	})

	return freqMap, err
}
