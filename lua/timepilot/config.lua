local M = {}

M.config = {
  timepilot_path = "timepilot",
  autoinstall = true,
  disabled_filetypes = {
    "aerial",
    "alpha",
    "Avante",
    "checkhealth",
    "dapui*",
    "db*",
    "Diffview*",
    "Dressing*",
    "fugitive",
    "help",
    "httpResult",
    "lazy",
    "lspinfo",
    "mason",
    "minifiles",
    "Neogit*",
    "neo%-tree*",
    "neotest%-summary",
    "netrw",
    "noice",
    "notify",
    "NvimTree",
    "oil",
    "prompt",
    "qf",
    "query",
    "TelescopePrompt",
    "Trouble",
    "trouble",
    "VoltWindow",
    "undotree",
    "blink-cmp-menu",
  },
  timeout = 0.5,
}

function M.set_config(user_config)
  for option, value in pairs(user_config) do
    M.config[option] = value
  end
end

function M.check_config() end

return M
