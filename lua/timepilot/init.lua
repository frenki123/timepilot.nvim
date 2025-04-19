local M = {}
local client = require("timepilot.client")
local config = require("timepilot.config")

local function install_timepilot()
  if M.config.autoinstall then
    return os.execute("go install github.com/frenki123/timepilot.nvim/timepilot@latest")
  end
  return false
end

function M.build()
  local required_version = "0.0.1"
  local timepilot = M.config.timepilot_path
  local handle = io.popen(timepilot .. " version")
  if not handle then
    if not install_timepilot() then
      print("Timepilot not installed")
      return false
    end
    return
  end
  local result = handle:read("*a")
  handle:close()
  local version = result:match("v(.+)")
  if version == required_version then
    return
  end
  if not install_timepilot() then
    print("Required Timepilot version is" .. required_version .. "installed" .. version)
  end
end

function M.setup(user_config)
  user_config = user_config or {}
  config.set_config(user_config)
  config.check_config()
  client.start()
  vim.api.nvim_create_autocmd("InsertEnter", {
    callback = function()
      local filename = vim.api.nvim_buf_get_name(0)
      local basename = vim.fn.fnamemodify(filename, ":t")
      local params = {
        filetype = vim.bo.filetype,
        filename = basename,
        action = "InsertEnter",
      }
      client.send_event("action", params)
    end,
  })

  vim.api.nvim_create_user_command("StartRPC", client.start, {})
end
return M
