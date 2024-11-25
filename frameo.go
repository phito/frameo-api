package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func isScreenAwake() (bool, error) {
	cmd := exec.Command("adb", "shell", "dumpsys", "power")

	out, err := cmd.CombinedOutput()

	if err != nil {
		return false, err
	}

	output := string(out)
	if strings.Contains(output, "mWakefulness=Awake") {
		return true, nil
	} else if strings.Contains(output, "mWakefulness=Asleep") {
		return false, nil
	}

	return false, fmt.Errorf("unable to determine screen state")
}

func toggleScreen() error {
	cmd := exec.Command("adb", "shell", "input", "keyevent", "26")
	return cmd.Run()
}

func getBrightness() (int, error) {
	cmd := exec.Command("adb", "shell", "dumpsys", "power")

	out, err := cmd.CombinedOutput()

	if err != nil {
		return 0, err
	}

	output := string(out)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.Contains(line, "mScreenBrightnessSetting=") {
			brightnessValue := strings.TrimSpace(strings.Split(line, "=")[1])
			return strconv.Atoi(brightnessValue)
		}
	}

	return 0, fmt.Errorf("unable to find screen brightness")
}

func handleBrightness(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		brightness, err := getBrightness()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(writer, "%d", brightness)
		writer.WriteHeader(http.StatusOK)
		return
	}

	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Unable to read request body", http.StatusBadRequest)
		return
	}

	defer request.Body.Close()

	brightness, err := strconv.Atoi(string(body))
	if err != nil {
		http.Error(writer, "Invalid integer", http.StatusBadRequest)
		return
	}

	if brightness < 0 || brightness > 255 {
		http.Error(writer, "Brightness must be between 0 and 255", http.StatusBadRequest)
		return
	}

	fmt.Println("Setting brightness to", brightness)

	cmd := exec.Command("adb", "shell", "settings", "put", "system", "screen_brightness", strconv.Itoa(brightness))
	err = cmd.Run()

	if err != nil {
		http.Error(writer, "Error occured while setting brightness", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func handleScreen(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		isScreenAwake, err := isScreenAwake()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(writer, "%t", isScreenAwake)
		writer.WriteHeader(http.StatusOK)
		return
	}

	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Unable to read request body", http.StatusBadRequest)
		return
	}

	defer request.Body.Close()

	bodyStr := string(body)
	shouldToggleScreen := false

	if bodyStr == "" {
		shouldToggleScreen = true
	} else {
		bodyBool, err := strconv.ParseBool(bodyStr)
		if err != nil {
			http.Error(writer, "Body should be either empty or boolean", http.StatusBadRequest)
			return
		}

		isScreenAwake, err := isScreenAwake()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		shouldToggleScreen = bodyBool != isScreenAwake
	}

	if shouldToggleScreen {
		err := toggleScreen()

		if err != nil {
			http.Error(writer, "Error while toggling sleep mode", http.StatusInternalServerError)
			return
		}
	}
	writer.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/brightness", handleBrightness)
	http.HandleFunc("/screen", handleScreen)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
