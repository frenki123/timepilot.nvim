local M = {}

local function get_go_bin()
  local gobin = os.getenv("GOBIN")
  if gobin ~= nil then
    return gobin
  end
  local gopath = os.getenv("GOPATH")
  if gopath ~= nil then
    return gopath .. "/bin"
  end
  local home = os.getenv("HOME")
  return home .. "/go/bin"
end

M.config = {
  timepilot_path = get_go_bin() .. "/timepilot",
}

function M.set_config(user_config)
  for option, value in pairs(user_config) do
    M.config[option] = value
  end
end

return M
