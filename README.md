<h1 align="center" id="top">Mandelbrot Magic ✨🌀</h1>

<p align="center">
    <strong>
        A magical ✨ web application written in Go 🐹 that lets you explore the infinite beauty of the Mandelbrot set interactively!
    </strong>
</p>

## What is the Mandelbrot Set? 🤔

The Mandelbrot set, discovered by Benoît Mandelbrot in the late 1970s, is one of the most famous examples of a **fractal** – a complex, infinitely detailed geometric shape. It lives in the **complex plane**, where numbers have both a real part and an imaginary part (like $a + bi$).

The set itself is defined by a very simple iterative process for each complex number $C$ in the plane:

$$ Z_{n+1} = Z_n^2 + C $$

We start with $Z_0 = 0$. Then we calculate $Z_1 = Z_0^2 + C = C$, then $Z_2 = Z_1^2 + C = C^2 + C$, then $Z_3 = Z_2^2 + C = (C^2+C)^2 + C$, and so on.

A complex number $C$ **belongs to the Mandelbrot set** if, when you repeatedly apply this formula, the magnitude (distance from zero) of the result $Z_n$ **never exceeds a certain bound** (usually 2), no matter how many times you iterate. If the magnitude stays small forever, the point $C$ is *inside* the set (typically colored black).

**What about the colors? 🌈**
If the magnitude of $Z_n$ *does* eventually exceed 2, the point $C$ is *outside* the set. The beautiful colors you often see in renderings are determined by **how quickly** the sequence escapes this bound. Points that escape faster might get one color, while points that take longer to escape get different colors, creating the intricate bands and filaments surrounding the main black shape. This is often called the "escape time algorithm".

**Why is it special? ✨**
*   **Infinite Complexity:** Despite the simple formula, the boundary of the Mandelbrot set is infinitely complex. Zooming into the boundary reveals ever more intricate details and smaller, distorted copies of the main shape (a property called **self-similarity**).
*   **Mathematical Beauty:** It connects various areas of mathematics in surprising ways.
*   **Computational Art:** It has become an icon of computer graphics and algorithmic art.

This project lets you explore this fascinating mathematical object right in your browser! 🎨

## Features 🚀

*   **Interactive Web UI:** Explore the fractal directly in your browser.
*   **Parameter Control:** Adjust position (X, Y), zoom level (Size ph), iteration depth, and anti-aliasing samples.
*   **Presets:** Jump to interesting pre-defined locations in the fractal.
*   **Live Rendering:** See the fractal update based on your parameters (shows a loading indicator).
*   **Pure Go Backend:** Uses Go's standard library for web serving and image generation.

## How to Run 🏃‍♀️💨

1.  **Prerequisites:** Ensure you have Go installed (see [golang.org](https://golang.org/dl/)).
2.  **Clone:** Clone this repository (or download the code).
    ```bash
    git clone https://github.com/VinsmokeSomya/Mandelbrot-Magic.git
    ```
3.  **Navigate:** Change into the project directory.
    ```bash
    cd Mandelbrot-Magic
    ```
4.  **Run:** Execute the application using the Go command.
    ```bash
    go run .
    ```
5.  **Explore:** Open your web browser and go to [http://localhost:8080](http://localhost:8080).

## Using the Interface 🖱️

*   Use the input fields on the left to change the rendering parameters.
*   Select a location from the "Presets ✨" dropdown to load specific settings.
*   Click the "Render!" button.
*   Wait for the "⏳ Rendering, please wait..." message to disappear and see the updated fractal on the right!

---
