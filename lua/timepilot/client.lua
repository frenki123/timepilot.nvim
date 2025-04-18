local M = {}

local uv = vim.uv
local handle, stdin, stdout, pid

function M.start()
	if handle then
		return
	end
	stdin = uv.new_pipe(false)
	stdout = uv.new_pipe(false)
	handle, pid = uv.spawn(
		"/home/frenki/projects/neovim/timepilot.nvim/bin/timepilot",
		{ stdio = { stdin, stdout, nil } },
		function(code, signal)
			print("EXITED with code", code, "signal", signal)
			handle:close()
		end
	)
	print("process opened", handle, pid)
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
