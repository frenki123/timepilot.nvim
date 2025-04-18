local client = require("timepilot.client")

vim.api.nvim_create_user_command("StartRPC", client.start, {})
vim.api.nvim_create_user_command("SendEvent", client.send_event, {})
