local M = {}
local client = require("timepilot.client")
local config = require("timepilot.config")

function M.setup(user_config)
    user_config = user_config or {}
    config.set_config(user_config)
	vim.api.nvim_create_user_command("StartRPC", client.start, {})
	vim.api.nvim_create_user_command("SendEvent", client.send_event, {})
end
return M
