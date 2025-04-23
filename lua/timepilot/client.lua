local M = {}

local config = require("timepilot.config").config
local uv = vim.uv
local handle, stdin, stdout
local id = 0

local function notify(msg, level)
    vim.schedule(function ()
      vim.notify(msg, level)
    end)
end

local function print_response(data)
  local ok, decoded = pcall(vim.json.decode, data)
  if not ok then
    local msg = string.format("Failed to decode json: '%s'", vim.inspect(data))
    notify(msg, vim.log.levels.WARN)
  end
  local result = decoded.result
  if not result then
    notify("Unknown Result", vim.log.levels.WARN)
    return
  end
  local kind = result.type
  local res_data = result.data
  if kind == "DEBUG" and config.debug then
    notify(res_data, vim.log.levels.DEBUG)
    return
  end
  if kind == "INFO/TIME" then
    local msg = string.format("Total time in project: %.2f min", res_data)
    notify(msg, vim.log.levels.INFO)
    return
  end
  if kind == "INFO/FILE" then
    local msg = string.format("Most edited file:\n- %s\n- %s\n- %s min", res_data.filepath, res_data.filetype, res_data.time)
    notify(msg, vim.log.levels.INFO)
    return
  end
end

function M.start()
  if handle then
    return
  end
  stdin = uv.new_pipe(false)
  stdout = uv.new_pipe(false)
  handle = uv.spawn(config.timepilot_path, { stdio = { stdin, stdout, nil } }, function(code, signal)
    print("EXITED with code", code, "signal", signal)
    handle:close()
  end)
  if handle == nil then
    local msg = "Timepilot process not started! "
    msg = msg .. "Is timepilot installed at '" .. config.timepilot_path
    msg = msg .. "'?"
    vim.notify(msg)
    return
  end
  vim.notify_once("Timepilot started")
  stdout:read_start(function(err, data)
    if err then
      print("STDOUT ERROR:", err)
    elseif data then
      print_response(data)
    else
      print("No data (maybe CLI exited)")
    end
  end)
end

function M.send_event(method, params)
  if not stdin then
    print("NOT RUNNING RPC")
    return
  end

  id = id + 1
  local res = {
    jsonrpc = "2.0",
    id = id,
    method = method,
    params = params,
  }
  local json = vim.fn.json_encode(res) .. "\n"
  stdin:write(json)
end

return M
