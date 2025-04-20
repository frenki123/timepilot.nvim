local M = {}
local config = require("timepilot.config")
local build = require("timepilot.build").build
local autocmd = require("timepilot.autocmds")

function M.build()
  build()
end

function M.setup(user_config)
  user_config = user_config or {}
  config.set_config(user_config)
  autocmd.autocmd()
end
return M
