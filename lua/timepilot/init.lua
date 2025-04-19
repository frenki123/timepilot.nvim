local M = {}
local client = require("timepilot.client")
local config = require("timepilot.config")

function M.build()
    config.build()
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
