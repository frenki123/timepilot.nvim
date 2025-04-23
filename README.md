# timepilot.nvim

Track coding time and stuff in Neovim

## 🚀 Installation

Using [lazy.nvim](https://github.com/folke/lazy.nvim):

```lua
{
	"frenki123/timepilot.nvim",
    build = function()
        require("timepilot").build() -- to install timepilot daemon
    end,
	config = function()
		require("timepilot").setup()
	end,
}
```

## Configuration
```lua
require("timepilot").setup({
  timepilot_path = "timepilot", -- Path to the `timepilot` daemon binary
  autoinstall = true,           -- Auto-install the daemon if missing
  timeout = 0.5,                -- Idle timeout in minutes
  debug = false,                -- Show debug notifications
})
```

## Usage
Commands 
1. `:GetProjectTime` - get total time spent on the project (project is git root directory or directory where you started neovim)
2. `:GetFileTime` - get the file on which you spent most time

## Disclaimer

Still WIP – bugs possible, especially with IDLE detection. Tests missing!
