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
