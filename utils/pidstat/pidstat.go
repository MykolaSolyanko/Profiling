package pidstat

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type PidstatCPUInfo struct {
	UID     int     `pidstat:"UID"`
	PID     int     `pidstat:"PID"`
	Usr     float32 `pidstat:"%usr"`
	System  float32 `pidstat:"%system"`
	Guest   float32 `pidstat:"%guest"`
	Wait    float32 `pidstat:"%wait"`
	CPU     float32 `pidstat:"%CPU"`
	CPUID   int     `pidstat:"CPU"`
	Command string  `pidstat:"Command"`
}

type PidstatMemoryInfo struct {
	UID     int     `pidstat:"UID"`
	PID     int     `pidstat:"PID"`
	Minflt  float32 `pidstat:"minflt/s"`
	Majflt  float32 `pidstat:"majflt/s"`
	VSZ     int64   `pidstat:"VSZ"`
	RSS     int64   `pidstat:"RSS"`
	MEM     float32 `pidstat:"%MEM"`
	Command string  `pidstat:"Command"`
}

var errNewPidstatTable = fmt.Errorf("New Pidstat table")

func GetPidstatMemoryMonitoring(pid, interval string, stat chan<- PidstatMemoryInfo, errs chan<- error) {
	cmd := exec.Command("pidstat", "-r", "-p", pid, interval)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errs <- err

		return
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		errs <- err

		return
	}

	r := bufio.NewReader(stdout)

	// skip first line system info
	if _, err = r.ReadString('\n'); err != nil {
		errs <- err

		return
	}

	// skip empty line
	if _, err = r.ReadString('\n'); err != nil {
		errs <- err

		return
	}

	// parse header line
	line, err := r.ReadString('\n')
	if err != nil {
		errs <- err

		return
	}

	header := strings.Fields(line)

	for {
		if line, err = r.ReadString('\n'); err != nil {
			errs <- err

			return
		}

		fields, err := parseLine(line, header)
		if err != nil {
			if errors.Is(err, errNewPidstatTable) {
				log.Infof("New Pidstat table")
				continue
			}

			errs <- err

			return
		}

		pRet := &PidstatMemoryInfo{}

		pErrs := fillLine(fields, pRet)
		if len(pErrs) > 0 {
			errs <- fmt.Errorf("Could not parse line: %s", pErrs)

			return
		}

		pRet.VSZ /= 1024
		pRet.RSS /= 1024

		stat <- *pRet
	}
}

func GetPidstatCPUMonitoring(pid, interval string, stat chan<- PidstatCPUInfo, errs chan<- error) {
	cmd := exec.Command("pidstat", "-u", "-p", pid, interval)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errs <- err

		return
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		errs <- err

		return
	}

	r := bufio.NewReader(stdout)

	// skip first line system info
	if _, err = r.ReadString('\n'); err != nil {
		errs <- err

		return
	}

	// skip empty line
	if _, err = r.ReadString('\n'); err != nil {
		errs <- err

		return
	}

	// parse header line
	line, err := r.ReadString('\n')
	if err != nil {
		errs <- err

		return
	}

	header := strings.Fields(line)

	for {
		if line, err = r.ReadString('\n'); err != nil {
			errs <- err

			return
		}

		fields, err := parseLine(line, header)
		if err != nil {
			if errors.Is(err, errNewPidstatTable) {
				log.Infof("New Pidstat table")
				continue
			}

			errs <- err

			return
		}

		pRet := &PidstatCPUInfo{}

		pErrs := fillLine(fields, pRet)
		if len(pErrs) > 0 {
			errs <- fmt.Errorf("Could not parse line: %s", pErrs)

			return
		}

		stat <- *pRet
	}
}

func fillLine(data map[string]string, pRet any) (errs []error) {
	errs = []error{}

	sv := reflect.Indirect(reflect.ValueOf(pRet))
	st := sv.Type()
	for i := 0; i < st.NumField(); i++ {
		fieldType := st.Field(i)
		fieldName, ok := fieldType.Tag.Lookup("pidstat")
		if !ok {
			continue
		}

		val, ok := data[fieldName]
		if !ok {
			errs = append(errs, fmt.Errorf("Missing field  %s", fieldName))

			continue
		}
		delete(data, fieldName)

		field := sv.FieldByIndex(fieldType.Index)

		switch fieldType.Type.Kind() {
		case reflect.String:
			field.SetString(val)

		case reflect.Float32:
			pVal, err := strconv.ParseFloat(val, 32)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: could not parse: %s",
					fieldName, err))

				continue
			}

			field.SetFloat(float64(pVal))

		case reflect.Int:
			pVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: could not parse: %s",
					fieldName, err))

				continue
			}

			field.SetInt(pVal)

		case reflect.Int64:
			pVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: could not parse: %s",
					fieldName, err))

				continue
			}

			field.SetInt(pVal)

		default:
			errs = append(errs, fmt.Errorf("Encountered unexpected fieldtype in Line struct"))
		}
	}

	return errs
}

func parseLine(line string, header []string) (map[string]string, error) {
	ret := make(map[string]string)

	fields := strings.Fields(line)

	if len(fields) == 0 || fields[1] == header[1] {
		log.Info(errNewPidstatTable)

		return nil, errNewPidstatTable
	}

	if len(header) != len(fields) {
		return nil, fmt.Errorf("header and fields length mismatch")
	}

	for i := 1; i < len(header); i++ {
		ret[header[i]] = fields[i]
	}

	return ret, nil
}
