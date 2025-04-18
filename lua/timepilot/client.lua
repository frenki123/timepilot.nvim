local M = {}

local config = require("timepilot.config").config
local uv = vim.uv
local handle, stdin, stdout

function M.start()
	if handle then
		return
	end
	stdin = uv.new_pipe(false)
	stdout = uv.new_pipe(false)
	handle = uv.spawn(
		config.timepilot_path,
		{ stdio = { stdin, stdout, nil } },
		function(code, signal)
			print("EXITED with code", code, "signal", signal)
			handle:close()
		end
	)
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
			print("STDOUT:", vim.inspect(data))
		else
			print("No data (maybe CLI exited)")
		end
	end)
end

function M.send_event()
	if not stdin then
		print("NOT RUNNING RPC")
		return
	end

	local msg = {
		jsonrpc = "2.0",
		id = 1,
		method = "track_event",
		params = {
			type = "insert_enter",
			filepath = vim.fn.expand("%:p"),
			timestamp = os.time(),
		},
	}

	local json = vim.fn.json_encode(msg) .. "\n"
	stdin:write(json)
end

return M
