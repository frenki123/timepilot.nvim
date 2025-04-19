local M = {}

M.config = {
  timepilot_path = "timepilot",
  autoinstall = true,
}

function M.set_config(user_config)
  for option, value in pairs(user_config) do
    M.config[option] = value
  end
end

function M.check_config() end

return M
