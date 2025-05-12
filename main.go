package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

// Configuration constants that were previously here for px, py, ph, imgWidth, etc., and ratio have been removed.

// REMAINING Configuration constants
const (
	linearMixing = true
	// showProgress = true // Not useful for web version
	profileCpu = false
)

// Default values for parameters
const (
	defaultWidth   = 512
	defaultHeight  = 512
	defaultPx      = -0.745
	defaultPy      = 0.113
	defaultPh      = 0.005
	defaultMaxIter = 1000
	defaultSamples = 10
)

// Preset definitions
type Preset struct {
	Name    string
	Px      float64
	Py      float64
	Ph      float64
	Iter    int
	Samples int
	W       int
	H       int
}

var presets = []Preset{
	{"Default View", -0.745, 0.113, 0.005, 1000, 10, 512, 512},
	{"Classic Full Set", -0.5, 0, 2.5, 500, 10, 512, 512},
	{"Seahorse Valley", -0.75, 0.11, 0.016, 1500, 50, 1024, 1024},
	{"Elephant Valley", 0.275, 0.005, 0.005, 1500, 50, 1024, 1024},
	{"Deep Zoom Point", -0.5557506, -0.55560, 0.000000001, 2500, 50, 1024, 1024},
}

// HTML form template
const htmlForm = `
<!DOCTYPE html>
<html>
<head>
    <title>Mandelbrot Magic ✨🌀</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
            margin: 0; /* Remove default margin */
            background-color: #f0f2f5;
            color: #333;
            display: flex;
            flex-direction: column; /* Stack title above main container */
            min-height: 100vh; /* Ensure body takes full viewport height */
        }
        h1 {
            color: #2c3e50;
            text-align: center;
            padding: 15px 0;
            margin: 0;
            background-color: #fff;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        #main-container {
            display: flex;
            flex-grow: 1; /* Allow main container to fill remaining space */
            padding: 20px;
        }
        #controls {
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-right: 20px; /* Space between controls and image */
            flex: 0 0 380px; /* Slightly wider controls */
            height: fit-content; /* Make controls container height fit its content */
        }
        #image-container {
            flex-grow: 1; /* Allow image container to take remaining space */
            display: flex;
            flex-direction: column;
            align-items: center; /* Center content horizontally */
            justify-content: center; /* Center content vertically */
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .form-row {
            display: flex;
            align-items: center;
            margin-bottom: 10px;
        }
        label {
            flex: 0 0 110px; /* Adjusted fixed width for labels */
            margin-right: 10px;
            font-weight: bold;
        }
        input[type=text], input[type=number], select {
            flex-grow: 1;
            padding: 8px 10px;
            border: 1px solid #ccc;
            border-radius: 4px;
            box-sizing: border-box;
            height: 36px; /* Ensure consistent height */
        }
        input[type=submit] {
            padding: 10px 20px;
            background-color: #3498db;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 1em;
            transition: background-color 0.2s;
            display: block;
            width: 100%;
            margin-top: 15px;
        }
        input[type=submit]:hover {
            background-color: #2980b9;
        }
        #loadingMessage {
            text-align: center;
            font-weight: bold;
            margin-top: 10px; /* Add some space */
            min-height: 1.2em; /* Reserve space */
        }
        #fractalImage {
            display: block; /* Remove auto margin for centering, container does it */
            max-width: 100%; /* Ensure image doesn't overflow container */
            max-height: 70vh; /* Limit image height for better layout */
            border: 1px solid #ddd;
            border-radius: 4px;
            /* background-color: #fff; - Removed, container has background */
            /* padding: 5px; - Removed, container has padding */
            box-shadow: 0 1px 3px rgba(0,0,0,0.05);
            margin-top: 10px; /* Space above image */
        }
        hr { display: none; } /* Hide the hr, less relevant in this layout */
        p {
             text-align: center; 
             margin-top: 0; /* Reduce space above note */
             margin-bottom: 10px;
             font-size: 0.9em;
             color: #555;
        }
        select {
             /* Make dropdown look more like other inputs */
             appearance: none; /* Remove default system appearance */
             /* Escaped % signs in SVG URL for Fprintf */
             background-image: url('data:image/svg+xml;charset=US-ASCII,%%3Csvg%%20xmlns%%3D%%22http%%3A%%2F%%2Fwww.w3.org%%2F2000%%2Fsvg%%22%%20width%%3D%%22292.4%%22%%20height%%3D%%22292.4%%22%%3E%%3Cpath%%20fill%%3D%%22%%23333%%22%%20d%%3D%%22M287%%2069.4a17.6%%2017.6%%200%%200%%200-13-5.4H18.4c-5%%200-9.3%%201.8-12.9%%205.4A17.6%%2017.6%%200%%200%%200%%200%%2082.2c0%%205%%201.8%%209.3%%205.4%%2012.9l128%%20127.9c3.6%%203.6%%207.8%%205.4%%2012.8%%205.4s9.2-1.8%%2012.8-5.4L287%%2095c3.5-3.5%%205.4-7.8%%205.4-12.8%%200-5-1.9-9.2-5.5-12.8z%%22%%2F%%3E%%3C%%2Fsvg%%3E');
             background-repeat: no-repeat;
             background-position: right .7em top 50%%; /* Escaped the 50% to 50%% */
             background-size: .65em auto;
             cursor: pointer;
        }
    </style>
</head>
<body>
    <h1>Mandelbrot Magic ✨🌀</h1>

    <div id="main-container">
        <div id="controls">
            <div class="form-row">
                <label for="presetSelect">Presets ✨:</label>
                <select id="presetSelect">
                    <option value="-1">-- Select a Preset --</option>
                    %[1]s
                </select>
            </div>

            <form id="renderForm" action="/" method="get">
                <div class="form-row">
                     <label for="w">Width ↔️:</label> <input type="number" id="w" name="w" value="%[2]d">
                </div>
                <div class="form-row">
                     <label for="h">Height ↕️:</label> <input type="number" id="h" name="h" value="%[3]d">
                </div>
                 <div class="form-row">
                     <label for="px">Pos X 📍:</label> <input type="text" id="px" name="px" value="%[4]f">
                </div>
                <div class="form-row">
                    <label for="py">Pos Y 📍:</label> <input type="text" id="py" name="py" value="%[5]f">
                </div>
                <div class="form-row">
                    <label for="ph">Size (ph) 🔍:</label> <input type="text" id="ph" name="ph" value="%[6]f">
                </div>
                <div class="form-row">
                    <label for="iter">Iterations 🔁:</label> <input type="number" id="iter" name="iter" value="%[7]d">
                </div>
                <div class="form-row">
                    <label for="samples">Samples ✨:</label> <input type="number" id="samples" name="samples" value="%[8]d">
                </div>
                <input type="submit" value="Render!">
            </form>
        </div>

        <div id="image-container">
            <p>Note: Rendering might take a few seconds depending on parameters.</p>
            <div id="loadingMessage" style="display:none; color: orange;">⏳ Rendering, please wait...</div>
            <img id="fractalImage" src="/render?w=%[2]d&h=%[3]d&px=%[4]f&py=%[5]f&ph=%[6]f&iter=%[7]d&samples=%[8]d" alt="Fractal Image">
        </div>
    </div>

    <script>
        const presetsData = %[9]s; // Explicitly placeholder 9 for JSON
        const presetSelect = document.getElementById('presetSelect');
        const form = document.getElementById('renderForm');
        const loadingMessage = document.getElementById('loadingMessage');
        const fractalImage = document.getElementById('fractalImage');

        // Input fields
        const inputW = document.getElementById('w');
        const inputH = document.getElementById('h');
        const inputPx = document.getElementById('px');
        const inputPy = document.getElementById('py');
        const inputPh = document.getElementById('ph');
        const inputIter = document.getElementById('iter');
        const inputSamples = document.getElementById('samples');

        presetSelect.addEventListener('change', function() {
            const selectedIndex = parseInt(this.value, 10);
            if (selectedIndex >= 0 && selectedIndex < presetsData.length) {
                const preset = presetsData[selectedIndex];
                inputW.value = preset.W;
                inputH.value = preset.H;
                inputPx.value = preset.Px;
                inputPy.value = preset.Py;
                inputPh.value = preset.Ph;
                inputIter.value = preset.Iter;
                inputSamples.value = preset.Samples;
            }
        });

        form.addEventListener('submit', function(event) {
            event.preventDefault(); // Stop default page reload
            loadingMessage.textContent = 'Rendering, please wait...';
            loadingMessage.style.display = 'block';

            const params = new URLSearchParams();
            params.append('w', inputW.value);
            params.append('h', inputH.value);
            params.append('px', inputPx.value);
            params.append('py', inputPy.value);
            params.append('ph', inputPh.value);
            params.append('iter', inputIter.value);
            params.append('samples', inputSamples.value);
            
            fractalImage.src = '/render?' + params.toString();
        });

        fractalImage.addEventListener('load', function() {
            loadingMessage.style.display = 'none';
        });

        fractalImage.addEventListener('error', function() {
            loadingMessage.textContent = 'Error loading image. Please try again.';
        });
    </script>
</body>
</html>
`

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/render", handleRender)

	port := "8080"
	log.Printf("Starting server on http://localhost:%s ...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Handler to display the form
func handleForm(w http.ResponseWriter, r *http.Request) {
	// Parse parameters or use defaults
	imgWidth := getIntParamHelper(r, "w", defaultWidth)
	imgHeight := getIntParamHelper(r, "h", defaultHeight)
	px := getFloatParamHelper(r, "px", defaultPx)
	py := getFloatParamHelper(r, "py", defaultPy)
	ph := getFloatParamHelper(r, "ph", defaultPh)
	maxIter := getIntParamHelper(r, "iter", defaultMaxIter)
	samples := getIntParamHelper(r, "samples", defaultSamples)

	// Generate preset options HTML
	var presetOptions strings.Builder
	for i, p := range presets {
		fmt.Fprintf(&presetOptions, `<option value="%d">%s</option>`, i, p.Name)
	}

	// Inject presets data as JSON for JavaScript
	presetsJSON, err := json.Marshal(presets)
	if err != nil {
		log.Printf("Error marshaling presets to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Populate the HTML form template using explicit numbered placeholders
	fmt.Fprintf(w, htmlForm,
		presetOptions.String(),
		imgWidth,
		imgHeight,
		px,
		py,
		ph,
		maxIter,
		samples,
		string(presetsJSON),
	)
}

// Helper function to get int query param, moved outside handleRender for use by handleForm
func getIntParamHelper(r *http.Request, name string, defaultValue int) int {
	valStr := r.URL.Query().Get(name)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s: %s. Using default %d.", name, valStr, defaultValue)
		return defaultValue
	}
	return val
}

// Helper function to get float query param, moved outside handleRender for use by handleForm
func getFloatParamHelper(r *http.Request, name string, defaultValue float64) float64 {
	valStr := r.URL.Query().Get(name)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		log.Printf("Warning: Invalid float value for %s: %s. Using default %f.", name, valStr, defaultValue)
		return defaultValue
	}
	return val
}

// Handler to parse parameters, render the fractal, and return the image
func handleRender(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /render request...")
	start := time.Now()

	// Parse parameters from query string using helper functions
	imgWidth := getIntParamHelper(r, "w", defaultWidth)
	imgHeight := getIntParamHelper(r, "h", defaultHeight)
	px := getFloatParamHelper(r, "px", defaultPx)
	py := getFloatParamHelper(r, "py", defaultPy)
	ph := getFloatParamHelper(r, "ph", defaultPh)
	maxIter := getIntParamHelper(r, "iter", defaultMaxIter)
	samples := getIntParamHelper(r, "samples", defaultSamples)

	// Basic validation/clamping (optional but good practice)
	if imgWidth <= 0 {
		imgWidth = defaultWidth
	}
	if imgHeight <= 0 {
		imgHeight = defaultHeight
	}
	if maxIter <= 0 {
		maxIter = defaultMaxIter
	}
	if samples <= 0 {
		samples = defaultSamples
	}
	if ph <= 0 {
		ph = defaultPh
	} // Size must be positive

	log.Printf("Rendering with params: w=%d, h=%d, px=%f, py=%f, ph=%f, iter=%d, samples=%d",
		imgWidth, imgHeight, px, py, ph, maxIter, samples)

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Render the fractal (using the refactored function)
	render(img, imgWidth, imgHeight, px, py, ph, maxIter, samples)
	renderEnd := time.Now()
	log.Printf("Rendering finished in %s", renderEnd.Sub(start))

	// Encode image to PNG in memory
	var imgBuffer bytes.Buffer
	err := png.Encode(&imgBuffer, img)
	if err != nil {
		log.Printf("Error encoding PNG: %v", err)
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
	encodeEnd := time.Now()
	log.Printf("Encoding finished in %s", encodeEnd.Sub(renderEnd))

	// Set response header and write image data
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(imgBuffer.Len()))
	_, err = imgBuffer.WriteTo(w)
	if err != nil {
		log.Printf("Error writing image to response: %v", err)
		// Don't write another header if one was already sent
	}

	log.Printf("Served /render request in %s", time.Since(start))
}

func render(img *image.RGBA, imgWidth, imgHeight int, px, py, ph float64, maxIter, samples int) {
	if profileCpu {
		f, err := os.Create("profile.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	jobs := make(chan int)
	ratio := float64(imgWidth) / float64(imgHeight)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for y := range jobs {
				for x := 0; x < imgWidth; x++ {
					var r, g, b int
					for i := 0; i < samples; i++ {
						nx := ph*ratio*((float64(x)+RandFloat64())/float64(imgWidth)) + px
						ny := ph*((float64(y)+RandFloat64())/float64(imgHeight)) + py
						c := paint(mandelbrotIter(nx, ny, maxIter))
						if linearMixing {
							r += int(RGBToLinear(c.R))
							g += int(RGBToLinear(c.G))
							b += int(RGBToLinear(c.B))
						} else {
							r += int(c.R)
							g += int(c.G)
							b += int(c.B)
						}
					}
					var cr, cg, cb uint8
					if linearMixing {
						cr = LinearToRGB(uint16(float64(r) / float64(samples)))
						cg = LinearToRGB(uint16(float64(g) / float64(samples)))
						cb = LinearToRGB(uint16(float64(b) / float64(samples)))
					} else {
						cr = uint8(float64(r) / float64(samples))
						cg = uint8(float64(g) / float64(samples))
						cb = uint8(float64(b) / float64(samples))
					}
					img.SetRGBA(x, y, color.RGBA{R: cr, G: cg, B: cb, A: 255})
				}
			}
		}()
	}

	for y := 0; y < imgHeight; y++ {
		jobs <- y
	}
	close(jobs)
}

func paint(r float64, n int) color.RGBA {
	insideSet := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	if r > 4 {
		return hslToRGB(float64(n)/800*r, 1, 0.5)
	}

	return insideSet
}

func mandelbrotIter(px, py float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + px
		y = 2*xy + py
	}

	return xx + yy, maxIter
}

/*

func mandelbrotIterComplex(px, py float64, maxIter int) (float64, int) {
	var current complex128
	pxpy := complex(px, py)

	for i := 0; i < maxIter; i++ {
		magnitude := cmplx.Abs(current)
		if magnitude > 2 {
			return magnitude * magnitude, i
		}
		current = current * current + pxpy
	}

	magnitude := cmplx.Abs(current)
	return magnitude * magnitude, maxIter
}

*/
