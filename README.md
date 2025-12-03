# HyprTask

> A lightweight, terminal-based task manager built specifically for the Hyprland window manager.

![License](https://img.shields.io/badge/license-MIT-blue.svg) ![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg) ![Hyprland](https://img.shields.io/badge/Hyprland-Supported-green)

**HyprTask** bridges the gap between your terminal and your window manager. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework, it provides a responsive, keyboard-driven interface to monitor workspaces, manage processes, and optimize your workflow without leaving the command line.

## ✨ Features

* **⚡ Real-Time Monitoring**: Live updates of your active workspaces and running processes.
* **🖥️ Responsive UI**: Dynamic padding and layout adjustments that respect your terminal dimensions.
* **⌨️ Keyboard-Centric**: Fully navigable using intuitive keybindings—no mouse required.
* **🔍 Workspace Selector**: Quickly filter and view processes specific to individual Hyprland workspaces.
* **🎨 Beautiful TUI**: Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) for a modern, clean aesthetic.

## 🚀 Getting Started

### Prerequisites

* **Go 1.25+**: Ensure you have Go installed.
* **Hyprland**: This tool is designed to interact directly with the Hyprland compositor.

### Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/paulvinueza30/hyprtask.git](https://github.com/paulvinueza30/hyprtask.git)
    cd hyprtask
    ```

2.  **Install dependencies:**
    ```bash
    go mod download
    ```

3.  **Build the binary:**
    ```bash
    go build -o hyprtask cmd/hyprtask/main.go
    ```

4.  **Run:**
    ```bash
    ./hyprtask
    ```

*> Note: The application is optimized for terminals with a minimum size of **65x20** characters.*

## 🤝 Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

### How to Contribute

1.  **Fork the Project**
2.  **Create your Feature Branch** (`git checkout -b feature/AmazingFeature`)
3.  **Commit your Changes** (`git commit -m 'Add some AmazingFeature'`)
4.  **Push to the Branch** (`git push origin feature/AmazingFeature`)
5.  **Open a Pull Request**

## 🙏 Acknowledgments

* [Charm](https://charm.sh/) for the incredible Bubble Tea and Lipgloss libraries.
* The Hyprland community for the window manager innovation.
