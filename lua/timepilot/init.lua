local timepilot = require("timepilot.hello")

vim.api.nvim_create_user_command("Hello", timepilot.say_hello, {})

return timepilot
